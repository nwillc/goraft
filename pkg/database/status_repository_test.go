package database

import (
	"github.com/nwillc/goraft/model"
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
	db := TempDB(suite.T())
	repo, err := NewStatusRepository(db)
	suite.NoError(err)
	suite.repo = repo
	err = suite.repo.Migrate()
	suite.NoError(err)
}

func TestStatusRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StatusRepositoryTestSuite))
}

func (suite *StatusRepositoryTestSuite) Test_Count() {
	count, err := suite.repo.RowCount()
	suite.NoError(err)
	records := count + 20
	for i := count + 1; i <= records; i++ {
		status := model.Status{
			Name: strconv.Itoa(i),
			Term: uint64(i),
		}
		err = suite.repo.Write(&status)
		suite.NoError(err, "failed on %d", i)
	}
}

func (suite *StatusRepositoryTestSuite) Test_WriteRead() {
	status1 := model.Status{
		Name: "foo",
		Term: 100,
	}
	err := suite.repo.Write(&status1)
	suite.NoError(err)
	status2, err := suite.repo.Read(status1.Name)
	suite.NoError(err)
	suite.Equal(status1.Term, status2.Term)
	suite.Equal(status1.Name, status2.Name)
}
