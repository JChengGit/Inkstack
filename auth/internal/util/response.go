package util

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondBadRequest responds with a 400 Bad Request
func RespondBadRequest(c *gin.Context, message string) {
	c.JSON(400, ErrorResponse{
		Error:   "bad_request",
		Message: message,
	})
}

// RespondUnauthorized responds with a 401 Unauthorized
func RespondUnauthorized(c *gin.Context, message string) {
	c.JSON(401, ErrorResponse{
		Error:   "unauthorized",
		Message: message,
	})
}

// RespondForbidden responds with a 403 Forbidden
func RespondForbidden(c *gin.Context, message string) {
	c.JSON(403, ErrorResponse{
		Error:   "forbidden",
		Message: message,
	})
}

// RespondNotFound responds with a 404 Not Found
func RespondNotFound(c *gin.Context, resource string) {
	c.JSON(404, ErrorResponse{
		Error:   "not_found",
		Message: resource + " not found",
	})
}

// RespondConflict responds with a 409 Conflict
func RespondConflict(c *gin.Context, message string) {
	c.JSON(409, ErrorResponse{
		Error:   "conflict",
		Message: message,
	})
}

// RespondTooManyRequests responds with a 429 Too Many Requests
func RespondTooManyRequests(c *gin.Context, message string) {
	c.JSON(429, ErrorResponse{
		Error:   "too_many_requests",
		Message: message,
	})
}

// RespondInternalError responds with a 500 Internal Server Error
func RespondInternalError(c *gin.Context, message string) {
	c.JSON(500, ErrorResponse{
		Error:   "internal_error",
		Message: message,
	})
}

// RespondSuccess responds with a 200 OK and optional data
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Message: message,
		Data:    data,
	}
	c.JSON(200, response)
}

// RespondCreated responds with a 201 Created
func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(201, data)
}

// RespondNoContent responds with a 204 No Content
func RespondNoContent(c *gin.Context) {
	c.Status(204)
}
