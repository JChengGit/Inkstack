package models

import (
	"time"
)

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	BaseModel
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null;size:500" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsRevoked bool      `gorm:"default:false" json:"is_revoked"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
}

// TableName specifies the table name for RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsValid checks if the token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && time.Now().Before(rt.ExpiresAt)
}
