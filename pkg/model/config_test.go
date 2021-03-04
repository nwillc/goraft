package model

import (
	"github.com/nwillc/goraft/conf"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
	config Config
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.T().Helper()
	config, err := ReadConfig("../../" + conf.ConfigFile)
	suite.Require().NoError(err)
	suite.config = config
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) Test_FailBadFile() {
	_, err := ReadConfig("foo.json")
	suite.Error(err)
}

func (suite *ConfigTestSuite) Test_ReadConfig() {
	suite.Less(1, suite.config.HeartbeatTimeout)
	suite.Less(1, suite.config.ElectionTimeout)
	suite.Less(0, suite.config.MinOffset)
	suite.Less(suite.config.MinOffset, suite.config.MaxOffset)
	for _, member := range suite.config.Members {
		suite.NotEmpty(member.Name)
		suite.Less(uint32(0), member.Port)
	}
}

func (suite *ConfigTestSuite) Test_ConfigPeers() {
	suite.NotEmpty(suite.config.Members)
	aMember := suite.config.Members[0].Name
	peers := suite.config.Peers(aMember)
	suite.Equal(len(suite.config.Members)-1, len(peers))
	var found bool
	for _, member := range peers {
		if member.Name == aMember {
			found = true
			break
		}
	}
	suite.False(found)
}

func (suite *ConfigTestSuite) Test_ConfigGetMember() {
	suite.NotEmpty(suite.config.Members)
	aMember := suite.config.Members[0].Name
	member, err := suite.config.GetMember(aMember)
	suite.NoError(err)
	suite.Equal(aMember, member.Name)
}

func (suite *ConfigTestSuite) Test_ConfigGetMemberNotPresent() {
	suite.NotEmpty(suite.config.Members)
	aMember := "foo"
	_, err := suite.config.GetMember(aMember)
	suite.Error(err)
}
