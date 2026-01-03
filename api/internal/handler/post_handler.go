package handler

import (
	"inkstack/internal/models"
	"inkstack/internal/service"
	"inkstack/internal/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PostHandler handles HTTP requests for posts
type PostHandler struct {
	service service.PostService
}

// NewPostHandler creates a new post handler
func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

// Request/Response DTOs

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
	Excerpt string `json:"excerpt"`
	Slug    string `json:"slug"`
	// AuthorID is extracted from JWT token, not from request body
}

type UpdatePostRequest struct {
	Title   *string `json:"title" binding:"omitempty,max=255"`
	Content *string `json:"content"`
	Excerpt *string `json:"excerpt"`
	Slug    *string `json:"slug"`
	Status  *string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type PostResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	AuthorID    uint       `json:"author_id"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
	ViewCount   int        `json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreatePost handles POST /api/posts
// @Summary Create a new post
// @Description Create a new blog post with title, content, and metadata
// @Tags posts
// @Accept json
// @Produce json
// @Param request body CreatePostRequest true "Post details"
// @Success 201 {object} PostResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	// Extract user ID from JWT token (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	post, err := h.service.CreatePost(req.Title, req.Content, req.Excerpt, req.Slug, userID.(uint))
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, toPostResponse(post))
}

// GetPost handles GET /api/posts/:id
// @Summary Get a post by ID
// @Description Retrieve a single post by its ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid post ID")
		return
	}

	post, err := h.service.GetPost(uint(id))
	if err != nil {
		util.RespondNotFound(c, "Post")
		return
	}

	c.JSON(http.StatusOK, toPostResponse(post))
}

// GetPostBySlug handles GET /api/posts/slug/:slug
// @Summary Get a post by slug
// @Description Retrieve a single post by its URL slug
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/posts/slug/{slug} [get]
func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		util.RespondBadRequest(c, "slug is required")
		return
	}

	post, err := h.service.GetPostBySlug(slug)
	if err != nil {
		util.RespondNotFound(c, "Post")
		return
	}

	c.JSON(http.StatusOK, toPostResponse(post))
}

// ListPosts handles GET /api/posts
// @Summary List posts
// @Description Get a paginated list of posts with optional filtering
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query string false "Filter by status (draft, published, archived)"
// @Param author_id query int false "Filter by author ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	authorID, _ := strconv.ParseUint(c.Query("author_id"), 10, 32)

	var posts []models.Post
	var total int64
	var err error

	if status == "published" {
		posts, total, err = h.service.ListPublishedPosts(page, pageSize)
	} else if authorID > 0 {
		posts, total, err = h.service.ListPostsByAuthor(uint(authorID), page, pageSize)
	} else {
		posts, total, err = h.service.ListPosts(page, pageSize)
	}

	if err != nil {
		util.RespondInternalError(c, "failed to retrieve posts")
		return
	}

	pagination := util.CalculatePagination(page, pageSize, total)
	postsResponse := toPostsResponse(posts)

	c.JSON(http.StatusOK, gin.H{
		"posts":      postsResponse,
		"pagination": pagination,
	})
}

// UpdatePost handles PUT /api/posts/:id
// @Summary Update a post
// @Description Update an existing post's information
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body UpdatePostRequest true "Updated post details"
// @Success 200 {object} PostResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid post ID")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Excerpt != nil {
		updates["excerpt"] = *req.Excerpt
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	post, err := h.service.UpdatePost(uint(id), updates)
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, toPostResponse(post))
}

// DeletePost handles DELETE /api/posts/:id
// @Summary Delete a post
// @Description Soft delete a post by ID
// @Tags posts
// @Param id path int true "Post ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid post ID")
		return
	}

	if err := h.service.DeletePost(uint(id)); err != nil {
		util.RespondNotFound(c, "Post")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// PublishPost handles POST /api/posts/:id/publish
// @Summary Publish a post
// @Description Change post status to published
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} PostResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/posts/{id}/publish [post]
func (h *PostHandler) PublishPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid post ID")
		return
	}

	post, err := h.service.PublishPost(uint(id))
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, toPostResponse(post))
}

// UnpublishPost handles POST /api/posts/:id/unpublish
// @Summary Unpublish a post
// @Description Change post status to draft
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} PostResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/posts/{id}/unpublish [post]
func (h *PostHandler) UnpublishPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid post ID")
		return
	}

	post, err := h.service.UnpublishPost(uint(id))
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, toPostResponse(post))
}

// Helper functions

func toPostResponse(post *models.Post) PostResponse {
	return PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Content:     post.Content,
		Excerpt:     post.Excerpt,
		AuthorID:    post.AuthorID,
		Status:      post.Status,
		PublishedAt: post.PublishedAt,
		ViewCount:   post.ViewCount,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}

func toPostsResponse(posts []models.Post) []PostResponse {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = toPostResponse(&post)
	}
	return responses
}
