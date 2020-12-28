package model

import (
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	"google.golang.org/grpc"
	"net"
)

type Member struct {
	Port uint   `json:"port"`
}

func (m *Member) Address() string {
	return fmt.Sprintf(":%d", m.Port)
}

func (m *Member) Listen() error {
	listen, err := net.Listen("tcp", m.Address())
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	raftapi.RegisterRaftServiceServer(srv, &server{})
	return srv.Serve(listen)
}
