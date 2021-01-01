package model

import (
	"github.com/nwillc/goraft/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.T().Helper()
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) TestFailBadFile() {
	_, err := ReadConfig("foo.json")
	assert.Error(suite.T(), err)
}

func (suite *ConfigTestSuite) TestReadConfig() {
	config, err := ReadConfig("../" + conf.ConfigFile)
	assert.NoError(suite.T(), err)
	assert.Less(suite.T(), 1, config.HeartbeatTimeout)
	assert.Less(suite.T(), 1, config.ElectionTimeout)
	assert.Less(suite.T(), 0, config.MinOffset)
	assert.Less(suite.T(), config.MinOffset, config.MaxOffset)
	for _, member := range config.Members {
		assert.NotEmpty(suite.T(), member.Name)
		assert.Less(suite.T(), uint32(0), member.Port)
	}
}
