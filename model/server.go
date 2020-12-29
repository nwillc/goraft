package model

import (
	"context"
	"github.com/nwillc/goraft/api/raftapi"
	"log"
	"os"
)

type Server struct {
	raftapi.UnimplementedRaftServiceServer
	Member
	Name string
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
