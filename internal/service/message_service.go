package service

import (
	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/repository"
)

type MessageService interface {
	SaveMessage(message domain.Message) error
}

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (s *messageService) SaveMessage(msg domain.Message) error {
	return s.repo.SaveMessage(msg)
}
