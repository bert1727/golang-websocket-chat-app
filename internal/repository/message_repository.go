package repository

import (
	"time"

	"github.com/bert1727/ChatApp/internal/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	FindMessageByID(messageID string) (string, error)
	SaveMessage(msg models.Message) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) FindMessageByID(messageID string) (string, error) {
	r.db.First(&messageID)
	return messageID, nil
}

func (r *messageRepository) SaveMessage(msg models.Message) error {
	model := models.MessageModels{
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		CreatedAt:  time.Now(),
	}
	return r.db.Create(&model).Error
}
