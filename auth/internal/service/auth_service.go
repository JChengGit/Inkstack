package service

import (
	"context"
	"fmt"
	"inkstack-auth/internal/database"
	"inkstack-auth/internal/models"
	"inkstack-auth/internal/repository"
	"inkstack-auth/internal/util"
	"time"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	jwtService *JWTService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtService *JWTService,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
	}
}

// RegisterInput contains registration data
type RegisterInput struct {
	Email    string
	Username string
	Password string
}

// LoginInput contains login data
type LoginInput struct {
	EmailOrUsername string
	Password        string
	IPAddress       string
	UserAgent       string
}

// TokenPair contains access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*models.User, *TokenPair, error) {
	// Validate input
	if err := util.ValidateEmail(input.Email); err != nil {
		return nil, nil, err
	}

	if err := util.ValidateUsername(input.Username); err != nil {
		return nil, nil, err
	}

	if err := util.ValidatePasswordStrength(input.Password); err != nil {
		return nil, nil, err
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(input.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, nil, fmt.Errorf("email already registered")
	}

	// Check if username already exists
	exists, err = s.userRepo.ExistsByUsername(input.Username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, nil, fmt.Errorf("username already taken")
	}

	// Hash password
	passwordHash, err := util.HashPassword(input.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: passwordHash,
		IsActive:     true,
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user, "", "")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, input LoginInput) (*models.User, *TokenPair, error) {
	// Check rate limiting
	attempts, err := database.GetLoginAttempts(ctx, input.EmailOrUsername)
	if err == nil && attempts >= 5 {
		return nil, nil, fmt.Errorf("too many failed login attempts, please try again in 15 minutes")
	}

	// Find user
	user, err := s.userRepo.FindByEmailOrUsername(input.EmailOrUsername)
	if err != nil {
		// Increment failed attempts
		database.IncrementLoginAttempts(ctx, input.EmailOrUsername)
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, fmt.Errorf("account is inactive")
	}

	// Verify password
	if !util.ComparePassword(user.PasswordHash, input.Password) {
		// Increment failed attempts
		database.IncrementLoginAttempts(ctx, input.EmailOrUsername)
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Reset login attempts on successful login
	database.ResetLoginAttempts(ctx, input.EmailOrUsername)

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user, input.IPAddress, input.UserAgent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, error) {
	// Validate refresh token JWT
	claims, err := s.jwtService.ValidateToken(refreshTokenString)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Check if token exists in database and is not revoked
	tokenRecord, err := s.tokenRepo.FindByToken(refreshTokenString)
	if err != nil {
		return "", fmt.Errorf("refresh token not found")
	}

	if !tokenRecord.IsValid() {
		return "", fmt.Errorf("refresh token is invalid or expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return "", fmt.Errorf("account is inactive")
	}

	// Generate new access token
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Optional: Implement token rotation (generate new refresh token)
	// For simplicity, we're not rotating refresh tokens here

	return accessToken, nil
}

// Logout revokes a refresh token
func (s *AuthService) Logout(ctx context.Context, refreshTokenString, accessTokenString string) error {
	// Revoke refresh token in database
	if err := s.tokenRepo.RevokeToken(refreshTokenString); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	// Blacklist access token in Redis
	// Get token expiry to set Redis TTL
	expiry, err := s.jwtService.GetTokenExpiry(accessTokenString)
	if err == nil {
		ttl := time.Until(expiry)
		if ttl > 0 {
			database.BlacklistToken(ctx, accessTokenString, ttl)
		}
	}

	return nil
}

// LogoutAll revokes all refresh tokens for a user
func (s *AuthService) LogoutAll(ctx context.Context, userID uint) error {
	return s.tokenRepo.RevokeAllUserTokens(userID)
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if !util.ComparePassword(user.PasswordHash, oldPassword) {
		return fmt.Errorf("invalid current password")
	}

	// Validate new password strength
	if err := util.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	passwordHash, err := util.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = passwordHash
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all existing refresh tokens (force re-login)
	s.tokenRepo.RevokeAllUserTokens(userID)

	return nil
}

// ValidateToken validates an access token and returns user info
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Check if token is blacklisted
	isBlacklisted, err := database.IsTokenBlacklisted(ctx, tokenString)
	if err == nil && isBlacklisted {
		return nil, fmt.Errorf("token has been revoked")
	}

	// Validate token
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	return user, nil
}

// generateTokenPair generates both access and refresh tokens
func (s *AuthService) generateTokenPair(ctx context.Context, user *models.User, ipAddress, userAgent string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, expiresAt, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	tokenRecord := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
		IsRevoked: false,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	if err := s.tokenRepo.Create(tokenRecord); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
