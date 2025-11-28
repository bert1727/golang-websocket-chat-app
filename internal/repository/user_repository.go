package repository

import (
	"github.com/bert1727/ChatApp/internal/models"
	"gorm.io/gorm"
)

func (r *userRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

type userRepo struct {
	db *gorm.DB
}

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}
