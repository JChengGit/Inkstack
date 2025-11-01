package models

// Comment represents a user comment on a post
type Comment struct {
	BaseModel
	PostID   uint   `gorm:"not null;index" json:"post_id" validate:"required"`
	UserID   uint   `gorm:"not null;index" json:"user_id" validate:"required"`
	ParentID *uint  `gorm:"index" json:"parent_id"`
	Content  string `gorm:"type:text;not null" json:"content" validate:"required,min=1,max=1000"`
	Status   string `gorm:"type:varchar(20);not null;default:'pending';index" json:"status" validate:"oneof=pending approved rejected spam"`
}

// TableName specifies the table name for the Comment model
func (Comment) TableName() string {
	return "comments"
}
