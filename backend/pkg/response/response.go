package response

import "github.com/gin-gonic/gin"

// Error writes a standardized JSON error response and aborts the request.
func Error(ginCtx *gin.Context, statusCode int, errorCode, message string) {
	ginCtx.AbortWithStatusJSON(statusCode, gin.H{
		"error":   errorCode,
		"message": message,
	})
}

// JSON writes a standardized JSON success response.
func JSON(ginCtx *gin.Context, statusCode int, payload any) {
	ginCtx.JSON(statusCode, payload)
}
