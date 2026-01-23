package service

import (
	"fmt"

	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/rs/zerolog/log"
)

type MessageService interface {
	SaveMessage(message *domain.Message) error
}

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (s *messageService) SaveMessage(msg *domain.Message) error {
	err := s.repo.SaveMessage(msg)
	if err != nil {
		log.Err(err).Msg("failed to save a message")
		return fmt.Errorf("failed to save a message")
	}

	log.Info().Msg("successfully saved a message")
	return err
}
