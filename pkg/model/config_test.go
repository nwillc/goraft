package model

import (
	"github.com/nwillc/goraft/conf"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(suite.T(), err)
	suite.config = config
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) TestFailBadFile() {
	_, err := ReadConfig("foo.json")
	assert.Error(suite.T(), err)
}

func (suite *ConfigTestSuite) TestReadConfig() {
	assert.Less(suite.T(), 1, suite.config.HeartbeatTimeout)
	assert.Less(suite.T(), 1, suite.config.ElectionTimeout)
	assert.Less(suite.T(), 0, suite.config.MinOffset)
	assert.Less(suite.T(), suite.config.MinOffset, suite.config.MaxOffset)
	for _, member := range suite.config.Members {
		assert.NotEmpty(suite.T(), member.Name)
		assert.Less(suite.T(), uint32(0), member.Port)
	}
}

func (suite *ConfigTestSuite) TestConfigPeers() {
	assert.NotEmpty(suite.T(), suite.config.Members)
	aMember := suite.config.Members[0].Name
	peers := suite.config.Peers(aMember)
	assert.Equal(suite.T(), len(suite.config.Members) - 1, len(peers))
	var found bool
	for _, member := range peers {
		if member.Name == aMember {
			found = true
			break
		}
	}
	assert.False(suite.T(), found)
}

func (suite *ConfigTestSuite) TestConfigGetMember() {
	assert.NotEmpty(suite.T(), suite.config.Members)
	aMember := suite.config.Members[0].Name
	member, err := suite.config.GetMember(aMember)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), aMember, member.Name)
}
