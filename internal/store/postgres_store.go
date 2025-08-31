package store

import (
	model "finalreg/internal/models"
	"finalreg/internal/providers"

	"gorm.io/gorm"
)

type PostgresStore struct {
	db *gorm.DB
}

func NewPostgresStore(db *gorm.DB) providers.RepoStore {
	return &PostgresStore{db: db}
}

//	func (s *PostgresStore) CreateUser(user *model.User) error {
//		return s.db.Create(user).Error
//	}
func (s *PostgresStore) CreateUser(user *model.User) (*model.User, error) {
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostgresStore) GetDB() *gorm.DB {
	return s.db
}
