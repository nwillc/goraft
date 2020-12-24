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

type TermTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *TermTestSuite) SetupTest() {
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
	err = suite.db.AutoMigrate(&Term{})
	assert.NoError(suite.T(), err)
}

func TestTermTestSuite(t *testing.T) {
	suite.Run(t, new(TermTestSuite))
}

func (suite *TermTestSuite) TestTermSanity() {
	assert.NotNil(suite.T(), suite.db)
	empty := Term{}
	assert.Equal(suite.T(), uint(0), empty.ID)
}

func (suite *TermTestSuite) TestTermWrite() {
	id := uint(42)
	// Create a person
	log := Term{ID: uint(id)}

	// Persist it to database
	suite.db.Create(&log)
	assert.Equal(suite.T(), uint(id), log.ID)

	// Select all
	var logs []Term
	suite.db.Where("id = ?", id).Find(&logs)
	assert.Equal(suite.T(), 1, len(logs))
	assert.Equal(suite.T(), log, logs[0])
}
