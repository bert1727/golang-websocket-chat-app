package models

import "time"

type MessageModels struct {
	ID         uint   `gorm:"primaryKey"`
	SenderID   string // user who sent
	ReceiverID uint   // user who receives (или RoomID для чата-комнаты)
	Content    string
	CreatedAt  time.Time
}

// DTO
type Message struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	SenderID   string
	ReceiverID uint
	Content    string
	Timestamp  time.Time
}
