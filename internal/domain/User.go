package domain

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Password  string
	Email     string `gorm:"unique"`
	Online    bool
	Messages  []Message `json:"messages" gorm:"foreignKey:UserID"`
	UpdatedAt time.Time
}

// `gorm:"primaryKey"`

type MessagePayload struct {
	SenderID   uint
	ReceiverID *uint
	Content    string
}
