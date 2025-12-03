package service

import (
	"github.com/bert1727/ChatApp/internal/models"
	"github.com/bert1727/ChatApp/internal/repository"
)

type MessageService interface {
	SaveMessage(message models.Message) error
}

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (s *messageService) SaveMessage(msg models.Message) error {
	return s.repo.SaveMessage(msg)
}
