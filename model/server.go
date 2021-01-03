package model

import (
	"context"
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
	"os"
	"time"
)

type Role string

const (
	Candidate Role = "CANDIDATE"
	Leader    Role = "LEADER"
	Follower  Role = "FOLLOWER"
)

type Server struct {
	raftapi.UnimplementedRaftServiceServer
	member             Member
	lastHeartbeat      time.Time
	role               Role
	peers              []Member
	votedOn            uint64
	log                *log.Entry
	databasePath       string
	db                 *gorm.DB
	electionCountdown  time.Duration
	heartbeatCountdown time.Duration
}

// Server implements fmt.Stringer
var _ fmt.Stringer = (*Server)(nil)

func NewServer(member Member, config Config, database string) *Server {
	var logger = log.New()
	logger.Out = os.Stdout
	logger.Level = log.DebugLevel
	entry := logger.WithFields(log.Fields{
		"server": member.Name,
	})
	if database == "" {
		database = member.Name + ".db"
	}
	var peers []Member
	for _, peer := range config.Members {
		if peer.Name == member.Name {
			continue
		}
		peers = append(peers, peer)
	}
	return &Server{
		member: member,
		// config:        config,
		lastHeartbeat:      time.Now(),
		role:               Follower,
		votedOn:            uint64(0),
		log:                entry,
		databasePath:       database,
		peers:              peers,
		electionCountdown:  config.ElectionCountdown(),
		heartbeatCountdown: config.HeartbeatCountDown(),
	}
}

/*
  Management Functions
*/

func (s *Server) Ping(_ context.Context, _ *raftapi.Empty) (*raftapi.WhoAmI, error) {
	s.log.Println("Ping")
	return &raftapi.WhoAmI{
		Name: s.member.Name,
		Port: s.member.Port,
		Role: string(s.role),
	}, nil
}

func (s *Server) Shutdown(_ context.Context, _ *raftapi.Empty) (*raftapi.Bool, error) {
	s.log.Warnln("Shutdown")
	defer func() {
		os.Exit(0)
	}()
	return &raftapi.Bool{Status: true}, nil
}

/*
  Raft Protocol Functions
*/

func (s *Server) RequestVote(_ context.Context, request *raftapi.RequestVoteMessage) (*raftapi.RequestVoteMessage, error) {
	s.log.Debugln("Received RequestVote")
	s.lastHeartbeat = time.Now()
	approve := s.votedOn < request.Term
	if approve {
		s.votedOn = request.Term
	}
	return request, nil
}

func (s *Server) AppendEntry(_ context.Context, request *raftapi.AppendEntryRequest) (*raftapi.AppendEntryResponse, error) {
	// TODO handle requests not from leader...?
	s.log.Debugln("Received AppendEntry from", request.Leader)
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

func (s *Server) String() string {
	return fmt.Sprintf("{ name: %s, port: %d, role: %s }", s.member.Name, s.member.Port, s.role)
}

func (s *Server) Run() error {
	if err := s.setupDB(); err != nil {
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
	s.log.Infoln("Starting:", s)
	return srv.Serve(listen)
}

func (s *Server) produceHeartbeat() {
	timeout := s.heartbeatCountdown
	s.log.Debugln("heartbeat timeout", timeout)
	for {
		time.Sleep(timeout)
		if s.role == Leader {
			s.lastHeartbeat = time.Now()
			for _, member := range s.peers {
				_, _ = member.AppendEntry(s)
			}
		}
	}
}

func (s *Server) monitorHeartbeat() {
	timeout := s.electionCountdown
	s.log.Debugln("election timeout", timeout)
	for {
		time.Sleep(50 * time.Millisecond)
		now := time.Now()
		if now.Sub(s.lastHeartbeat) > timeout {
			s.log.Debugf("Last Heartbeat: %d, now: %d", s.lastHeartbeat.Unix(), now.Unix())
			s.log.Debugln("Delta: ", now.Sub(s.lastHeartbeat))
			s.role = Candidate
			if s.runElection() {
				s.lastHeartbeat = time.Now()
				s.log.Infoln("Role now", Leader)
				s.role = Leader
			}
		}
	}
}

func (s *Server) runElection() bool {
	s.log.Infoln("kicking off vote")
	term, _ := s.getTerm()
	term += 1
	_ = s.setTerm(term)
	s.votedOn = term
	var votes = 1
	for _, member := range s.peers {
		resp, err := member.RequestVote(s)
		if err != nil {
			s.log.Errorf("%s: No response from %s\n", s.String(), member.String())
			continue
		}
		if resp.Term > term {
			s.log.Errorf("%s: Response %v indicate election term conflict", s.String(), resp)
			_ = s.setTerm(resp.Term)
			s.role = Follower
			return false
		}
		votes += 1
	}
	return votes > (len(s.peers)+1)/2
}

func (s *Server) setupDB() error {
	db, err := gorm.Open(sqlite.Open(s.databasePath), &gorm.Config{})
	if err != nil {
		return err
	}
	s.db = db
	if err := db.AutoMigrate(&Status{}); err != nil {
		return err
	}
	return nil
}

func (s *Server) getTerm() (uint64, error) {
	var statuses []Status
	tx := s.db.Where("name = ?", s.member.Name).Find(&statuses)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if len(statuses) == 0 {
		return 0, nil
	}
	return statuses[0].Term, nil
}

func (s *Server) setTerm(term uint64) error {
	status := Status{
		Name: s.member.Name,
		Term: term,
	}
	tx := s.db.Model(&status).Update("term", status.Term)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		tx2 := s.db.Create(&status)
		if tx2.Error != nil {
			return tx2.Error
		}
	}
	return nil
}
