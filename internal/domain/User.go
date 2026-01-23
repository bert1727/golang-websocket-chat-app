package domain

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `validate:"required,min=2,max=50"`

	Password string `validate:"required,min=8"`

	Email            string `gorm:"unique" validate:"required,email"`
	Online           bool
	SentMessages     []Message `gorm:"foreignKey:SenderID;references:ID" json:"sent_messages,omitempty"`
	ReceivedMessages []Message `gorm:"foreignKey:ReceiverID;references:ID" json:"received_messages,omitempty"`
	UpdatedAt        time.Time

	// gorm.Model //NOTE: for soft delete
}

// NOTE: Role: admin, user
