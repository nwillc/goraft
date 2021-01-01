package model

import (
	"context"
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	member        Member
	config        Config
	lastHeartbeat time.Time
	term          uint32
	role          Role
	votedOn       uint32
	log           *log.Entry
}

// Server implements fmt.Stringer
var _ fmt.Stringer = (*Server)(nil)

func NewServer(member Member, config Config, offset int) *Server {
	var logger = log.New()
	logger.Out = os.Stdout
	logger.Level = log.DebugLevel
	entry := logger.WithFields(log.Fields{
		"server": member.Name,
	})
	return &Server{
		member:        member,
		config:        config,
		lastHeartbeat: time.Now(),
		role:          Follower,
		votedOn:       uint32(0),
		log:           entry,
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

func (s *Server) AppendEntry(_ context.Context, request *raftapi.AppendEntryRequest) (*raftapi.Bool, error) {
	s.log.Debugln("Received AppendEntry")
	s.lastHeartbeat = time.Now()
	s.term = request.Term
	return &raftapi.Bool{Status: true}, nil
}

/*
  Other functions
*/

func (s *Server) String() string {
	return fmt.Sprintf("{ name: %s, port: %d, role: %s }", s.member.Name, s.member.Port, s.role)
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", s.member.Address())
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	raftapi.RegisterRaftServiceServer(srv, s)
	go s.monitorHeartbeat()
	s.log.Infoln("Starting:", s)
	return srv.Serve(listen)
}

func (s *Server) monitorHeartbeat() {
	timeout := s.config.ElectionCountdown()
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
	s.term += 1
	s.votedOn = s.term
	var votes = 1
	for _, member := range s.config.Members {
		if member.Name == s.member.Name {
			continue
		}
		resp, err := member.RequestVote(s)
		if err != nil {
			s.log.Errorf("%s: No response from %s\n", s.String(), member.String())
			continue
		}
		if resp.Term > s.term {
			s.log.Errorf("%s: Response %v indicate election term conflict", s.String(), resp)
			s.term = resp.Term
			s.role = Follower
			return false
		}
		votes += 1
	}
	return votes > len(s.config.Members)/2
}
