package service

import (
	"errors"

	"github.com/bert1727/ChatApp/internal/models"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	FindUserByID(userID uint) (*models.User, error)
	Register(username, password string) (*models.User, error)
	DeleteUserByID(userID uint) error
	UpdateUser(user *models.User) error
	FindUserByName(name string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) FindUserByID(userID uint) (*models.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateNewUser(u *models.User) error {
	return s.repo.CreateNewUser(u)
}

func (s *userService) FindAllUsers() (*models.User, error) {
	return s.repo.FindAllUsers()
}

func (s *userService) DeleteUserByID(userID uint) error {
	return s.repo.DeleteUserByID(userID)
}

func (s *userService) FindUserByName(name string) (*models.User, error) {
	u, err := s.repo.FindUserByName(name)
	if err != nil {
		log.Info("user was not find", err)
		return nil, err
	}
	return u, nil
}

func (s *userService) Register(username, password string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// NOTE: add hashing and validaton for password
	user := &models.User{
		Name:     username,
		Password: password,
	}

	err := s.repo.CreateNewUser(user)
	return user, err
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.repo.UpdateUser(user)
}
