package repository

import (
	"inkstack/internal/models"

	"gorm.io/gorm"
)

// PostRepository defines the interface for post data operations
type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id uint) (*models.Post, error)
	FindBySlug(slug string) (*models.Post, error)
	FindAll(limit, offset int) ([]models.Post, error)
	FindByAuthor(authorID uint, limit, offset int) ([]models.Post, error)
	FindByStatus(status string, limit, offset int) ([]models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
	IncrementViewCount(id uint) error
	Count() (int64, error)
	CountByAuthor(authorID uint) (int64, error)
}

// postRepository implements PostRepository
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// Create creates a new post
func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// FindByID finds a post by ID
func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// FindBySlug finds a post by slug
func (r *postRepository) FindBySlug(slug string) (*models.Post, error) {
	var post models.Post
	err := r.db.Where("slug = ?", slug).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// FindAll retrieves all posts with pagination
func (r *postRepository) FindAll(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// FindByAuthor retrieves posts by author with pagination
func (r *postRepository) FindByAuthor(authorID uint, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Where("author_id = ?", authorID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// FindByStatus retrieves posts by status with pagination
func (r *postRepository) FindByStatus(status string, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Where("status = ?", status).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// Update updates a post
func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// Delete soft deletes a post by ID
func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

// IncrementViewCount increments the view count of a post
func (r *postRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// Count returns the total number of posts
func (r *postRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Count(&count).Error
	return count, err
}

// CountByAuthor returns the number of posts by an author
func (r *postRepository) CountByAuthor(authorID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Where("author_id = ?", authorID).Count(&count).Error
	return count, err
}
