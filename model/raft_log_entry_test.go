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

type RaftLogTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *RaftLogTestSuite) SetupTest() {
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

func TestRaftLogTestSuite(t *testing.T) {
	suite.Run(t, new(RaftLogTestSuite))
}

func (suite *RaftLogTestSuite) TestSanity() {
	assert.NotNil(suite.T(), suite.db)
	empty := RaftLogEntry{}
	assert.Equal(suite.T(), uint(0), empty.Value)
}

func (suite *RaftLogTestSuite) TestWrite() {
	// Create a person
	log := RaftLogEntry{Value: uint(42)}
	assert.Equal(suite.T(), uint(0), log.ID)

	// Persist it to database
	suite.db.Create(&log)
	assert.NotEqual(suite.T(), uint(0), log.ID)

	// Select all
	var logs []RaftLogEntry
	suite.db.Where("id = ?", log.ID).Find(&logs)
	assert.Equal(suite.T(), 1, len(logs))
	assert.Equal(suite.T(), log, logs[0])
}
