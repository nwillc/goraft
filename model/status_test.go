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

func TestStatusTestSuite(t *testing.T)  {
	suite.Run(t, new(StatusTestSuite))
}

func (suite *StatusTestSuite) TestStatusSanity() {
	assert.NotNil(suite.T(), suite.db)
	status := Status{}
	assert.Equal(suite.T(), uint64(0), status.Term)
}
