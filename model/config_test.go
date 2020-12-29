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
	assert.Less(suite.T(), uint(0), config.Term)
	assert.Equal(suite.T(), config.Term, uint(len(config.Members)))
	for _, member := range config.Members {
		assert.Less(suite.T(), uint32(0), member.Port)
	}
}
