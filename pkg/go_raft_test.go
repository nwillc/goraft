package main

import (
	"context"
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

func (suite *GoRaftTestSuite) Test_SingleServer() {
	svr := startEmbedded(suite.T(), "one")
	suite.NotNil(svr)
	time.Sleep(10 * time.Second)
	_, err := svr.Shutdown(context.TODO(), nil)
	suite.NoError(err)
	time.Sleep(5 * time.Second)
}

func (suite *GoRaftTestSuite) Test_AllServer() {
	_ = startRaft(suite.T())
}

func (suite *GoRaftTestSuite) Test_BasicHappyPathLogEntries() {
	servers := startRaft(suite.T())
	leader := hasLeader(servers)
	suite.NotNil(leader)
	v := int64(42)
	success, err := leader.AppendValue(context.TODO(), &raftapi.Value{Value: v})
	suite.NoError(err)
	suite.True(success.Status)
	success, err = leader.AppendValue(context.TODO(), &raftapi.Value{Value: v + 1})
	suite.NoError(err)
	suite.True(success.Status)
	for _, svr := range servers {
		response, err := svr.ListEntries(context.TODO(), nil)
		suite.NoError(err)
		suite.Len(response.Entries, 2)
		suite.Equal(v, response.Entries[0].Value)
		suite.Equal(v+1, response.Entries[1].Value)
	}
}

func (suite *GoRaftTestSuite) Test_BasicMemberDownLogEntries() {
	servers := startRaft(suite.T())
	leader := hasLeader(servers)
	suite.NotNil(leader)
	v := int64(42)
	success, err := leader.AppendValue(context.TODO(), &raftapi.Value{Value: v})
	suite.NoError(err)
	suite.True(success.Status)
	member := nonLeader(servers)
	_, _ = member.Shutdown(context.TODO(), nil)
	suite.True(until(untilAttempts, untilRetryDelay, func() bool {
		return member.GetState() == server.Shutdown
	}))
	success, err = leader.AppendValue(context.TODO(), &raftapi.Value{Value: v + 1})
	suite.NoError(err)
	suite.True(success.Status)
	for _, svr := range servers {
		if svr.GetState() == server.Shutdown {
			continue
		}
		response, err := svr.ListEntries(context.TODO(), nil)
		suite.NoError(err)
		suite.Len(response.Entries, 2)
		suite.Equal(v, response.Entries[0].Value)
		suite.Equal(v+1, response.Entries[1].Value)
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
			_, _ = svr.Shutdown(context.TODO(), nil)
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
