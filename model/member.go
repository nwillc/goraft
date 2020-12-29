package model

import (
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	"google.golang.org/grpc"
	"net"
)

type Member struct {
	Port uint32   `json:"port"`
}

func (m *Member) Address() string {
	return fmt.Sprintf(":%d", m.Port)
}

func (m *Member) Listen(name string) error {
	listen, err := net.Listen("tcp", m.Address())
	if err != nil {
		return err
	}
	s := Server{
		Member: Member{
			Port: m.Port,
		},
		Name: name,
	}
	srv := grpc.NewServer()
	raftapi.RegisterRaftServiceServer(srv, &s)
	return srv.Serve(listen)
}
