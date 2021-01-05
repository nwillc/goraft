package database

import (
	"github.com/nwillc/goraft/model"
	"gorm.io/gorm"
)

type StatusRepository struct {
	db *gorm.DB
}

var _ GormRepository = (*StatusRepository)(nil)

func NewStatusRepository(db *gorm.DB) (*StatusRepository, error) {
	repo := StatusRepository{
		db: db,
	}
	return &repo, nil
}

func (s *StatusRepository) GetDB() *gorm.DB {
	return s.db
}

func (s *StatusRepository) RowCount() (int, error) {
	var count int64
	s.db.Model(&model.Status{}).Count(&count)
	return int(count), nil
}

func (s *StatusRepository) Migrate() error {
	if err := s.db.AutoMigrate(&model.Status{}); err != nil {
		return err
	}
	return nil
}

func (s *StatusRepository) Write(status *model.Status) error {
	tx := s.db.Create(status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *StatusRepository) Read(name string) (*model.Status, error) {
	var statuses []model.Status
	tx := s.db.Where("name = ?", name).Find(&statuses)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &statuses[0], nil
}
