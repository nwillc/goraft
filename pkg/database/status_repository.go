package database

import (
	"github.com/nwillc/goraft/model"
	"gorm.io/gorm"
)

// StatusRepository is a GormRepository that supports model.Status.
type StatusRepository struct {
	db *gorm.DB
}

var _ GormRepository = (*StatusRepository)(nil)

// NewStatusRepository instantiates a new StatusRepository.
func NewStatusRepository(db *gorm.DB) (*StatusRepository, error) {
	repo := StatusRepository{
		db: db,
	}
	return &repo, nil
}

// RowCount returns the row count of model.Status entries.
func (s *StatusRepository) RowCount() (int, error) {
	var count int64
	s.db.Model(&model.Status{}).Count(&count)
	return int(count), nil
}

// Migrate the scheme for this repoitory.
func (s *StatusRepository) Migrate() error {
	if err := s.db.AutoMigrate(&model.Status{}); err != nil {
		return err
	}
	return nil
}

func (s *StatusRepository) Write(status *model.Status) error {
	tx := s.db.Model(&status).Update("term", status.Term)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		tx := s.db.Create(status)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}

func (s *StatusRepository) Read(name string) (*model.Status, error) {
	var statuses []model.Status
	tx := s.db.Where("name = ?", name).Find(&statuses)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return &statuses[0], nil
}
