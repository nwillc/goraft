package database

import (
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
	db := TempDB(suite.T())
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
	truncate(suite.repo)
	var term uint64 = 64
	var value int64 = 234
	entryNo, err := suite.repo.Create(term, value)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), entryNo)
	entry2, err := suite.repo.Read(entryNo)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), term, entry2.Term)
	assert.Equal(suite.T(), value, entry2.Value)
}

func (suite *LogEntryRepositoryTestSuite) TestUpdate() {
	var term uint64 = 64
	var value int64 = 234
	entryNo, err := suite.repo.Create(term, value)
	assert.NoError(suite.T(), err)
	err = suite.repo.Update(entryNo, term+1, value+1)
	assert.NoError(suite.T(), err)
	entry2, err := suite.repo.Read(entryNo)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), term+1, entry2.Term)
	assert.Equal(suite.T(), value+1, entry2.Value)
}

func (suite *LogEntryRepositoryTestSuite) TestList() {
	truncate(suite.repo)
	err := suite.repo.TruncateToEntryNo(-1)
	assert.NoError(suite.T(), err)
	entries := 10
	for i := 1; i <= entries; i++ {
		_, err = suite.repo.Create(uint64(i), int64(i))
		assert.NoError(suite.T(), err)
	}

	list, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entries, len(list))
	for i, entry := range list {
		ii := i + 1
		assert.Equal(suite.T(), int64(ii), entry.EntryNo)
		assert.Equal(suite.T(), uint64(ii), entry.Term)
		assert.Equal(suite.T(), int64(ii), entry.Value)
	}
}

func (suite *LogEntryRepositoryTestSuite) TestCount() {
	count, err := suite.repo.RowCount()
	assert.NoError(suite.T(), err)
	records := count + 20
	for i := count + 1; i <= records; i++ {
		_, err = suite.repo.Create(uint64(i), int64(i))
		assert.NoError(suite.T(), err, "failed on %d", i)
	}
	count2, err := suite.repo.RowCount()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count+records, count2)
}

func (suite *LogEntryRepositoryTestSuite) TestLogSize() {
	var entries uint64 = 10
	for i := entries; i > 0; i-- {
		_, err := suite.repo.Create(i, int64(i))
		assert.NoError(suite.T(), err)
	}
	size, err := suite.repo.LogSize()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entries, size)
}

func (suite *LogEntryRepositoryTestSuite) TestLogSizeEmpty() {
	truncate(suite.repo)
	size, err := suite.repo.LogSize()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint64(0), size)
}

func (suite *LogEntryRepositoryTestSuite) TestMaxEntryNoEmptyTable() {
	truncate(suite.repo)
	max, _ := suite.repo.MaxEntryNo()
	assert.Equal(suite.T(), int64(-1), max)
}

func (suite *LogEntryRepositoryTestSuite) TestMaxEntryNoEmpty() {
	truncate(suite.repo)
	entryNo, err := suite.repo.MaxEntryNo()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(-1), entryNo)
}

func (suite *LogEntryRepositoryTestSuite) TestMaxEntryNo() {
	truncate(suite.repo)
	count, err := suite.repo.RowCount()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
	var term uint64 = 80
	var value int64 = 80
	id, err := suite.repo.Create(term, value)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), id)
	max, err := suite.repo.MaxEntryNo()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), id, max)
}

func (suite *LogEntryRepositoryTestSuite) TestLastEntryNoEntries() {
	truncate(suite.repo)
	_, err := suite.repo.LastEntry()
	assert.Error(suite.T(), err)
}

func (suite *LogEntryRepositoryTestSuite) TestLastEntry() {
	truncate(suite.repo)
	var term uint64 = 80
	var value int64 = 80
	_, err := suite.repo.Create(term, value)
	assert.NoError(suite.T(), err)
	term++
	value++
	_, err = suite.repo.Create(term, value)
	assert.NoError(suite.T(), err)
	entry, err := suite.repo.LastEntry()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), term, entry.Term)
	assert.Equal(suite.T(), value, entry.Value)
}

func truncate(repo *LogEntryRepository) {
	_ = repo.TruncateToEntryNo(-1)
}
