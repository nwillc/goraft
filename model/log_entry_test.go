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

type LogEntryTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *LogEntryTestSuite) SetupTest() {
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
	err = suite.db.AutoMigrate(&LogEntry{})
	assert.NoError(suite.T(), err)
}

func TestLogEntryTestSuite(t *testing.T) {
	suite.Run(t, new(LogEntryTestSuite))
}

func (suite *LogEntryTestSuite) TestLogEntrySanity() {
	assert.NotNil(suite.T(), suite.db)
	empty := LogEntry{}
	assert.Equal(suite.T(), uint(0), empty.Value)
}

func (suite *LogEntryTestSuite) TestLogEntryWrite() {
	term := uint(42)
	// Create a person
	log := LogEntry{Term: term, Value: uint(term)}

	// Persist it to database
	suite.db.Create(&log)
	assert.Equal(suite.T(), uint(term), log.Term)

	// Select all
	var logs []LogEntry
	suite.db.Where("term = ?", term).Find(&logs)
	assert.Equal(suite.T(), 1, len(logs))
	assert.Equal(suite.T(), log, logs[0])
}
