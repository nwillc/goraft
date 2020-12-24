package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type RaftLogTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *RaftLogTestSuite) SetupTest() {
	suite.T().Helper()
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	suite.T().Cleanup(func() {
		sqlDb, _ := db.DB()
		_ = sqlDb.Close()
	})
	suite.db = db
	err = suite.db.AutoMigrate(&RaftLog{})
	assert.NoError(suite.T(), err)
}

func TestRaftLogTestSuite(t *testing.T) {
	suite.Run(t, new(RaftLogTestSuite))
}

func (suite *RaftLogTestSuite) TestSanity() {
	assert.NotNil(suite.T(), suite.db)
	empty := RaftLog{}
	assert.Equal(suite.T(), uint(0), empty.Value)
}

func (suite *RaftLogTestSuite) TestWrite() {
	// Create a person
	log := RaftLog{Value: uint(42)}

	// Persist it to database
	suite.db.Create(&log)

	// Select all
	var logs []RaftLog
	suite.db.Find(&logs)
	assert.Greater(suite.T(), logs, 0)
	assert.Equal(suite.T(), log, logs[0])
}
