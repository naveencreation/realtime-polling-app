package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"polling-backend/pkg/response"
)

func RequireAuth(service Service, cookie string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := c.Cookie(cookie)
		if err != nil || raw == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				raw = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
				err = nil
			}
		}
		if err != nil || raw == "" {
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
