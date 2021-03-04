package database

import (
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
	suite.NoError(err)
	suite.repo = repo
	err = suite.repo.Migrate()
	suite.NoError(err)
}

func TestLogEntryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LogEntryRepositoryTestSuite))
}

func (suite *LogEntryRepositoryTestSuite) Test_WriteRead() {
	truncate(suite.repo)
	var term uint64 = 64
	var value int64 = 234
	entryNo, err := suite.repo.Create(term, value)
	suite.NoError(err)
	suite.Equal(int64(1), entryNo)
	entry2, err := suite.repo.Read(entryNo)
	suite.NoError(err)
	suite.Equal(term, entry2.Term)
	suite.Equal(value, entry2.Value)
}

func (suite *LogEntryRepositoryTestSuite) Test_Update() {
	var term uint64 = 64
	var value int64 = 234
	entryNo, err := suite.repo.Create(term, value)
	suite.NoError(err)
	err = suite.repo.Update(entryNo, term+1, value+1)
	suite.NoError(err)
	entry2, err := suite.repo.Read(entryNo)
	suite.NoError(err)
	suite.Equal(term+1, entry2.Term)
	suite.Equal(value+1, entry2.Value)
}

func (suite *LogEntryRepositoryTestSuite) Test_List() {
	truncate(suite.repo)
	err := suite.repo.TruncateToEntryNo(-1)
	suite.NoError(err)
	entries := 10
	for i := 1; i <= entries; i++ {
		_, err = suite.repo.Create(uint64(i), int64(i))
		suite.NoError(err)
	}

	list, err := suite.repo.List()
	suite.NoError(err)
	suite.Equal(entries, len(list))
	for i, entry := range list {
		ii := i + 1
		suite.Equal(int64(ii), entry.EntryNo)
		suite.Equal(uint64(ii), entry.Term)
		suite.Equal(int64(ii), entry.Value)
	}
}

func (suite *LogEntryRepositoryTestSuite) Test_Count() {
	count, err := suite.repo.RowCount()
	suite.NoError(err)
	records := count + 20
	for i := count + 1; i <= records; i++ {
		_, err = suite.repo.Create(uint64(i), int64(i))
		suite.NoError(err, "failed on %d", i)
	}
	count2, err := suite.repo.RowCount()
	suite.NoError(err)
	suite.Equal(count+records, count2)
}

func (suite *LogEntryRepositoryTestSuite) Test_LogSize() {
	var entries uint64 = 10
	for i := entries; i > 0; i-- {
		_, err := suite.repo.Create(i, int64(i))
		suite.NoError(err)
	}
	size, err := suite.repo.LogSize()
	suite.NoError(err)
	suite.Equal(entries, size)
}

func (suite *LogEntryRepositoryTestSuite) Test_LogSizeEmpty() {
	truncate(suite.repo)
	size, err := suite.repo.LogSize()
	suite.NoError(err)
	suite.Equal(uint64(0), size)
}

func (suite *LogEntryRepositoryTestSuite) Test_MaxEntryNoEmptyTable() {
	truncate(suite.repo)
	max, _ := suite.repo.MaxEntryNo()
	suite.Equal(int64(-1), max)
}

func (suite *LogEntryRepositoryTestSuite) Test_MaxEntryNoEmpty() {
	truncate(suite.repo)
	entryNo, err := suite.repo.MaxEntryNo()
	suite.NoError(err)
	suite.Equal(int64(-1), entryNo)
}

func (suite *LogEntryRepositoryTestSuite) Test_MaxEntryNo() {
	truncate(suite.repo)
	count, err := suite.repo.RowCount()
	suite.NoError(err)
	suite.Equal(0, count)
	var term uint64 = 80
	var value int64 = 80
	id, err := suite.repo.Create(term, value)
	suite.NoError(err)
	suite.Equal(int64(1), id)
	max, err := suite.repo.MaxEntryNo()
	suite.NoError(err)
	suite.Equal(id, max)
}

func (suite *LogEntryRepositoryTestSuite) Test_LastEntryNoEntries() {
	truncate(suite.repo)
	_, err := suite.repo.LastEntry()
	suite.Error(err)
}

func (suite *LogEntryRepositoryTestSuite) Test_LastEntry() {
	truncate(suite.repo)
	var term uint64 = 80
	var value int64 = 80
	_, err := suite.repo.Create(term, value)
	suite.NoError(err)
	term++
	value++
	_, err = suite.repo.Create(term, value)
	suite.NoError(err)
	entry, err := suite.repo.LastEntry()
	suite.NoError(err)
	suite.Equal(term, entry.Term)
	suite.Equal(value, entry.Value)
}

func truncate(repo *LogEntryRepository) {
	_ = repo.TruncateToEntryNo(-1)
}
