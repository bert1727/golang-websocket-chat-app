package service

import (
	"errors"

	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	FindUserByID(userID uint) (*domain.User, error)
	Register(username, password string) (*domain.User, error)
	DeleteUserByID(userID uint) error
	UpdateUser(user *domain.User) error
	FindUserByName(name string) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) FindUserByID(userID uint) (*domain.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateNewUser(u *domain.User) (*domain.User, error) {
	return s.repo.CreateNewUser(u)
}

func (s *userService) FindAllUsers() (*domain.User, error) {
	return s.repo.FindAllUsers()
}

func (s *userService) DeleteUserByID(userID uint) error {
	return s.repo.DeleteUserByID(userID)
}

func (s *userService) FindUserByName(name string) (*domain.User, error) {
	u, err := s.repo.FindUserByName(name)
	if err != nil {
		log.Info("user was not find", err)
		return nil, err
	}
	return u, nil
}

func (s *userService) Register(username, password string) (*domain.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// NOTE: add hashing and validaton for password
	user := &domain.User{
		Username: username,
		Password: password,
	}

	user, err := s.repo.CreateNewUser(user)
	if err != nil {
		log.Info("couldn't create a user", err)
	}
	return user, err
}

func (s *userService) UpdateUser(user *domain.User) error {
	return s.repo.UpdateUser(user)
}

func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Info("user wasn't found")
		return nil, err
	}
	return user, err
}
