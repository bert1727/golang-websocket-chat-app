package repository

import (
	"errors"

	"github.com/bert1727/ChatApp/internal/config"
	"github.com/bert1727/ChatApp/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(userID uint) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
	Update(user *domain.User) error
	FindAll() (*domain.User, error)
	DeleteByID(userID uint) error
	FindByName(name string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: config.New().DB,
	}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	// if user.Email == "" || user.Username == "" {
	// 	log.Print("email and username are required")
	// 	return nil, fmt.Errorf("email and username are required")
	// }

	result := r.db.Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, ErrUserDuplicatedKey
			// return nil, fmt.Errorf("user with this email already exists")
		}
		return nil, ErrInternalServer
		// return nil, result.Error
	}

	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	if err := r.db.Model(user).Updates(user).Error; err != nil {
		return ErrInternalServer
	}
	return nil
}

func (r *userRepository) FindByID(userID uint) (*domain.User, error) {
	var user domain.User
	// NOTE: it will return nothing if there are no such user
	err := r.db.First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServer
	}

	return &user, nil
}

func (r *userRepository) FindAll() (*domain.User, error) {
	var users domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, ErrInternalServer
	}

	return &users, nil
}

func (r *userRepository) DeleteByID(userID uint) error {
	var user domain.User
	err := r.db.Delete(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return ErrInternalServer
	}

	return nil
}

func (r *userRepository) FindByName(name string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", name).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound || user.ID == 0 {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServer
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		} else {
			return nil, ErrInternalServer
		}
	}
	return &user, nil
}
