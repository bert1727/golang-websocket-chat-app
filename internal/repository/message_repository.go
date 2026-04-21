package repository

import (
	"github.com/bert1727/ChatApp/internal/config"
	"github.com/bert1727/ChatApp/internal/domain"
	"gorm.io/gorm"
)

type MessageRepository interface {
	FindMessageByID(messageID string) (*domain.Message, error)
	SaveMessage(msg *domain.Message) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository() MessageRepository {
	return &messageRepository{
		db: config.New().DB,
	}
}

func (r *messageRepository) FindMessageByID(messageID string) (*domain.Message, error) {
	var msg domain.Message

	if messageID == "" {
		return nil, ErrIDRequired
	}

	err := r.db.Find(&msg, messageID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrMessageNotFound
			// return nil, errors.New("no such message with given ID")
		} else {
			return nil, ErrInternalServer
			// return nil, errors.New("failed to get a message")
		}
	}
	return &msg, nil
}

func (r *messageRepository) SaveMessage(msg *domain.Message) error {
	err := r.db.Create(msg).Error
	if err != nil {
		return ErrInternalServer
	}
	return err
}
