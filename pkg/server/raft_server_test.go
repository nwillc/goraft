package server

import (
	"github.com/nwillc/goraft/conf"
	"github.com/nwillc/goraft/model"
	log "github.com/sirupsen/logrus"
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

	config, err := model.ReadConfig("../../" + conf.ConfigFile)
	suite.NoError(err)
	member := config.Members[0]
	tempFile, err := ioutil.TempFile("", member.Name+"*.db")
	suite.NoError(err)
	suite.T().Cleanup(func() {
		_ = os.Remove(tempFile.Name())
	})
	log.SetLevel(log.WarnLevel)
	suite.server = NewRaftServer(member, config, tempFile.Name())
}

func TestRaftServerTestSuite(t *testing.T) {
	suite.Run(t, new(RaftServerTestSuite))
}

func (suite *RaftServerTestSuite) Test_RaftServerSanity() {
	suite.NotEmpty(suite.server.member.Name)
	suite.Less(0, len(suite.server.peers))
	suite.Nil(suite.server.statusRepo)
	suite.NoError(suite.server.setupRepositories())
	suite.NotNil(suite.server.statusRepo)
	term := suite.server.getTerm()
	suite.Equal(uint64(0), term)
}

func (suite *RaftServerTestSuite) Test_PersistTerm() {
	term := uint64(23)
	suite.NoError(suite.server.setupRepositories())
	suite.NoError(suite.server.setTerm(term))
	t := suite.server.getTerm()
	suite.Equal(term, t)
}
