package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"polling-backend/pkg/response"
)

// Recovery returns a Gin middleware that recovers from panics and logs structured details.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(ginCtx *gin.Context, recoveredErr any) {
		logger.Error("panic recovered",
			"error", recoveredErr,
			"path", ginCtx.Request.URL.Path,
			"stack", string(debug.Stack()),
		)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
	})
}
