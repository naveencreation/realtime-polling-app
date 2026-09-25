package auth

import (
	"github.com/gin-gonic/gin"
	"polling-backend/pkg/response"
)

func RequireAuth(service Service, cookie string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := c.Cookie(cookie)
		if err != nil {
			response.Error(c, 401, "unauthorized", "authentication required")
			return
		}
		id, err := service.VerifyToken(raw)
		if err != nil {
			response.Error(c, 401, "unauthorized", "authentication required")
			return
		}
		c.Set("userId", id)
		c.Next()
	}
}
