package database

import (
	"github.com/nwillc/goraft/model"
	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
}

var _ GormRepository = (*ServerRepository)(nil)

func NewServerRepository(db *gorm.DB) (*ServerRepository, error) {
	repo := ServerRepository{
		db: db,
	}
	return &repo, nil
}

func (s ServerRepository) GetDB() *gorm.DB {
	return s.db
}

func (s ServerRepository) RowCount() (int, error) {
	var count int64
	s.db.Model(&model.Status{}).Count(&count)
	return int(count), nil
}

func (s ServerRepository) Migrate() error {
	if err := s.db.AutoMigrate(&model.Status{}); err != nil {
		return err
	}
	return nil
}
