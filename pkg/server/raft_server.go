package server

import (
	"context"
	"fmt"
	"github.com/nwillc/goraft/database"
	"github.com/nwillc/goraft/model"
	"github.com/nwillc/goraft/raftapi"
	"github.com/nwillc/goraft/util"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
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

type State string

const (
	Ready    State = "READY"
	Running  State = "RUNNING"
	Shutdown State = "SHUTDOWN"
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
	leaderID           string
	onExit             *util.FunctionChain
	state              State
}

// RaftServer implements fmt.Stringer
var _ fmt.Stringer = (*RaftServer)(nil)

// RaftServer implements raftapi.RaftServiceServer
var _ raftapi.RaftServiceServer = (*RaftServer)(nil)

// NewRaftServer function to instantiate a RaftServer
func NewRaftServer(member model.Member, config model.Config, database string) *RaftServer {
	if database == "" {
		database = member.Name + ".db"
	}
	rf := RaftServer{
		member:             member,
		lastHeartbeat:      time.Now(),
		role:               Follower,
		votedOn:            uint64(0),
		databasePath:       database,
		peers:              config.Peers(member.Name),
		electionCountdown:  config.ElectionCountdown(),
		heartbeatCountdown: config.HeartbeatCountDown(),
		ctx:                context.Background(),
		state:              Ready,
		onExit:             &util.FunctionChain{},
	}
	rf.onExit.Add(func() {
		rf.state = Shutdown
		log.WithFields(rf.LogFields()).Infoln("State set to", rf.state)
	})
	return &rf
}

/*
  Management Functions
*/

// Ping the RaftServer
func (s *RaftServer) Ping(_ context.Context, _ *raftapi.Empty) (*raftapi.WhoAmI, error) {
	log.WithFields(s.LogFields()).Debugln("Received Ping")
	s.lastHeartbeat = time.Now()
	return &raftapi.WhoAmI{
		Name: s.member.Name,
		Port: s.member.Port,
		Role: string(s.role),
	}, nil
}

// Shutdown the RaftServer
func (s *RaftServer) Shutdown(_ context.Context, _ *raftapi.Empty) (*raftapi.Bool, error) {
	log.WithFields(s.LogFields()).Warnln("Shutdown")
	s.onExit.Add(func() {
		s.state = Shutdown
	})
	s.onExit.InvokeReverse()
	return &raftapi.Bool{Status: true}, nil
}

// AppendValue requests the leader to append a log value to the raft log
func (s *RaftServer) AppendValue(_ context.Context, value *raftapi.Value) (*raftapi.Bool, error) {
	log.WithFields(s.LogFields()).Infoln("received request to append", value.Value)
	if s.role != Leader {
		msg := "request to append log to non leader"
		log.WithFields(s.LogFields()).Errorf(msg)
		return nil, model.NewRaftError(&s.member, fmt.Errorf(msg))
	}
	term, _ := s.getTerm()
	lastEntry, err := s.logRepo.LastEntry()
	if err != nil {
		lastEntry = &model.LogEntry{
			EntryNo: -1,
			Term:    0,
			Value:   0,
		}
	}
	// My entry
	_, err = s.logRepo.Create(term, value.Value)
	if err != nil {
		return nil, model.NewRaftError(&s.member, err)
	}
	// TODO: Retry
	var succeeded = 0.0
	for _, member := range s.peers {
		success, err := member.AppendEntry(s.member.Name, term, value.Value, lastEntry.EntryNo)
		if err != nil {
			log.Errorln(err)
			continue
		}
		if success {
			succeeded++
		}
	}

	quorum := succeeded >= float64(len(s.peers))/2.0
	if !quorum {
		_ = s.logRepo.TruncateToEntryNo(lastEntry.EntryNo)
	}
	return &raftapi.Bool{Status: quorum}, nil
}

// ListEntries returns the RaftServer's persisted log entries.
func (s *RaftServer) ListEntries(_ context.Context, _ *raftapi.Empty) (*raftapi.EntryListResponse, error) {
	list, err := s.logRepo.List()
	if err != nil {
		return nil, model.NewRaftError(&s.member, err)
	}
	var logEntries = make([]*raftapi.LogEntry, 0)
	for _, entry := range list {
		logEntries = append(logEntries, &raftapi.LogEntry{
			Term:  entry.Term,
			Value: entry.Value,
		})
	}
	response := &raftapi.EntryListResponse{Entries: logEntries}
	return response, nil
}

/*
  Raft Protocol Functions
*/

// RequestVote Raft request to ask peers to participate in a vote.
func (s *RaftServer) RequestVote(_ context.Context, request *raftapi.RequestVoteMessage) (*raftapi.RequestVoteResponse, error) {
	log.WithFields(s.LogFields()).Debugln("Received RequestVote")
	s.lastHeartbeat = time.Now()
	term, err := s.getTerm()
	if err  != nil {
		log.WithFields(s.LogFields()).Errorln("Could not look up term", err)
		return nil, model.NewRaftError(&s.member, err)
	}
	approve := s.votedOn < request.Term
	if approve {
		s.votedOn = request.Term
	}
	return &raftapi.RequestVoteResponse{
		Term:     term,
		Approved: approve,
	}, nil
}

// AppendEntry Raft request to append a LogEntry to the log.
func (s *RaftServer) AppendEntry(_ context.Context, request *raftapi.AppendEntryRequest) (*raftapi.AppendEntryResponse, error) {
	// TODO handle requests not from leader...?
	log.WithFields(s.LogFields()).Debugln("Received AppendEntry from", request.Leader)
	s.lastHeartbeat = time.Now()
	term, err := s.getTerm()
	if err != nil {
		log.WithFields(s.LogFields()).Errorln("Could not look up term", err)
		return nil, model.NewRaftError(&s.member, err)
	}
	if request.Term < term {
		log.WithFields(s.LogFields()).Warnln("Term", request.Term, "Less than my term", term)
		return &raftapi.AppendEntryResponse{Term: term}, nil
	} else if request.Term >= term {
		s.role = Follower
		if err := s.setTerm(request.Term); err != nil {
			log.WithFields(s.LogFields()).Errorln("Unable to update my term")
			return nil, model.NewRaftError(&s.member, err)
		}
	}
	s.leaderID = request.Leader
	maxID, err := s.logRepo.MaxEntryNo()
	if err != nil {
		log.WithFields(s.LogFields()).Errorln("Could not look up MaxEntryNo", err)
		return nil, model.NewRaftError(&s.member, err)
	}
	if request.PrevLogId == -1 {
		switch entry := request.LogEntry.(type) {
		case *raftapi.AppendEntryRequest_Entry:
			_, err = s.logRepo.Create(entry.Entry.Term, entry.Entry.Value)
			if err != nil {
				return nil, model.NewRaftError(&s.member, err)
			}
			return &raftapi.AppendEntryResponse{
				Term:    term,
				Success: true,
			}, nil
		}
	} else if request.PrevLogId <= maxID {
		entry, _ := s.logRepo.Read(request.PrevLogId)
		if entry.Term != request.PrevLogTerm {
			log.WithFields(s.LogFields()).Warnln("Entry term", entry.Term, "not equal to previous term", request.PrevLogTerm)
			return &raftapi.AppendEntryResponse{Term: term}, nil
		}
		switch x := request.LogEntry.(type) {
		case *raftapi.AppendEntryRequest_Entry:
			log.Infoln("Appending Term", x.Entry.Term, "Value", x.Entry.Value)
			_, _ = s.logRepo.Create(x.Entry.Term, x.Entry.Value)
		case nil:
			// No entry
			log.WithFields(s.LogFields()).Warnln("No entry")
		default:
			log.WithFields(s.LogFields()).Errorf("Unknown request log entry type %T", x)
			return nil, model.NewRaftError(&s.member, err)
		}
		log.WithFields(s.LogFields()).Infoln("Returning success")
		return &raftapi.AppendEntryResponse{Term: term, Success: true}, nil
	} else {
		log.WithFields(s.LogFields()).Warnf("Previous ID: %d, MaxID: %d, returning nil", request.PrevLogId, maxID)
		return &raftapi.AppendEntryResponse{Term: term}, nil
	}
	return nil, model.NewRaftError(&s.member, fmt.Errorf("unknown code path"))
}

/*
  Other functions
*/

func (s *RaftServer) String() string {
	return fmt.Sprintf("{ name: %s, port: %d, role: %s }", s.member.Name, s.member.Port, s.role)
}

// Run the RaftServer
func (s *RaftServer) Run() error {
	log.WithFields(s.LogFields()).Infoln("Running")
	if err := s.setupRepositories(); err != nil {
		return err
	}
	listen, err := net.Listen("tcp", s.member.Address())
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	s.onExit = s.onExit.Add(func() {
		log.WithFields(s.LogFields()).Infoln("Stopping grpc")
		server.GracefulStop()
		s.state = Shutdown
	})
	raftapi.RegisterRaftServiceServer(server, s)
	go s.monitorHeartbeat()
	go s.produceHeartbeat()
	s.state = Running
	return server.Serve(listen)
}

func (s *RaftServer) produceHeartbeat() {
	timeout := s.heartbeatCountdown
	log.WithFields(s.LogFields()).Debugln("heartbeat timeout", timeout)
	for {
		time.Sleep(timeout)
		if s.state == Shutdown {
			break
		}
		if s.role == Leader {
			s.lastHeartbeat = time.Now()
			for _, member := range s.peers {
				_ = member.Ping()
			}
		}
	}
	log.WithFields(s.LogFields()).Infoln("No longer producing heartbeats")
}

func (s *RaftServer) monitorHeartbeat() {
	timeout := s.electionCountdown
	log.WithFields(s.LogFields()).Debugln("election timeout", timeout)
	for {
		time.Sleep(50 * time.Millisecond)
		if s.state == Shutdown {
			break
		}
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
	log.WithFields(s.LogFields()).Infoln("No longer monitoring heartbeats")
}

func (s *RaftServer) runElection() bool {
	log.WithFields(s.LogFields()).Infoln("kicking off vote")
	term, _ := s.getTerm()
	term++
	_ = s.setTerm(term)
	logSize, _ := s.logRepo.MaxTerm()
	s.votedOn = term
	var votes = 1
	for _, member := range s.peers {
		resp, err := member.RequestVote(s.ctx, s.member.Name, term, logSize)
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
	s.onExit = s.onExit.Add(func() {
		log.WithFields(s.LogFields()).Infoln("Closing database")
		sqlDb, _ := db.DB()
		_ = sqlDb.Close()
	})
	return nil
}

func (s *RaftServer) getTerm() (uint64, error) {
	status, err := s.statusRepo.Read(s.member.Name)
	if err != nil || status == nil {
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

// LogFields creates a log.Fields with this RaftServer's name and port info
func (s *RaftServer) LogFields() log.Fields {
	return log.Fields{"server_name": s.member.Name, "server_port": s.member.Port}
}

// GetRole returns the RaftServer's Role.
func (s *RaftServer) GetRole() Role {
	return s.role
}

// GetState returns the RaftServer's State.
func (s *RaftServer) GetState() State {
	return s.state
}
