package handler

import (
	"inkstack/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HelloWorld handles the hello world endpoint
func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
		"service": "Inkstack API",
	})
}

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	// Check database health
	dbHealth, err := database.CheckHealth()

	response := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"database":  dbHealth,
	}

	// If database is unhealthy, return 503
	if err != nil || dbHealth.Status != "connected" {
		response["status"] = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
