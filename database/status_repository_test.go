package database

import (
	"github.com/nwillc/goraft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StatusRepositoryTestSuite struct {
	suite.Suite
	repo *StatusRepository
}

func (suite *StatusRepositoryTestSuite) SetupTest() {
	suite.T().Helper()
	db := tempDB(suite.T())
	repo, err := NewStatusRepository(db)
	assert.NoError(suite.T(), err)
	suite.repo = repo
	err = suite.repo.Migrate()
	assert.NoError(suite.T(), err)
}

func TestStatusRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StatusRepositoryTestSuite))
}

func (suite *StatusRepositoryTestSuite) TestCount() {
	count, err := suite.repo.RowCount()
	assert.NoError(suite.T(), err)
	records := count + 20
	for i := count + 1; i <= records; i++ {
		status := model.Status{
			Name: "foo",
			Term: uint64(i),
		}
		err = suite.repo.Write(&status)
		assert.NoError(suite.T(), err, "failed on %d", i)
	}
}
