package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"testing"
)

type RaftLogEntryTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *RaftLogEntryTestSuite) SetupTest() {
	suite.T().Helper()
	f, err := ioutil.TempFile("", "test*.db")
	assert.NoError(suite.T(), err)
	db, err := gorm.Open(sqlite.Open(f.Name()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	suite.T().Cleanup(func() {
		sqlDb, _ := db.DB()
		_ = sqlDb.Close()
		_ = os.Remove(f.Name())
	})
	suite.db = db
	err = suite.db.AutoMigrate(&RaftLogEntry{})
	assert.NoError(suite.T(), err)
}

func TestRaftLogEntryTestSuite(t *testing.T) {
	suite.Run(t, new(RaftLogEntryTestSuite))
}

func (suite *RaftLogEntryTestSuite) TestRaftLogEntrySanity() {
	assert.NotNil(suite.T(), suite.db)
	empty := RaftLogEntry{}
	assert.Equal(suite.T(), uint(0), empty.Value)
}

func (suite *RaftLogEntryTestSuite) TestRaftLogEntryWrite() {
	term := uint(42)
	// Create a person
	log := RaftLogEntry{Term: term, Value: uint(term)}

	// Persist it to database
	suite.db.Create(&log)
	assert.Equal(suite.T(), uint(term), log.Term)

	// Select all
	var logs []RaftLogEntry
	suite.db.Where("term = ?", term).Find(&logs)
	assert.Equal(suite.T(), 1, len(logs))
	assert.Equal(suite.T(), log, logs[0])
}
