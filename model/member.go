package model

import (
	"fmt"
	"github.com/nwillc/goraft/api/raftapi"
	"google.golang.org/grpc"
	"log"
	"context"
)

type Member struct {
	Name string `json:"name"`
	Port uint32 `json:"port"`
}

// Member implements fmt.Stringer
var _ fmt.Stringer = (* Member)(nil)

func (m *Member) String() string {
	return fmt.Sprintf("{ name: %s, port: %d }", m.Name, m.Port)
}

func (m *Member) Address() string {
	return fmt.Sprintf(":%d", m.Port)
}

func (m *Member) RequestVote(s *Server) (bool, error){
	log.Println("Requesting vote:", m.String())

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	ctx := context.Background()
	response, err := api.RequestVote(ctx, &raftapi.RequestVoteRequest{
		Term:        s.term + 1,
		Candidate:   s.member.Name,
		LogSize:     0,
		LastLogTerm: 0,
	})
	if err != nil {
		return false, err
	}

	return response.Status, nil
}
