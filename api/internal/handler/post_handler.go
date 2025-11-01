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
	Title    string `json:"title" binding:"required,max=255"`
	Content  string `json:"content" binding:"required"`
	Excerpt  string `json:"excerpt"`
	Slug     string `json:"slug"`
	AuthorID uint   `json:"author_id" binding:"required"`
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
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	post, err := h.service.CreatePost(req.Title, req.Content, req.Excerpt, req.Slug, req.AuthorID)
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, toPostResponse(post))
}

// GetPost handles GET /api/posts/:id
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
