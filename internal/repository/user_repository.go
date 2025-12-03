package repository

import (
	"github.com/bert1727/ChatApp/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByID(userID uint) (*models.User, error)
	CreateNewUser(user *models.User) error
	UpdateUser(user *models.User) error
	FindAllUsers() (*models.User, error)
	DeleteUserByID(userID uint) error
	FindUserByName(name string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) CreateNewUser(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) UpdateUser(user *models.User) error {
	return u.db.Model(user).Updates(user).Error
}

func (u *userRepository) FindUserByID(userID uint) (*models.User, error) {
	var user models.User
	// NOTE: it will return nothing if there are no such user
	if err := u.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) FindAllUsers() (*models.User, error) {
	var users models.User
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func (u *userRepository) DeleteUserByID(userID uint) error {
	var user models.User
	return u.db.Delete(&user, userID).Error
}

func (u *userRepository) FindUserByName(name string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
