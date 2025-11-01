package models

// User represents a user in the system
type User struct {
	BaseModel
	Email    string `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username string `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Password string `gorm:"not null;size:255" json:"-"` // Password is excluded from JSON
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}