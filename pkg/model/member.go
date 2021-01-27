package model

import (
	"context"
	"fmt"
	"github.com/nwillc/goraft/raftapi"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Member is a Raft cluster member.
type Member struct {
	Name string `json:"name"`
	Port uint32 `json:"port"`
}

// Member implements fmt.Stringer
var _ fmt.Stringer = (*Member)(nil)

func (m *Member) String() string {
	return fmt.Sprintf("{ name: %s, port: %d }", m.Name, m.Port)
}

// Address off the Member.
func (m *Member) Address() string {
	return fmt.Sprintf(":%d", m.Port)
}

func (m *Member) Ping() error {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
	if err != nil {
		return RaftError{Member: m, Err: err}
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	ctx := context.Background()
	_, err = api.Ping(ctx, &raftapi.Empty{})
	if err != nil {
		return RaftError{Member: m, Err: err}
	}
	return nil
}

// AppendEntry request of a Member.
func (m *Member) AppendEntry(leader string, term uint64, value int64, prevLogId int64) (bool, error) {
	log.Infoln("Requesting log entry of", m.Name, "Value", value)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
	if err != nil {
		return false, RaftError{Member: m, Err: err}
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	ctx := context.Background()
	ee := raftapi.AppendEntryRequest_Entry{
		Entry: &raftapi.LogEntry{
			Term:  term,
			Value: value,
		},
	}
	response, err := api.AppendEntry(ctx, &raftapi.AppendEntryRequest{
		Term:        term,
		Leader:      leader,
		PrevLogId:   prevLogId,
		PrevLogTerm: term,
		LogEntry:    &ee,
	})
	if err != nil {
		return false, RaftError{Member: m, Err: err}
	}

	return response.Success, nil
}

// RequestVote request of a Member.
func (m *Member) RequestVote(ctx context.Context, leader string, term uint64, logSize uint64) (*raftapi.RequestVoteMessage, error) {
	log.WithFields(log.Fields{"member_name": m.Name}).Debugln("Requesting vote from")
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(m.Address(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	api := raftapi.NewRaftServiceClient(conn)
	response, err := api.RequestVote(ctx, &raftapi.RequestVoteMessage{
		Term:        term,
		Candidate:   leader,
		LogSize:     logSize,
		LastLogTerm: 0,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
