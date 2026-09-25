package response

import "github.com/gin-gonic/gin"

func Error(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": code, "message": message})
}
func JSON(c *gin.Context, status int, value any) { c.JSON(status, value) }
