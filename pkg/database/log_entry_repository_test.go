package database

import (
	"github.com/nwillc/goraft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LogEntryRepositoryTestSuite struct {
	suite.Suite
	repo *LogEntryRepository
}

func (suite *LogEntryRepositoryTestSuite) SetupTest() {
	suite.T().Helper()
	db := tempDB(suite.T())
	repo, err := NewLogEntryRepository(db)
	assert.NoError(suite.T(), err)
	suite.repo = repo
	err = suite.repo.Migrate()
	assert.NoError(suite.T(), err)
}

func TestLogEntryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LogEntryRepositoryTestSuite))
}

func (suite *LogEntryRepositoryTestSuite) TestWriteRead() {
	var term uint64 = 64
	var value = 234
	id, err := suite.repo.Write(term, value)
	assert.NoError(suite.T(), err)
	entry2, err := suite.repo.Read(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), term, entry2.Term)
	assert.Equal(suite.T(), value, entry2.Value)
}

func (suite *LogEntryRepositoryTestSuite) TestTermConstraint() {
	var term uint64 = 80
	var value = 80
	_, err := suite.repo.Write(term, value)
	assert.NoError(suite.T(), err)
	_, err = suite.repo.Write(term, value + 1)
	assert.Error(suite.T(), err)
}

func (suite *LogEntryRepositoryTestSuite) TestCount() {
	count, err := suite.repo.RowCount()
	assert.NoError(suite.T(), err)
	records := count + 20
	for i := count + 1; i <= records; i++ {
		_, err = suite.repo.Write(uint64(i), i)
		assert.NoError(suite.T(), err, "failed on %d", i)
	}
	count2, err := suite.repo.RowCount()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count+records, count2)
}

func (suite *LogEntryRepositoryTestSuite) TestMaxTerm() {
	maxTerm := 10
	for i := maxTerm; i > 0; i-- {
		_, err := suite.repo.Write(uint64(i), i)
		assert.NoError(suite.T(), err)
	}
	max, err := suite.repo.MaxTerm()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint64(maxTerm), max)
}

func (suite *LogEntryRepositoryTestSuite) TestMaxIdEmpty() {
	truncate(suite.repo)
	max, _ := suite.repo.MaxId()
	assert.Equal(suite.T(), int64(-1), max)
}

func (suite *LogEntryRepositoryTestSuite) TestMaxId() {
	truncate(suite.repo)
	var term uint64 = 80
	var value = 80
	id, err := suite.repo.Write(term, value)
	assert.NoError(suite.T(), err)
	max, _ := suite.repo.MaxId()
	assert.Equal(suite.T(), id, max)
}

func truncate(repo *LogEntryRepository)  {
	repo.db.Where("1 = 1").Delete(&model.LogEntry{})
}
