package repository

import (
	"fmt"
	"inkstack-auth/internal/models"
	"time"

	"gorm.io/gorm"
)

// TokenRepository defines the interface for refresh token operations
type TokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	FindByUserID(userID uint) ([]models.RefreshToken, error)
	RevokeToken(token string) error
	RevokeAllUserTokens(userID uint) error
	DeleteExpired() error
	CleanupRevokedTokens(olderThan time.Duration) error
}

type tokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

// Create creates a new refresh token
func (r *tokenRepository) Create(token *models.RefreshToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

// FindByToken finds a refresh token by token string
func (r *tokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}
	return &refreshToken, nil
}

// FindByUserID finds all refresh tokens for a user
func (r *tokenRepository) FindByUserID(userID uint) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	if err := r.db.Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to find refresh tokens: %w", err)
	}
	return tokens, nil
}

// RevokeToken revokes a specific refresh token
func (r *tokenRepository) RevokeToken(token string) error {
	if err := r.db.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error; err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (r *tokenRepository) RevokeAllUserTokens(userID uint) error {
	if err := r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND is_revoked = false", userID).
		Update("is_revoked", true).Error; err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}
	return nil
}

// DeleteExpired deletes all expired tokens
func (r *tokenRepository) DeleteExpired() error {
	if err := r.db.Where("expires_at < ?", time.Now()).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}
	return nil
}

// CleanupRevokedTokens deletes revoked tokens older than specified duration
func (r *tokenRepository) CleanupRevokedTokens(olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)
	if err := r.db.Where("is_revoked = true AND updated_at < ?", cutoffTime).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup revoked tokens: %w", err)
	}
	return nil
}
