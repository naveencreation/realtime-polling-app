package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"polling-backend/pkg/response"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered", "path", c.Request.URL.Path)
		response.Error(c, 500, "server_error", "something went wrong")
	})
}
