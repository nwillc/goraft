package model

import "fmt"

// RaftError is an error noting the Raft Member
type RaftError struct {
	Member *Member
	Err    error
}

// RaftError implements error
var _ error = RaftError{}

func (err RaftError) Error() string {
	return fmt.Sprintf("member %s: error %v", err.Member.Name, err.Err)
}
