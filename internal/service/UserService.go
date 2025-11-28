package service

import (
	"github.com/bert1727/ChatApp/internal/models"
	"github.com/bert1727/ChatApp/internal/repository"
)

type UserService interface {
	FindUserByID(userID uint) (*models.User, error)
	// DeleteUserByID(userID string) error
}

type userService struct {
	repo repository.UserRepository
}

func (s *userService) FindUserByID(userID uint) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
