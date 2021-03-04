package server

import (
	"context"
	"github.com/nwillc/goraft/conf"
	"github.com/nwillc/goraft/model"
	"github.com/nwillc/goraft/raftapi"
	log "github.com/sirupsen/logrus"
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
	suite.NoError(err)
	member := config.Members[0]
	tempDB, err := ioutil.TempFile("", "test*.db")
	suite.NoError(err)
	suite.T().Cleanup(func() {
		_ = os.Remove(tempDB.Name())
	})
	suite.server = NewRaftServer(member, config, tempDB.Name())
	err = suite.server.setupRepositories()
	suite.NoError(err)
	suite.ctx = context.Background()
	log.SetLevel(log.WarnLevel)
}

func TestRaftServerAppendEntrySuite(t *testing.T) {
	suite.Run(t, new(RaftServerAppendEntryTestSuite))
}

func (suite *RaftServerAppendEntryTestSuite) Test_FirstEntry() {
	enrtyNo, err := suite.server.logRepo.MaxEntryNo()
	suite.NoError(err)
	suite.Equal(int64(-1), enrtyNo)
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
	suite.NoError(err)
	suite.Equal(true, response.Success)
	suite.Equal(term, response.Term)
	enrtyNo, err = suite.server.logRepo.MaxEntryNo()
	suite.NoError(err)
	suite.Equal(int64(1), enrtyNo)
	logEntry, err := suite.server.logRepo.Read(enrtyNo)
	suite.NoError(err)
	suite.Equal(term, logEntry.Term)
	suite.Equal(value, logEntry.Value)
}

func (suite *RaftServerAppendEntryTestSuite) Test_TwoSuccessiveEntries() {
	entryNo, err := suite.server.logRepo.MaxEntryNo()
	suite.NoError(err)
	suite.Equal(int64(-1), entryNo)
	previousLogEntry := int64(-1)
	term := uint64(0)
	value := int64(42)
	request := &raftapi.AppendEntryRequest{
		Term:        term,
		Leader:      "",
		PrevLogId:   previousLogEntry,
		PrevLogTerm: term,
		Entry: &raftapi.LogEntry{
			Term:  term,
			Value: value,
		},
	}
	response, err := suite.server.AppendEntry(suite.ctx, request)
	suite.NoError(err)
	suite.Equal(true, response.Success)
	suite.Equal(term, response.Term)
	entryNo, err = suite.server.logRepo.MaxEntryNo()
	suite.NoError(err)
	suite.Equal(int64(1), entryNo)
	logEntry, err := suite.server.logRepo.Read(entryNo)
	suite.NoError(err)
	suite.Equal(term, logEntry.Term)
	suite.Equal(value, logEntry.Value)
	value++
	request2 := &raftapi.AppendEntryRequest{
		Term:        term,
		Leader:      "",
		PrevLogId:   logEntry.EntryNo,
		PrevLogTerm: logEntry.Term,
		Entry: &raftapi.LogEntry{
			Term:  term,
			Value: value,
		},
	}
	response2, err := suite.server.AppendEntry(suite.ctx, request2)
	suite.NoError(err)
	suite.Equal(true, response2.Success)
	suite.Equal(term, response2.Term)
	entries, err := suite.server.ListEntries(suite.ctx, &raftapi.Empty{})
	suite.NoError(err)
	suite.Equal(2, len(entries.Entries))
	for i, entry := range entries.Entries {
		suite.Equal(int64(42+i), entry.Value)
	}
}
