package database

import (
	"github.com/nwillc/goraft/model"
	"gorm.io/gorm"
)

type StatusRepository struct {
	db *gorm.DB
}

var _ GormRepository = (*StatusRepository)(nil)

func NewServerRepository(db *gorm.DB) (*StatusRepository, error) {
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
