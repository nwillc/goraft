package model

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RaftErrorTestSuite struct {
	suite.Suite
	member *Member
}

func (suite *RaftErrorTestSuite) SetupTest() {
	suite.T().Helper()
	suite.member = &Member{
		Name: "foo",
		Port: 100,
	}
}

func TestRaffTestSuite(t *testing.T) {
	suite.Run(t, new(RaftErrorTestSuite))
}

func (suite RaftErrorTestSuite) Test_SingleError() {
	err := fmt.Errorf("foo")
	raftError := NewRaftError(suite.member, err)
	suite.NotNil(raftError)
	suite.Equal(fmt.Sprintf("member: %s error_0: %s", suite.member.Name, err.Error()), raftError.Error())
}

func (suite RaftErrorTestSuite) Test_DoubleError() {
	err1 := fmt.Errorf("foo")
	err2 := fmt.Errorf("bar")
	raftError := NewRaftError(suite.member, err1, err2)
	suite.NotNil(raftError)
	suite.Equal(fmt.Sprintf("member: %s error_0: %s error_1: %s",
		suite.member.Name, err1.Error(), err2.Error()), raftError.Error())
}
