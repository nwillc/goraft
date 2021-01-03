package model

import (
	"github.com/nwillc/goraft/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

type ServerTestSuite struct {
	suite.Suite
	server *Server
}

func (suite *ServerTestSuite) SetupTest() {
	suite.T().Helper()

	conf, err := ReadConfig("../" + conf.ConfigFile)
	assert.NoError(suite.T(), err)
	member := conf.Members[0]
	tempFile, err := ioutil.TempFile("", member.Name+"*.db")
	assert.NoError(suite.T(), err)
	suite.T().Cleanup(func() {
		_ = os.Remove(tempFile.Name())
	})
	suite.server = NewServer(member, conf, tempFile.Name())
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) TestServerSanity() {
	assert.NotEmpty(suite.T(), suite.server.member.Name)
	assert.Less(suite.T(), 0, len(suite.server.peers))
	assert.Nil(suite.T(), suite.server.db)
	assert.NoError(suite.T(), suite.server.setupDB())
	assert.NotNil(suite.T(), suite.server.db)
}

func (suite *ServerTestSuite) TestPersistTerm() {
	term := uint64(23)
	assert.NoError(suite.T(), suite.server.setupDB())
	assert.NoError(suite.T(), suite.server.setTerm(term))
	t, err := suite.server.getTerm()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), term, t)
}
