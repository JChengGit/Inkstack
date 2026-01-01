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

// CommentHandler handles HTTP requests for comments
type CommentHandler struct {
	service service.CommentService
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(service service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

// Request/Response DTOs

type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

type CommentResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	ParentID  *uint     `json:"parent_id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateComment handles POST /api/posts/:id/comments
// @Summary Create a new comment
// @Description Add a comment to a post, optionally as a reply to another comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body CreateCommentRequest true "Comment details"
// @Success 201 {object} CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/posts/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid post ID")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	// TODO: Get userID from auth context
	// For now, using a hardcoded value
	userID := uint(1)

	comment, err := h.service.CreateComment(uint(postID), userID, req.Content, req.ParentID)
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, toCommentResponse(comment))
}

// GetComment handles GET /api/comments/:id
// @Summary Get a comment by ID
// @Description Retrieve a single comment by its ID
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid comment ID")
		return
	}

	comment, err := h.service.GetComment(uint(id))
	if err != nil {
		util.RespondNotFound(c, "Comment")
		return
	}

	c.JSON(http.StatusOK, toCommentResponse(comment))
}

// ListCommentsByPost handles GET /api/posts/:id/comments
// @Summary List comments for a post
// @Description Get all comments for a specific post
// @Tags comments
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/posts/{id}/comments [get]
func (h *CommentHandler) ListCommentsByPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid post ID")
		return
	}

	comments, err := h.service.ListCommentsByPost(uint(postID))
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": toCommentsResponse(comments),
	})
}

// UpdateComment handles PUT /api/comments/:id
// @Summary Update a comment
// @Description Update a comment's content
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param request body UpdateCommentRequest true "Updated comment content"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid comment ID")
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	comment, err := h.service.UpdateComment(uint(id), req.Content)
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, toCommentResponse(comment))
}

// DeleteComment handles DELETE /api/comments/:id
// @Summary Delete a comment
// @Description Soft delete a comment by ID
// @Tags comments
// @Param id path int true "Comment ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid comment ID")
		return
	}

	if err := h.service.DeleteComment(uint(id)); err != nil {
		util.RespondNotFound(c, "Comment")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ApproveComment handles POST /api/comments/:id/approve
// @Summary Approve a comment
// @Description Change comment status to approved
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/comments/{id}/approve [post]
func (h *CommentHandler) ApproveComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid comment ID")
		return
	}

	comment, err := h.service.ApproveComment(uint(id))
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, toCommentResponse(comment))
}

// RejectComment handles POST /api/comments/:id/reject
// @Summary Reject a comment
// @Description Change comment status to rejected
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/comments/{id}/reject [post]
func (h *CommentHandler) RejectComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.RespondBadRequest(c, "invalid comment ID")
		return
	}

	comment, err := h.service.RejectComment(uint(id))
	if err != nil {
		util.RespondBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, toCommentResponse(comment))
}

// Helper functions

func toCommentResponse(comment *models.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		Status:    comment.Status,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func toCommentsResponse(comments []models.Comment) []CommentResponse {
	responses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = toCommentResponse(&comment)
	}
	return responses
}
