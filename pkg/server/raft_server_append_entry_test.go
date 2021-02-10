package server

import (
	"context"
	"github.com/nwillc/goraft/conf"
	"github.com/nwillc/goraft/model"
	"github.com/nwillc/goraft/raftapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

type RaftServerAppendEntryTestSuite struct {
	suite.Suite
	server *RaftServer
	ctx    context.Context
}

func (suite *RaftServerAppendEntryTestSuite) SetupTest() {
	suite.T().Helper()
	config, err := model.ReadConfig("../../" + conf.ConfigFile)
	assert.NoError(suite.T(), err)
	member := config.Members[0]
	tempDB, err := ioutil.TempFile("", "test*.db")
	assert.NoError(suite.T(), err)
	suite.T().Cleanup(func() {
		_ = os.Remove(tempDB.Name())
	})
	suite.server = NewRaftServer(member, config, tempDB.Name())
	err = suite.server.setupRepositories()
	assert.NoError(suite.T(), err)
	suite.ctx = context.Background()
}

func TestRaftServerAppendEntrySuite(t *testing.T) {
	suite.Run(t, new(RaftServerAppendEntryTestSuite))
}

func (suite *RaftServerAppendEntryTestSuite) TestFirstEntry() {
	enrtyNo, err := suite.server.logRepo.MaxEntryNo()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(-1), enrtyNo)
	term := uint64(0)
	value := int64(42)
	request := &raftapi.AppendEntryRequest{
		Term:        term,
		Leader:      "",
		PrevLogId:   -1,
		PrevLogTerm: 0,
		Entry: &raftapi.LogEntry{
			Term:  term,
			Value: value,
		},
	}
	response, err := suite.server.AppendEntry(suite.ctx, request)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response.Success)
	assert.Equal(suite.T(), term, response.Term)
	enrtyNo, err = suite.server.logRepo.MaxEntryNo()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), enrtyNo)
	logEntry, err := suite.server.logRepo.Read(enrtyNo)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), term, logEntry.Term)
	assert.Equal(suite.T(), value, logEntry.Value)
}

func (suite *RaftServerAppendEntryTestSuite) TestTwoSuccessiveEntries() {
	entryNo, err := suite.server.logRepo.MaxEntryNo()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(-1), entryNo)
	term := uint64(0)
	value := int64(42)
	request := &raftapi.AppendEntryRequest{
		Term:        term,
		Leader:      "",
		PrevLogId:   -1,
		PrevLogTerm: 0,
		Entry: &raftapi.LogEntry{
			Term:  term,
			Value: value,
		},
	}
	response, err := suite.server.AppendEntry(suite.ctx, request)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response.Success)
	assert.Equal(suite.T(), term, response.Term)
	entryNo, err = suite.server.logRepo.MaxEntryNo()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), entryNo)
	logEntry, err := suite.server.logRepo.Read(entryNo)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), term, logEntry.Term)
	assert.Equal(suite.T(), value, logEntry.Value)
	value++
	request2 := &raftapi.AppendEntryRequest{
		Term:        term,
		Leader:      "",
		PrevLogId:   0,
		PrevLogTerm: term,
		Entry: &raftapi.LogEntry{
			Term:  term,
			Value: value,
		},
	}
	response2, err := suite.server.AppendEntry(suite.ctx, request2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response2.Success)
	assert.Equal(suite.T(), term, response2.Term)
	entries, err := suite.server.ListEntries(suite.ctx, &raftapi.Empty{})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(entries.Entries))
	for i, entry := range entries.Entries {
		assert.Equal(suite.T(), int64(42+i), entry.Value)
	}
}
