package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response structure
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, code int, message string, detail ...string) {
	response := ErrorResponse{
		Error: message,
	}
	if len(detail) > 0 {
		response.Message = detail[0]
	}
	c.JSON(code, response)
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// RespondWithData sends a data response without wrapper
func RespondWithData(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

// RespondNotFound sends a 404 response
func RespondNotFound(c *gin.Context, resource string) {
	RespondWithError(c, http.StatusNotFound, resource+" not found")
}

// RespondBadRequest sends a 400 response
func RespondBadRequest(c *gin.Context, message string) {
	RespondWithError(c, http.StatusBadRequest, message)
}

// RespondInternalError sends a 500 response
func RespondInternalError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusInternalServerError, message)
}

// CalculatePagination calculates pagination values
func CalculatePagination(page, pageSize int, total int64) PaginationResponse {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return PaginationResponse{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
