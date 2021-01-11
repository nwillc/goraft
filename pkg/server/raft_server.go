package server

import (
	"context"
	"fmt"
	"github.com/nwillc/goraft/database"
	"github.com/nwillc/goraft/model"
	"github.com/nwillc/goraft/raftapi"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
	"os"
	"time"
)

// Role currently of a RaftServer
type Role string

const (
	// Candidate Role
	Candidate Role = "CANDIDATE"
	// Leader Role
	Leader Role = "LEADER"
	// Follower Role
	Follower Role = "FOLLOWER"
)

// RaftServer holds data needed by a Raft Server
type RaftServer struct {
	raftapi.UnimplementedRaftServiceServer
	member             model.Member
	lastHeartbeat      time.Time
	role               Role
	peers              []model.Member
	votedOn            uint64
	databasePath       string
	electionCountdown  time.Duration
	heartbeatCountdown time.Duration
	statusRepo         *database.StatusRepository
	logRepo            *database.LogEntryRepository
	ctx                context.Context
}

// RaftServer implements fmt.Stringer
var _ fmt.Stringer = (*RaftServer)(nil)

// NewRaftServer function to instantiate a RaftServer
func NewRaftServer(member model.Member, config model.Config, database string) *RaftServer {
	var peers []model.Member
	for _, peer := range config.Members {
		if peer.Name == member.Name {
			continue
		}
		peers = append(peers, peer)
	}
	if database == "" {
		database = member.Name + ".db"
	}
	fields := log.Fields{
		"server_name": member.Name,
		"server_port": member.Port,
	}
	ctx := context.WithValue(context.Background(), "fields", fields)
	return &RaftServer{
		member:             member,
		lastHeartbeat:      time.Now(),
		role:               Follower,
		votedOn:            uint64(0),
		databasePath:       database,
		peers:              peers,
		electionCountdown:  config.ElectionCountdown(),
		heartbeatCountdown: config.HeartbeatCountDown(),
		ctx:                ctx,
	}
}

/*
  Management Functions
*/

// Ping the RaftServer
func (s *RaftServer) Ping(_ context.Context, _ *raftapi.Empty) (*raftapi.WhoAmI, error) {
	log.Println("Ping")
	return &raftapi.WhoAmI{
		Name: s.member.Name,
		Port: s.member.Port,
		Role: string(s.role),
	}, nil
}

// Shutdown the RaftServer
func (s *RaftServer) Shutdown(_ context.Context, _ *raftapi.Empty) (*raftapi.Bool, error) {
	log.WithFields(s.LogFields()).Warnln("Shutdown")
	defer func() {
		os.Exit(0)
	}()
	return &raftapi.Bool{Status: true}, nil
}

/*
  Raft Protocol Functions
*/

// RequestVote Raft request to ask peers to participate in a vote.
func (s *RaftServer) RequestVote(_ context.Context, request *raftapi.RequestVoteMessage) (*raftapi.RequestVoteMessage, error) {
	log.WithFields(s.LogFields()).Debugln("Received RequestVote")
	s.lastHeartbeat = time.Now()
	approve := s.votedOn < request.Term
	if approve {
		s.votedOn = request.Term
	}
	return request, nil
}

// AppendEntry Raft request to append a LogEntry to the log.
func (s *RaftServer) AppendEntry(_ context.Context, request *raftapi.AppendEntryRequest) (*raftapi.AppendEntryResponse, error) {
	// TODO handle requests not from leader...?
	log.WithFields(s.LogFields()).Debugln("Received AppendEntry from", request.Leader)
	s.lastHeartbeat = time.Now()
	var term uint64
	if _, err := s.getTerm(); err != nil {
		return nil, err
	}
	if request.Term < term {
		return &raftapi.AppendEntryResponse{Term: term}, nil
	} else if request.Term >= term {
		s.role = Follower
		if err := s.setTerm(request.Term); err != nil {
			return nil, err
		}
	}
	return &raftapi.AppendEntryResponse{Term: term, Success: true}, nil
}

/*
  Other functions
*/

func (s *RaftServer) String() string {
	return fmt.Sprintf("{ name: %s, port: %d, role: %s }", s.member.Name, s.member.Port, s.role)
}

// Run the RaftServer
func (s *RaftServer) Run() error {
	if err := s.setupRepositories(); err != nil {
		return err
	}
	listen, err := net.Listen("tcp", s.member.Address())
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	raftapi.RegisterRaftServiceServer(srv, s)
	go s.monitorHeartbeat()
	go s.produceHeartbeat()
	return srv.Serve(listen)
}

func (s *RaftServer) produceHeartbeat() {
	timeout := s.heartbeatCountdown
	log.WithFields(s.LogFields()).Debugln("heartbeat timeout", timeout)
	for {
		time.Sleep(timeout)
		if s.role == Leader {
			// TODO handle error
			term, _ := s.getTerm()
			s.lastHeartbeat = time.Now()
			for _, member := range s.peers {
				_, _ = member.AppendEntry(s.member.Name, term)
			}
		}
	}
}

func (s *RaftServer) monitorHeartbeat() {
	timeout := s.electionCountdown
	log.WithFields(s.LogFields()).Debugln("election timeout", timeout)
	for {
		time.Sleep(50 * time.Millisecond)
		now := time.Now()
		if now.Sub(s.lastHeartbeat) > timeout {
			log.Debugf("Last Heartbeat: %d, now: %d", s.lastHeartbeat.Unix(), now.Unix())
			log.Debugln("Delta: ", now.Sub(s.lastHeartbeat))
			s.role = Candidate
			if s.runElection() {
				s.lastHeartbeat = time.Now()
				log.WithFields(s.LogFields()).Infoln("Role now", Leader)
				s.role = Leader
			}
		}
	}
}

func (s *RaftServer) runElection() bool {
	log.WithFields(s.LogFields()).Infoln("kicking off vote")
	term, _ := s.getTerm()
	term++
	_ = s.setTerm(term)
	s.votedOn = term
	var votes = 1
	for _, member := range s.peers {
		resp, err := member.RequestVote(s.ctx, s.member.Name, term)
		if err != nil {
			log.Errorf("%s: No response from %s\n", s.String(), member.String())
			continue
		}
		if resp.Term > term {
			log.Errorf("%s: Response %v indicate election term conflict", s.String(), resp)
			_ = s.setTerm(resp.Term)
			s.role = Follower
			return false
		}
		votes++
	}
	return votes > (len(s.peers)+1)/2
}

func (s *RaftServer) setupRepositories() error {
	db, err := gorm.Open(sqlite.Open(s.databasePath), &gorm.Config{})
	if err != nil {
		return err
	}
	sRepo, err := database.NewStatusRepository(db)
	if err != nil {
		return err
	}
	if err = sRepo.Migrate(); err != nil {
		return err
	}
	lRepo, err := database.NewLogEntryRepository(db)
	if err != nil {
		return err
	}
	if err = lRepo.Migrate(); err != nil {
		return err
	}
	s.statusRepo = sRepo
	s.logRepo = lRepo
	return nil
}

func (s *RaftServer) getTerm() (uint64, error) {
	status, err := s.statusRepo.Read(s.member.Name)
	if err != nil {
		return 0, err
	}
	if status == nil {
		return 0, nil
	}
	return status.Term, nil
}

func (s *RaftServer) setTerm(term uint64) error {
	status := model.Status{
		Name: s.member.Name,
		Term: term,
	}
	return s.statusRepo.Write(&status)
}

func (s *RaftServer) LogFields() log.Fields {
	return log.Fields{"server_name": s.member.Name, "server_port": s.member.Port}
}
