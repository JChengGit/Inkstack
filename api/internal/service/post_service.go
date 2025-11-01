package service

import (
	"errors"
	"fmt"
	"inkstack/internal/models"
	"inkstack/internal/repository"
	"inkstack/internal/util"
	"time"

	"gorm.io/gorm"
)

// PostService defines the interface for post business logic
type PostService interface {
	CreatePost(title, content, excerpt, slug string, authorID uint) (*models.Post, error)
	GetPost(id uint) (*models.Post, error)
	GetPostBySlug(slug string) (*models.Post, error)
	ListPosts(page, pageSize int) ([]models.Post, int64, error)
	ListPostsByAuthor(authorID uint, page, pageSize int) ([]models.Post, int64, error)
	ListPublishedPosts(page, pageSize int) ([]models.Post, int64, error)
	UpdatePost(id uint, updates map[string]interface{}) (*models.Post, error)
	DeletePost(id uint) error
	PublishPost(id uint) (*models.Post, error)
	UnpublishPost(id uint) (*models.Post, error)
	GenerateSlug(title string) string
}

// postService implements PostService
type postService struct {
	repo repository.PostRepository
}

// NewPostService creates a new post service
func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

// CreatePost creates a new post
func (s *postService) CreatePost(title, content, excerpt, slug string, authorID uint) (*models.Post, error) {
	// Validate inputs
	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if authorID == 0 {
		return nil, errors.New("author_id is required")
	}

	// Generate slug if not provided
	if slug == "" {
		slug = util.GenerateSlug(title)
	} else if !util.IsValidSlug(slug) {
		return nil, errors.New("invalid slug format")
	}

	// Check if slug already exists
	existingPost, _ := s.repo.FindBySlug(slug)
	if existingPost != nil {
		return nil, errors.New("slug already exists")
	}

	post := &models.Post{
		Title:     title,
		Slug:      slug,
		Content:   content,
		Excerpt:   excerpt,
		AuthorID:  authorID,
		Status:    "draft",
		ViewCount: 0,
	}

	if err := s.repo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetPost retrieves a post by ID and increments view count
func (s *postService) GetPost(id uint) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// Increment view count (ignore errors)
	_ = s.repo.IncrementViewCount(id)

	return post, nil
}

// GetPostBySlug retrieves a post by slug and increments view count
func (s *postService) GetPostBySlug(slug string) (*models.Post, error) {
	post, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// Increment view count (ignore errors)
	_ = s.repo.IncrementViewCount(post.ID)

	return post, nil
}

// ListPosts retrieves all posts with pagination
func (s *postService) ListPosts(page, pageSize int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	posts, err := s.repo.FindAll(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// ListPostsByAuthor retrieves posts by author with pagination
func (s *postService) ListPostsByAuthor(authorID uint, page, pageSize int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	posts, err := s.repo.FindByAuthor(authorID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByAuthor(authorID)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// ListPublishedPosts retrieves published posts with pagination
func (s *postService) ListPublishedPosts(page, pageSize int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	posts, err := s.repo.FindByStatus("published", pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// For published posts count, we'd need a CountByStatus method
	// For now, using regular Count
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// UpdatePost updates a post
func (s *postService) UpdatePost(id uint, updates map[string]interface{}) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// Apply updates
	if title, ok := updates["title"].(string); ok && title != "" {
		post.Title = title
	}
	if content, ok := updates["content"].(string); ok {
		post.Content = content
	}
	if excerpt, ok := updates["excerpt"].(string); ok {
		post.Excerpt = excerpt
	}
	if slug, ok := updates["slug"].(string); ok && slug != "" {
		if !util.IsValidSlug(slug) {
			return nil, errors.New("invalid slug format")
		}
		// Check slug uniqueness if changing
		if slug != post.Slug {
			existingPost, _ := s.repo.FindBySlug(slug)
			if existingPost != nil && existingPost.ID != id {
				return nil, errors.New("slug already exists")
			}
		}
		post.Slug = slug
	}
	if status, ok := updates["status"].(string); ok {
		if status != "draft" && status != "published" && status != "archived" {
			return nil, errors.New("invalid status")
		}
		post.Status = status
	}

	if err := s.repo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

// DeletePost soft deletes a post
func (s *postService) DeletePost(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

// PublishPost publishes a post
func (s *postService) PublishPost(id uint) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	if post.Status == "published" {
		return post, nil // Already published
	}

	post.Status = "published"
	now := time.Now()
	post.PublishedAt = &now

	if err := s.repo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to publish post: %w", err)
	}

	return post, nil
}

// UnpublishPost unpublishes a post (sets to draft)
func (s *postService) UnpublishPost(id uint) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	post.Status = "draft"
	post.PublishedAt = nil

	if err := s.repo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to unpublish post: %w", err)
	}

	return post, nil
}

// GenerateSlug generates a slug from a title
func (s *postService) GenerateSlug(title string) string {
	return util.GenerateSlug(title)
}
