package middleware

import (
	"inkstack-auth/internal/service"
	"inkstack-auth/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates middleware that validates JWT tokens
func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.RespondUnauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			util.RespondUnauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			util.RespondUnauthorized(c, "Token is required")
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			util.RespondUnauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user info in context for handlers
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireRole creates middleware that checks user role
func RequireRole(role string, jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			util.RespondForbidden(c, "Role information not found")
			c.Abort()
			return
		}

		if userRole.(string) != role {
			util.RespondForbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth is middleware that extracts user info if token exists, but doesn't require it
func OptionalAuth(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtService.ValidateToken(token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("email", claims.Email)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
			}
		}

		c.Next()
	}
}
