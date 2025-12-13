package repository

import (
	"errors"
	"time"

	"github.com/bert1727/ChatApp/internal/domain"
	"gorm.io/gorm"
)

type MessageRepository interface {
	FindMessageByID(messageID string) (*domain.Message, error)
	SaveMessage(msg domain.Message) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) FindMessageByID(messageID string) (*domain.Message, error) {
	var msg domain.Message

	if messageID == "" {
		return nil, errors.New("ID is required")
	}

	err := r.db.Find(&msg, messageID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no such message with given ID")
		} else {
			return nil, errors.New("failed to get a message")
		}
	}
	return &msg, nil
}

func (r *messageRepository) SaveMessage(msg domain.Message) error {
	model := domain.MessageModels{
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		CreatedAt:  time.Now(),
	}

	return r.db.Create(&model).Error
}
