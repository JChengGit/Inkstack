package main

import (
	"context"
	"errors"
	"fmt"
	"inkstack/internal/config"
	"inkstack/internal/database"
	"inkstack/internal/handler"
	"inkstack/internal/repository"
	"inkstack/internal/service"
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

	log.Printf("Starting Inkstack API in %s mode", cfg.App.Env)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Run migrations
	// Use relative path - requires running from api/ directory
	// Alternative: Use environment variable MIGRATIONS_PATH for flexibility
	migrationsPath := "./migrations"

	log.Printf("Running database migrations from: %s", migrationsPath)
	if err := database.RunMigrations(migrationsPath); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
		log.Println("Continuing anyway - migrations may need to be run manually")
		log.Println("Note: Ensure you're running from the api/ directory")
	}

	// Set up Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	// Initialize repositories
	postRepo := repository.NewPostRepository(database.GetDB())
	commentRepo := repository.NewCommentRepository(database.GetDB())

	// Initialize services
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)

	// Initialize handlers
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	// Health check endpoint
	r.GET("/health", handler.HealthCheck)

	// Hello world endpoint
	r.GET("/hello", handler.HelloWorld)

	// API routes
	api := r.Group("/api")
	{
		// Posts routes
		posts := api.Group("/posts")
		{
			posts.GET("", postHandler.ListPosts)
			posts.POST("", postHandler.CreatePost)
			posts.GET("/:id", postHandler.GetPost)
			posts.GET("/slug/:slug", postHandler.GetPostBySlug)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
			posts.POST("/:id/publish", postHandler.PublishPost)
			posts.POST("/:id/unpublish", postHandler.UnpublishPost)

			// Comments for a specific post
			posts.GET("/:id/comments", commentHandler.ListCommentsByPost)
			posts.POST("/:id/comments", commentHandler.CreateComment)
		}

		// Comments routes
		comments := api.Group("/comments")
		{
			comments.GET("/:id", commentHandler.GetComment)
			comments.PUT("/:id", commentHandler.UpdateComment)
			comments.DELETE("/:id", commentHandler.DeleteComment)
			comments.POST("/:id/approve", commentHandler.ApproveComment)
			comments.POST("/:id/reject", commentHandler.RejectComment)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.App.Port)
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
