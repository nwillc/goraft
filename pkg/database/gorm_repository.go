package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormRepository is a database repository interface which uses GORM for database access.
type GormRepository interface {
	RowCount() (int, error)
	Migrate() error
}

// OpenSqlite opens a Sqlite database located at filename and returns a gorm.DB reference.
func OpenSqlite(filename string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db, nil
}
