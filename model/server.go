package model

import (
	"context"
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

type Role string

const (
	Leader   Role = "LEADER"
	Follower Role = "FOLLOWER"
)

type Server struct {
	raftapi.UnimplementedRaftServiceServer
	member        Member
	config        Config
	lastHeartbeat time.Time
	term          uint32
	role          Role
}

// Server implements fmt.Stringer
var _ fmt.Stringer = (*  Server)(nil)

func NewServer(member Member, config Config) *Server {
	return &Server{
		member:        member,
		config:        config,
		lastHeartbeat: time.Now(),
		role:          Follower,
	}
}

/*
  Management Functions
*/

func (s *Server) Ping(ctx context.Context, _ *raftapi.Empty) (*raftapi.WhoAmI, error) {
	log.Println("Ping")
	return &raftapi.WhoAmI{
		Name: s.member.Name,
		Port: s.member.Port,
	}, nil
}

func (s *Server) Shutdown(ctx context.Context, _ *raftapi.Empty) (*raftapi.Bool, error) {
	log.Println("Shutdown")
	defer func() {
		os.Exit(0)
	}()
	return &raftapi.Bool{Status: true}, nil
}

/*
  Raft Protocol Functions
*/

func (s *Server) RequestVote(ctx context.Context, request *raftapi.RequestVoteRequest) (*raftapi.Bool, error) {
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
	log.Println("Starting:", s)
	return srv.Serve(listen)
}

func (s *Server) monitorHeartbeat() {
	for ; ; {
		now := time.Now()
		if now.Sub(s.lastHeartbeat) > (time.Duration(s.config.Election) * time.Millisecond) {
			log.Println("Election")
			s.runElection()
			s.lastHeartbeat = now
		}
	}
}

func (s *Server) runElection() {
	s.term += 1
	for _, member := range s.config.Members {
		if s.member.Name == member.Name {
			continue
		}
		response, err := member.RequestVote(s)
		log.Println("Response:", response, "Error:", err)
	}
}
