// NOTE: change the directory of this file
package ws

import "gorm.io/gorm"

type MessageService interface{}

type messageService struct {
	db  *gorm.DB
	hub *Hub
}

// func NewMessageService(db *gorm.DB) *MessageService
// func (s *messageService) AttachHub(h *Hub)
// func (s *messageService) ProcessIncoming(raw []byte) error
