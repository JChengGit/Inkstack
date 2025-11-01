package service

import (
	"errors"
	"fmt"
	"inkstack/internal/models"
	"inkstack/internal/repository"

	"gorm.io/gorm"
)

// CommentService defines the interface for comment business logic
type CommentService interface {
	CreateComment(postID, userID uint, content string, parentID *uint) (*models.Comment, error)
	GetComment(id uint) (*models.Comment, error)
	ListCommentsByPost(postID uint) ([]models.Comment, error)
	ListCommentsByUser(userID uint, page, pageSize int) ([]models.Comment, int64, error)
	UpdateComment(id uint, content string) (*models.Comment, error)
	DeleteComment(id uint) error
	ApproveComment(id uint) (*models.Comment, error)
	RejectComment(id uint) (*models.Comment, error)
	MarkAsSpam(id uint) (*models.Comment, error)
}

// commentService implements CommentService
type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

// NewCommentService creates a new comment service
func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// CreateComment creates a new comment
func (s *commentService) CreateComment(postID, userID uint, content string, parentID *uint) (*models.Comment, error) {
	// Validate inputs
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(content) > 1000 {
		return nil, errors.New("content exceeds maximum length of 1000 characters")
	}
	if postID == 0 {
		return nil, errors.New("post_id is required")
	}
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	// Verify post exists
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// If replying to a comment, verify parent comment exists
	if parentID != nil && *parentID > 0 {
		parentComment, err := s.commentRepo.FindByID(*parentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent comment not found")
			}
			return nil, err
		}
		// Verify parent comment belongs to the same post
		if parentComment.PostID != postID {
			return nil, errors.New("parent comment does not belong to this post")
		}
	}

	comment := &models.Comment{
		PostID:   postID,
		UserID:   userID,
		ParentID: parentID,
		Content:  content,
		Status:   "pending",
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetComment retrieves a comment by ID
func (s *commentService) GetComment(id uint) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}
	return comment, nil
}

// ListCommentsByPost retrieves all comments for a post
func (s *commentService) ListCommentsByPost(postID uint) ([]models.Comment, error) {
	// Verify post exists
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	comments, err := s.commentRepo.FindByPostID(postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// ListCommentsByUser retrieves comments by user with pagination
func (s *commentService) ListCommentsByUser(userID uint, page, pageSize int) ([]models.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	comments, err := s.commentRepo.FindByUserID(userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Note: We'd need a CountByUser method in repository for accurate total
	// For now, returning length of results
	total := int64(len(comments))

	return comments, total, nil
}

// UpdateComment updates a comment's content
func (s *commentService) UpdateComment(id uint, content string) (*models.Comment, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(content) > 1000 {
		return nil, errors.New("content exceeds maximum length of 1000 characters")
	}

	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	comment.Content = content

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return comment, nil
}

// DeleteComment soft deletes a comment
func (s *commentService) DeleteComment(id uint) error {
	_, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("comment not found")
		}
		return err
	}

	return s.commentRepo.Delete(id)
}

// ApproveComment approves a comment
func (s *commentService) ApproveComment(id uint) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	if err := s.commentRepo.UpdateStatus(id, "approved"); err != nil {
		return nil, fmt.Errorf("failed to approve comment: %w", err)
	}

	comment.Status = "approved"
	return comment, nil
}

// RejectComment rejects a comment
func (s *commentService) RejectComment(id uint) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	if err := s.commentRepo.UpdateStatus(id, "rejected"); err != nil {
		return nil, fmt.Errorf("failed to reject comment: %w", err)
	}

	comment.Status = "rejected"
	return comment, nil
}

// MarkAsSpam marks a comment as spam
func (s *commentService) MarkAsSpam(id uint) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	if err := s.commentRepo.UpdateStatus(id, "spam"); err != nil {
		return nil, fmt.Errorf("failed to mark comment as spam: %w", err)
	}

	comment.Status = "spam"
	return comment, nil
}
