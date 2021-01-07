package database

import (
	"github.com/nwillc/goraft/pkg/model"
	"gorm.io/gorm"
)

type LogEntryRepository struct {
	db *gorm.DB
}

// LogEntryRepository implements Repository
var _ GormRepository = (*LogEntryRepository)(nil)

func NewLogEntryRepository(db *gorm.DB) (*LogEntryRepository, error) {
	repo := LogEntryRepository{
		db: db,
	}
	return &repo, nil
}

func (l *LogEntryRepository) GetDB() *gorm.DB {
	return l.db
}

func (l *LogEntryRepository) RowCount() (int, error) {
	var count int64
	l.db.Model(&model.LogEntry{}).Count(&count)
	return int(count), nil
}

func (l *LogEntryRepository) Migrate() error {
	if err := l.db.AutoMigrate(&model.LogEntry{}); err != nil {
		return err
	}
	return nil
}

func (l *LogEntryRepository) Write(logEntry *model.LogEntry) error {
	tx := l.db.Create(logEntry)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (l *LogEntryRepository) Read(term uint64) (*model.LogEntry, error) {
	var logs []model.LogEntry
	tx := l.db.Where("term = ?", term).Find(&logs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &logs[0], nil
}
