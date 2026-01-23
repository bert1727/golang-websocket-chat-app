package domain

import (
	"time"
)

type MessageModels struct {
	ID         uint   `gorm:"primaryKey"`
	SenderID   string // user who sent
	ReceiverID uint   // user who receives
	CreatedAt  time.Time
	Content    string
	// gorm.Model
}

// DTO
// FIX: implement all relations in db
type Message struct {
	ID   uint `gorm:"primaryKey"`
	Type string

	SenderID uint `gorm:"not null" validate:"required"`
	Sender   User `json:"-" validate:"-" gorm:"foreignKey:SenderID;references:ID" `

	ReceiverID uint `gorm:"not null" validate:"required"`
	Receiver   User `json:"-" validate:"-" gorm:"foreignKey:ReceiverID;references:ID"`

	Content   string
	Timestamp time.Time
}
