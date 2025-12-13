package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/bert1727/ChatApp/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByID(userID uint) (*domain.User, error)
	CreateNewUser(user *domain.User) (*domain.User, error)
	UpdateUser(user *domain.User) error
	FindAllUsers() (*domain.User, error)
	DeleteUserByID(userID uint) error
	FindUserByName(name string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) CreateNewUser(user *domain.User) (*domain.User, error) {
	if user.Email == "" || user.Username == "" {
		log.Print("email and username are required")
		return nil, fmt.Errorf("email and username are required")
	}

	result := r.db.Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("user with this email already exists")
		}
		return nil, result.Error
	}

	return user, nil
}

func (r *userRepository) UpdateUser(user *domain.User) error {
	return r.db.Model(user).Updates(user).Error
}

func (r *userRepository) FindUserByID(userID uint) (*domain.User, error) {
	var user domain.User
	// NOTE: it will return nothing if there are no such user
	err := r.db.First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		} else {
			return nil, errors.New("failed to find a user")
		}
	}

	return &user, nil
}

func (r *userRepository) FindAllUsers() (*domain.User, error) {
	var users domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *userRepository) DeleteUserByID(userID uint) error {
	var user domain.User
	return r.db.Delete(&user, userID).Error
}

func (r *userRepository) FindUserByName(name string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("username = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("user not found")
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		} else {
			return nil, errors.New("failed to get user by email")
		}
	}
	return &user, nil
}
