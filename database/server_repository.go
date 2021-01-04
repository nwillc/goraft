package database

import "gorm.io/gorm"

type ServerRepository struct {

}

var _ GormRepository = (*ServerRepository)(nil)

func NewServerRepository() *ServerRepository {
	return nil
}

func (s ServerRepository) GetDB() *gorm.DB {
	panic("implement me")
}

func (s ServerRepository) RowCount() (int,error) {
	panic("implement me")
}

func (s ServerRepository) Migrate() error {
	panic("implement me")
}
