package handler

import (
	"inkstack-auth/internal/service"
	"inkstack-auth/internal/util"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest represents registration request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest represents login request body
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents refresh token request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest represents logout request body
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest represents change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         interface{} `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

// Register handles POST /api/auth/register
// @Summary Register a new user
// @Description Create a new user account with email, username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	user, tokens, err := h.authService.Register(c.Request.Context(), service.RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	util.RespondCreated(c, AuthResponse{
		User:         user.ToPublic(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Login handles POST /api/auth/login
// @Summary User login
// @Description Authenticate user with email/username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	user, tokens, err := h.authService.Login(c.Request.Context(), service.LoginInput{
		EmailOrUsername: req.EmailOrUsername,
		Password:        req.Password,
		IPAddress:       c.ClientIP(),
		UserAgent:       c.Request.UserAgent(),
	})

	if err != nil {
		util.RespondUnauthorized(c, err.Error())
		return
	}

	util.RespondSuccess(c, "Login successful", AuthResponse{
		User:         user.ToPublic(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// RefreshToken handles POST /api/auth/refresh
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	accessToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		util.RespondUnauthorized(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"access_token": accessToken,
	})
}

// Logout handles POST /api/auth/logout
// @Summary Logout user
// @Description Invalidate refresh token and access token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LogoutRequest true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	// Get access token from header
	accessToken := c.GetHeader("Authorization")
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken, accessToken); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	util.RespondSuccess(c, "Logged out successfully", nil)
}

// GetMe handles GET /api/auth/me
// @Summary Get current user profile
// @Description Get the authenticated user's profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	_, exists := c.Get("user_id")
	if !exists {
		util.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Get full user info from token validation
	accessToken := c.GetHeader("Authorization")
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	user, err := h.authService.ValidateToken(c.Request.Context(), accessToken)
	if err != nil {
		util.RespondUnauthorized(c, "Invalid token")
		return
	}

	c.JSON(200, gin.H{
		"user": user.ToPublic(),
	})
}

// ChangePassword handles POST /api/auth/change-password
// @Summary Change password
// @Description Change the authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Old and new password"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		util.RespondUnauthorized(c, "User not authenticated")
		return
	}

	if err := h.authService.ChangePassword(
		c.Request.Context(),
		userID.(uint),
		req.OldPassword,
		req.NewPassword,
	); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	util.RespondSuccess(c, "Password changed successfully", nil)
}

// ValidateToken handles POST /api/auth/validate (for API service)
// @Summary Validate JWT token
// @Description Validate a JWT token and return user info (used by API service)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Token to validate"
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/validate [post]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	type ValidateRequest struct {
		Token string `json:"token" binding:"required"`
	}

	var req ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	user, err := h.authService.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(200, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"valid":   true,
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	})
}

// HealthCheck handles GET /health
// @Summary Health check
// @Description Check if the auth service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "auth",
	})
}
