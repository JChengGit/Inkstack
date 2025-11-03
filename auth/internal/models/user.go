package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	BaseModel
	Email         string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username      string     `gorm:"uniqueIndex;not null;size:50" json:"username"`
	PasswordHash  string     `gorm:"not null;size:255" json:"-"` // Never expose in JSON
	DisplayName   string     `gorm:"size:100" json:"display_name"`
	Bio           string     `gorm:"type:text" json:"bio"`
	AvatarURL     string     `gorm:"size:500" json:"avatar_url"`
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	Role          string     `gorm:"default:'user';size:20" json:"role"` // user, admin
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// PublicUser returns user data safe for public consumption
type PublicUser struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
}

// ToPublic converts User to PublicUser
func (u *User) ToPublic() PublicUser {
	return PublicUser{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		AvatarURL:   u.AvatarURL,
		Role:        u.Role,
	}
}
