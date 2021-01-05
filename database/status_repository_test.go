package database

import (
	"github.com/nwillc/goraft/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
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
			Name: strconv.Itoa(i),
			Term: uint64(i),
		}
		err = suite.repo.Write(&status)
		assert.NoError(suite.T(), err, "failed on %d", i)
	}
}

func (suite *StatusRepositoryTestSuite) TestWriteRead() {
	status1 := model.Status{
		Name: "foo",
		Term: 100,
	}
	err := suite.repo.Write(&status1)
	assert.NoError(suite.T(), err)
	status2, err := suite.repo.Read(status1.Name)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), status1.Term, status2.Term)
	assert.Equal(suite.T(), status1.Name, status2.Name)
}
