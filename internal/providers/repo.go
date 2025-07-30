package providers

import (
	model "finalreg/internal/models"

	"gorm.io/gorm"
)

type RepoStore interface {
	Repository
}

type Repository interface {
	UserRepositories
}

type UserRepositories interface {
	CreateUser(user *model.User) error
	GetDB() *gorm.DB
}
