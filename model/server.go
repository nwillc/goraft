package model

import (
	"context"
	"github.com/nwillc/goraft/api/raftapi"
	"log"
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
