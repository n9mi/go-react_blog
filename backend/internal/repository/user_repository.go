package repository

import (
	"backend/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByEmail(tx *gorm.DB, user *entity.User, email string) error {
	return tx.First(user, "email = ?", email).Error
}
