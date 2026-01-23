package service

import (
	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
)

var validate *validator.Validate

type UserService interface {
	FindUserByID(userID uint) (*domain.User, error)
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
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// func (s *userService) CreateNewUser(user *domain.User) (*domain.User, error) {
// 	err := validate.Struct(*user)
// 	validationErrors := err.(validator.ValidationErrors)
// 	for _, e := range validationErrors {
// 		fmt.Println(e.Namespace())
// 		fmt.Println(e.Field())
// 	}
//
// 	return s.repo.Create(user)
// }

func (s *userService) FindAllUsers() (*domain.User, error) {
	return s.repo.FindAll()
}

func (s *userService) DeleteUserByID(userID uint) error {
	return s.repo.DeleteByID(userID)
}

func (s *userService) FindUserByName(name string) (*domain.User, error) {
	u, err := s.repo.FindByName(name)
	if err != nil {
		log.Info("user was not find", err)
		return nil, err
	}
	return u, nil
}

func (s *userService) UpdateUser(user *domain.User) error {
	// err := validate.Struct(*user)

	return s.repo.Update(user)
}

func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		log.Info("user wasn't found")
		return nil, err
	}
	return user, err
}
