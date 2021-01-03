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

type StatusTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *StatusTestSuite) SetupTest() {
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
	err = suite.db.AutoMigrate(&Status{})
	assert.NoError(suite.T(), err)
}

func TestStatusTestSuite(t *testing.T) {
	suite.Run(t, new(StatusTestSuite))
}

func (suite *StatusTestSuite) TestStatusSanity() {
	assert.NotNil(suite.T(), suite.db)
	status := Status{}
	assert.Equal(suite.T(), uint64(0), status.Term)
}

func (suite *StatusTestSuite) TestStatusWrite() {
	term := uint64(42)
	name := "one"
	// Create a person
	stutus := Status{
		Name: name,
		Term: term,
	}

	// Persist it to database
	suite.db.Create(&stutus)
	assert.Equal(suite.T(), term, stutus.Term)

	// Select all
	var statuses []Status
	suite.db.Where("name = ?", name).Find(&statuses)
	assert.Equal(suite.T(), 1, len(statuses))
	assert.Equal(suite.T(), stutus, statuses[0])
}
