package models

import "time"

// Post represents a blog post
type Post struct {
	BaseModel
	Title       string     `gorm:"type:varchar(255);not null" json:"title" validate:"required,max=255"`
	Slug        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug" validate:"required,max=255"`
	Content     string     `gorm:"type:text;not null" json:"content" validate:"required"`
	Excerpt     string     `gorm:"type:text" json:"excerpt"`
	AuthorID    uint       `gorm:"not null;index" json:"author_id" validate:"required"`
	Status      string     `gorm:"type:varchar(20);not null;default:'draft';index" json:"status" validate:"oneof=draft published archived"`
	PublishedAt *time.Time `gorm:"index" json:"published_at"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
}

// TableName specifies the table name for the Post model
func (Post) TableName() string {
	return "posts"
}
