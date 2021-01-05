package server

import (
	"github.com/nwillc/goraft/conf"
	"github.com/nwillc/goraft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

type RaftServerTestSuite struct {
	suite.Suite
	server *RaftServer
}

func (suite *RaftServerTestSuite) SetupTest() {
	suite.T().Helper()

	conf, err := model.ReadConfig("../" + conf.ConfigFile)
	assert.NoError(suite.T(), err)
	member := conf.Members[0]
	tempFile, err := ioutil.TempFile("", member.Name+"*.db")
	assert.NoError(suite.T(), err)
	suite.T().Cleanup(func() {
		_ = os.Remove(tempFile.Name())
	})
	suite.server = NewRaftServer(member, conf, tempFile.Name())
}

func TestRaftServerTestSuite(t *testing.T) {
	suite.Run(t, new(RaftServerTestSuite))
}

func (suite *RaftServerTestSuite) TestRaftServerSanity() {
	assert.NotEmpty(suite.T(), suite.server.member.Name)
	assert.Less(suite.T(), 0, len(suite.server.peers))
	assert.Nil(suite.T(), suite.server.statusRepo)
	assert.NoError(suite.T(), suite.server.setupRepositories())
	assert.NotNil(suite.T(), suite.server.statusRepo)
}

func (suite *RaftServerTestSuite) TestPersistTerm() {
	term := uint64(23)
	assert.NoError(suite.T(), suite.server.setupRepositories())
	assert.NoError(suite.T(), suite.server.setTerm(term))
	t, err := suite.server.getTerm()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), term, t)
}
