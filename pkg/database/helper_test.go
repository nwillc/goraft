package database

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"testing"
)

func tempDB(t *testing.T) *gorm.DB {
	t.Helper()
	f, err := ioutil.TempFile("", "test*.db")
	assert.NoError(t, err)
	db, err := OpenSqlite(f.Name())
	if err != nil {
		panic("failed to open database")
	}
	t.Cleanup(func() {
		sqlDb, _ := db.DB()
		_ = sqlDb.Close()
		_ = os.Remove(f.Name())
	})
	return db
}
