package repository

import (
	"inkstack/internal/models"

	"gorm.io/gorm"
)

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByID(id uint) (*models.Comment, error)
	FindByPostID(postID uint) ([]models.Comment, error)
	FindByUserID(userID uint, limit, offset int) ([]models.Comment, error)
	FindReplies(parentID uint) ([]models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uint) error
	UpdateStatus(id uint, status string) error
	CountByPost(postID uint) (int64, error)
}

// commentRepository implements CommentRepository
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create creates a new comment
func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

// FindByID finds a comment by ID
func (r *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// FindByPostID retrieves all comments for a post
func (r *commentRepository) FindByPostID(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("post_id = ?", postID).
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// FindByUserID retrieves comments by user with pagination
func (r *commentRepository) FindByUserID(userID uint, limit, offset int) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}

// FindReplies retrieves replies to a comment
func (r *commentRepository) FindReplies(parentID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// Update updates a comment
func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

// Delete soft deletes a comment by ID
func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

// UpdateStatus updates the status of a comment
func (r *commentRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).
		Update("status", status).Error
}

// CountByPost returns the number of comments for a post
func (r *commentRepository) CountByPost(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}
