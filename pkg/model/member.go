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

func (m *Member) AppendEntry(leader string, term uint64) (uint64, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	ctx := context.Background()
	response, err := api.AppendEntry(ctx, &raftapi.AppendEntryRequest{
		Term:        term,
		Leader:      leader,
		LogSize:     0,
		LastLogTerm: 0,
	})
	if err != nil {
		return 0, err
	}

	return response.Term, nil
}

func (m *Member) RequestVote(logger *log.Entry, leader string, term uint64) (*raftapi.RequestVoteMessage, error) {
	logger.WithFields(log.Fields{"member": m.Name}).Debugln("Requesting vote from")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	ctx := context.Background()
	response, err := api.RequestVote(ctx, &raftapi.RequestVoteMessage{
		Term:        term,
		Candidate:   leader,
		LogSize:     0,
		LastLogTerm: 0,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
