package model

import (
	"context"
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Member struct {
	Name string `json:"name"`
	Port uint32 `json:"port"`
}

// Member implements fmt.Stringer
var _ fmt.Stringer = (*Member)(nil)

func (m *Member) String() string {
	return fmt.Sprintf("{ name: %s, port: %d }", m.Name, m.Port)
}

func (m *Member) Address() string {
	return fmt.Sprintf(":%d", m.Port)
}

func (m *Member) AppendEntry(s *Server) (uint64, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	ctx := context.Background()
	term, err := s.getTerm()
	if err != nil {
		return 0, err
	}
	response, err := api.AppendEntry(ctx, &raftapi.AppendEntryRequest{
		Term:        term,
		Leader:      s.member.Name,
		LogSize:     0,
		LastLogTerm: 0,
	})
	if err != nil {
		return 0, err
	}

	return response.Term, nil
}

func (m *Member) RequestVote(s *Server) (*raftapi.RequestVoteMessage, error) {
	s.log.WithFields(log.Fields{"member": m.Name}).Debugln("Requesting vote from")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	ctx := context.Background()
	term, err := s.getTerm()
	if err != nil {
		return nil, err
	}
	response, err := api.RequestVote(ctx, &raftapi.RequestVoteMessage{
		Term:        term,
		Candidate:   s.member.Name,
		LogSize:     0,
		LastLogTerm: 0,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
