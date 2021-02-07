package model

import (
	"fmt"
)

// RaftError is an error noting the Raft Member
type RaftError struct {
	Member *Member
	Errors []error
}

// RaftError implements error
var _ error = RaftError{}

// NewRaftError new RaftError
func NewRaftError(member *Member, errors ...error) RaftError {
	return RaftError{
		Member: member,
		Errors: errors,
	}
}

func (er RaftError) Error() string {
	msg := fmt.Sprintf("member: %s", er.Member.Name)
	for i, err := range er.Errors {
		msg += fmt.Sprintf(" error_%d: %s", i, err.Error())
	}
	return msg
}
