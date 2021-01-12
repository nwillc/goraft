package database

import (
	"github.com/nwillc/goraft/model"
	"gorm.io/gorm"
)

// LogEntryRepository is a GormRepository focused on model.LogEntry.
type LogEntryRepository struct {
	db *gorm.DB
}

// LogEntryRepository implements Repository
var _ GormRepository = (*LogEntryRepository)(nil)

// NewLogEntryRepository instantiates a LogEntryRepository.
func NewLogEntryRepository(db *gorm.DB) (*LogEntryRepository, error) {
	repo := LogEntryRepository{
		db: db,
	}
	return &repo, nil
}

// RowCount returns the row count LogEntryRepository.
func (l *LogEntryRepository) RowCount() (int, error) {
	var count int64
	l.db.Model(&model.LogEntry{}).Count(&count)
	return int(count), nil
}

// Migrate the schema for the LogEntryRepository.
func (l *LogEntryRepository) Migrate() error {
	if err := l.db.AutoMigrate(&model.LogEntry{}); err != nil {
		return err
	}
	return nil
}

func (l *LogEntryRepository) Write(term uint64, value int) (int64, error) {
	entry := model.LogEntry{
		Term:  term,
		Value: value,
	}
	tx := l.db.Create(&entry)
	if tx.Error != nil {
		return -1, tx.Error
	}
	return entry.ID, nil
}

func (l *LogEntryRepository) Read(id int64) (*model.LogEntry, error) {
	var logs []model.LogEntry
	tx := l.db.Where("id = ?", id).Find(&logs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &logs[0], nil
}

func (l *LogEntryRepository) MaxTerm() (uint64, error) {
	var result uint64
	tx := l.db.Model(&model.LogEntry{}).Select("max(term)").Row()
	err := tx.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (l *LogEntryRepository) MaxId() (int64, error) {
	var result int64
	tx := l.db.Model(&model.LogEntry{}).Select("max(id)").Row()
	err := tx.Scan(&result)
	if err != nil {
		return -1, err
	}
	return result, nil
}
