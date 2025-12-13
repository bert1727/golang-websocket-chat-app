package domain

import "time"

type RefreshToken struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID string `gorm:"index;not null" json:"user_id"`

	TokenHash string `gorm:"uniqueIndex;not null" json:"token_hash"`
	TokenJTI  string `gorm:"uniqueIndex;not null" json:"jti"`

	IssuedAt  time.Time `gorm:"not null;index" json:"issued_at"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`

	IsRevoked bool `gorm:"default:false;index" json:"is_revoked"`
	IsUsed    bool `gorm:"default:false" json:"is_used"`

	// IPAddress string `gorm:"size:50" json:"ip_address"`
	// UserAgent string `gorm:"size:500" json:"user_agent"`
	// DeviceID  string `gorm:"size:100" json:"device_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
