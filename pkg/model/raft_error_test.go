package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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

func (suite RaftErrorTestSuite) TestSingleError() {
	err := fmt.Errorf("foo")
	raftError := NewRaftError(suite.member, err)
	assert.NotNil(suite.T(), raftError)
	assert.Equal(suite.T(), fmt.Sprintf("member: %s error_0: %s", suite.member.Name, err.Error()), raftError.Error())
}

func (suite RaftErrorTestSuite) TestDoubleError() {
	err1 := fmt.Errorf("foo")
	err2 := fmt.Errorf("bar")
	raftError := NewRaftError(suite.member, err1, err2)
	assert.NotNil(suite.T(), raftError)
	assert.Equal(suite.T(), fmt.Sprintf("member: %s error_0: %s error_1: %s",
		suite.member.Name, err1.Error(), err2.Error()), raftError.Error())
}
