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
	entry1 := model.LogEntry{Term: uint64(42), Value: 42}
	err := suite.repo.Write(&entry1)
	assert.NoError(suite.T(), err)
	entry2, err := suite.repo.Read(entry1.Term)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entry1.Term, entry2.Term)
	assert.Equal(suite.T(), entry1.Value, entry2.Value)
}

func (suite *LogEntryRepositoryTestSuite) TestCount() {
	count, err := suite.repo.RowCount()
	assert.NoError(suite.T(), err)
	records := count + 20
	for i := count + 1; i <= records; i++ {
		entry := model.LogEntry{Term: uint64(i), Value: i}
		err = suite.repo.Write(&entry)
		assert.NoError(suite.T(), err, "failed on %d", i)
	}
	count2, err := suite.repo.RowCount()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count+records, count2)
}
