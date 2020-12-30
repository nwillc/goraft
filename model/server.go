package model

import (
	"context"
	"github.com/nwillc/goraft/api/raftapi"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type Server struct {
	raftapi.UnimplementedRaftServiceServer
	Member
}

func NewServer(member Member) *Server {
	return &Server{
		Member: member,
	}
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", s.Address())
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	raftapi.RegisterRaftServiceServer(srv, s)
	return srv.Serve(listen)
}

func (s *Server) Ping(ctx context.Context, request *raftapi.Empty) (*raftapi.WhoAmI, error) {
	log.Println("Ping")
	return &raftapi.WhoAmI{
		Name: s.Name,
		Port: s.Port,
	}, nil
}

func (s *Server) Shutdown(ctx context.Context, empty *raftapi.Empty) (*raftapi.Bool, error)  {
	log.Println("Shutdown")
	defer func() {
		os.Exit(0)
	}()
	return &raftapi.Bool{Status: true}, nil
}
