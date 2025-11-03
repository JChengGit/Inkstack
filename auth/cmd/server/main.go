package main

import (
	"context"
	"errors"
	"fmt"
	"inkstack-auth/internal/config"
	"inkstack-auth/internal/database"
	"inkstack-auth/internal/handler"
	"inkstack-auth/internal/middleware"
	"inkstack-auth/internal/repository"
	"inkstack-auth/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	log.Printf("Starting Inkstack Auth Service in %s mode", cfg.App.Env)

	// Connect to PostgreSQL
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Connect to Redis
	if err := database.ConnectRedis(cfg); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer func() {
		if err := database.CloseRedis(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}()

	// Run migrations (optional - you can also use SQL migrations manually)
	if err := database.RunMigrations(); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	// Set up Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	// Initialize dependencies
	db := database.GetDB()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Services
	jwtService := service.NewJWTService(cfg)
	authService := service.NewAuthService(userRepo, tokenRepo, jwtService)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)

	// Health check endpoint
	r.GET("/health", handler.HealthCheck)

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken) // For API service

			// Protected routes (require authentication)
			protected := auth.Group("")
			protected.Use(middleware.AuthMiddleware(jwtService))
			{
				protected.GET("/me", authHandler.GetMe)
				protected.POST("/logout", authHandler.Logout)
				protected.POST("/change-password", authHandler.ChangePassword)
			}
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Auth service starting on port %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
