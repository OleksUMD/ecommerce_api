// Package utils contains some useful tools and common entities
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is base struct for every api response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

// PaginatedResponse is a base response with pagination parameters
type PaginatedResponse struct {
	Response
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta is container for the pagination parameters
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// SuccessResponse is a standard successful GET response
func SuccessResponse(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse is a standard successful POST response
func CreatedResponse(c *gin.Context, message string, data any) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessNoContentResponse is a standard successful DELETE response or any other where response data is nil
func SuccessNoContentResponse(c *gin.Context, message string, data any) {
	c.JSON(http.StatusNoContent, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse is base error response for any exception
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	c.JSON(statusCode, response)
}

// BadRequestResponse is a standard response for 400 error
func BadRequestResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusBadRequest, message, err)
}

// UnauthorizedResponse is a standard response for 401 error
func UnauthorizedResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusUnauthorized, message, err)
}

// ForbiddenResponse is a standard response for 403 error
func ForbiddenResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusForbidden, message, err)
}

// NotFoundResponse is a standard response for 404 error
func NotFoundResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusNotFound, message, err)
}

// InternalServerErrorResponse is a standard response for 500 error
func InternalServerErrorResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusInternalServerError, message, err)
}

// PaginatedSuccessResponse is a standard successful GET response with pagination params
func PaginatedSuccessResponse(c *gin.Context, message string, data any, meta PaginationMeta) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Meta: meta,
	})
}
