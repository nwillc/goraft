package main

import (
	"github.com/nwillc/goraft/model"
	"github.com/nwillc/goraft/raftapi"
	"github.com/nwillc/goraft/server"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	untilRetryDelay = 500 * time.Millisecond
	untilAttempts   = 30
)

type GoRaftTestSuite struct {
	suite.Suite
}

func (suite *GoRaftTestSuite) SetupTest() {
	suite.T().Helper()
	log.SetLevel(log.WarnLevel)
}

func TestGoRaftTestSuite(t *testing.T) {
	suite.Run(t, new(GoRaftTestSuite))
}

func (suite *GoRaftTestSuite) TestSingleServer() {
	svr := startEmbedded(suite.T(), "one")
	assert.NotNil(suite.T(), svr)
	time.Sleep(10 * time.Second)
	_, err := svr.Shutdown(nil, nil)
	assert.NoError(suite.T(), err)
	time.Sleep(5 * time.Second)
}

func (suite *GoRaftTestSuite) TestAllServer() {
	_ = startRaft(suite.T())
}

func (suite *GoRaftTestSuite) TestBasicHappyPathLogEntries() {
	servers := startRaft(suite.T())
	leader := hasLeader(servers)
	assert.NotNil(suite.T(), leader)
	v := int64(42)
	success, err := leader.AppendValue(nil, &raftapi.Value{Value: v})
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), success.Status)
	success, err = leader.AppendValue(nil, &raftapi.Value{Value: v + 1})
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), success.Status)
	for _, svr := range servers {
		response, err := svr.ListEntries(nil, nil)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), response.Entries, 2)
		assert.Equal(suite.T(), v, response.Entries[0].Value)
		assert.Equal(suite.T(), v+1, response.Entries[1].Value)
	}
}

func (suite *GoRaftTestSuite) TestBasicMemberDownLogEntries() {
	servers := startRaft(suite.T())
	leader := hasLeader(servers)
	assert.NotNil(suite.T(), leader)
	v := int64(42)
	success, err := leader.AppendValue(nil, &raftapi.Value{Value: v})
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), success.Status)
	member := nonLeader(servers)
	member.Shutdown(nil, nil)
	assert.True(suite.T(), until(untilAttempts, untilRetryDelay, func() bool {
		return member.GetState() == server.Shutdown
	}))
	success, err = leader.AppendValue(nil, &raftapi.Value{Value: v + 1})
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), success.Status)
	for _, svr := range servers {
		if svr.GetState() == server.Shutdown {
			continue
		}
		response, err := svr.ListEntries(nil, nil)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), response.Entries, 2)
		assert.Equal(suite.T(), v, response.Entries[0].Value)
		assert.Equal(suite.T(), v+1, response.Entries[1].Value)
	}
}

func startRaft(t *testing.T) []*server.RaftServer {
	t.Helper()
	var names = []string{"one", "two", "three", "four", "five"}
	var servers []*server.RaftServer
	for _, name := range names {
		servers = append(servers, startEmbedded(t, name))
	}
	assert.True(t, until(untilAttempts, untilRetryDelay, func() bool { return allRunning(servers) }))
	assert.True(t, until(untilAttempts, untilRetryDelay, func() bool { return hasLeader(servers) != nil }))
	t.Cleanup(func() {
		for _, svr := range servers {
			_, _ = svr.Shutdown(nil, nil)
		}
	})
	return servers
}

func startEmbedded(t *testing.T, name string) *server.RaftServer {
	t.Helper()
	config, err := model.ReadConfig("../config.json")
	assert.NoError(t, err)
	member, err := config.GetMember(name)
	assert.NoError(t, err)
	tempDB, err := ioutil.TempFile("", "test*.db")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(tempDB.Name())
	}()
	svr := server.NewRaftServer(*member, config, tempDB.Name())
	assert.NotNil(t, svr)
	go func() {
		err = svr.Run()
		assert.NoError(t, err)
	}()
	return svr
}

func allRunning(servers []*server.RaftServer) bool {
	for _, s := range servers {
		if s.GetState() != server.Running {
			return false
		}
	}
	return true
}

func hasLeader(servers []*server.RaftServer) *server.RaftServer {
	for _, s := range servers {
		if s.GetRole() == server.Leader {
			return s
		}
	}
	return nil
}

func nonLeader(servers []*server.RaftServer) *server.RaftServer {
	for _, s := range servers {
		if s.GetRole() != server.Leader {
			return s
		}
	}
	return nil
}

func until(attempts int, delay time.Duration, f func() bool) bool {
	for i := 0; i < attempts; i++ {
		if f() {
			return true
		}
		time.Sleep(delay)
	}
	return false
}
