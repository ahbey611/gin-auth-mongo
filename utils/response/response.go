package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 400 BadRequest
func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Bad Request",
		"status":  http.StatusBadRequest,
	})
}

// 400 BadRequestWithMessage
func BadRequestWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": message,
		"status":  http.StatusBadRequest,
	})
}

// custom Failure
func Failure(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"message": message,
		"status":  statusCode,
	})
}

// 401 Unauthorized
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Unauthorized",
		"status":  http.StatusUnauthorized,
	})
}

// 403 Forbidden
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"message": "Forbidden",
		"status":  http.StatusForbidden,
	})
}

func PermissionDenied(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"message": "Permission denied",
		"status":  http.StatusForbidden,
	})
}

// 500 InternalServerError
func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": "Internal Server Error",
		"status":  http.StatusInternalServerError,
	})
}

// 200 Success
func Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"status":  http.StatusOK,
	})
}

// SuccessWithData
func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"status":  http.StatusOK,
		"data":    data,
	})
}
