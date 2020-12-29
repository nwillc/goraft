package model

import (
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	"google.golang.org/grpc"
	"net"
)

type Member struct {
	Name string   `json:"name"`
	Port uint32   `json:"port"`
}

var _ fmt.Stringer = (* Member)(nil)

func (m *Member) String() string {
	return fmt.Sprintf("{ name: %s, port: %d }", m.Name, m.Port)
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
