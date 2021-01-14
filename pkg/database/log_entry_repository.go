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

func (l *LogEntryRepository) Create(term uint64, value int64) (int64, error) {
	index, err := l.RowCount()
	if err != nil {
		return -1, err
	}
	entry := model.LogEntry{
		EntryNo: int64(index),
		Term:    term,
		Value:   value,
	}
	tx := l.db.Create(&entry)
	if tx.Error != nil || tx.RowsAffected != 1 {
		return -1, tx.Error
	}
	return entry.EntryNo, nil
}

func (l *LogEntryRepository) Read(id int64) (*model.LogEntry, error) {
	var logs []model.LogEntry
	tx := l.db.Where("entry_no = ?", id).Find(&logs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &logs[0], nil
}

func (l *LogEntryRepository) Update(id int64, term uint64, value int64) error {
	entry := model.LogEntry{
		EntryNo: id,
		Term:    term,
		Value:   value,
	}
	tx := l.db.Model(&entry).Where("entry_no = ?", id).Updates(entry)
	return tx.Error
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

func (l *LogEntryRepository) MaxEntryNo() (int64, error) {
	var result int64
	tx := l.db.Model(&model.LogEntry{}).Select("max(entry_no)").Row()
	err := tx.Scan(&result)
	if err != nil {
		return -1, err
	}
	return result, nil
}

func (l *LogEntryRepository) TruncateToEntryNo(entryNo int64) error {
	tx := l.db.Where("row_entry > ?", entryNo).Delete(&model.LogEntry{})
	return tx.Error
}
