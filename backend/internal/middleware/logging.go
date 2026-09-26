package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the HTTP header used to carry correlation IDs.
	RequestIDHeader = "X-Request-ID"
	loggerKey       = "slog_logger"
)

// RequestLogger generates or extracts a request correlation ID, injects a scoped
// logger into the Gin context, and logs request lifecycle metrics upon completion.
func RequestLogger(baseLogger *slog.Logger) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		requestID := strings.TrimSpace(ginCtx.GetHeader(RequestIDHeader))
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ginCtx.Header(RequestIDHeader, requestID)

		// Create a child logger bound to this specific request ID
		reqLogger := baseLogger.With("request_id", requestID)
		ginCtx.Set(loggerKey, reqLogger)

		startTime := time.Now()
		path := ginCtx.Request.URL.Path

		ginCtx.Next()

		// Skip access logging for successful health check probes to avoid log noise
		statusCode := ginCtx.Writer.Status()
		if path == "/healthz" && statusCode == http.StatusOK {
			return
		}

		duration := time.Since(startTime)
		durationMs := float64(duration.Microseconds()) / 1000.0

		fields := []any{
			"method", ginCtx.Request.Method,
			"path", path,
			"status", statusCode,
			"duration_ms", durationMs,
			"client_ip", ginCtx.ClientIP(),
			"bytes", ginCtx.Writer.Size(),
		}

		switch {
		case statusCode >= http.StatusInternalServerError:
			reqLogger.Error("http request completed with server error", fields...)
		case statusCode >= http.StatusBadRequest:
			reqLogger.Warn("http request completed with client error", fields...)
		default:
			reqLogger.Info("http request completed", fields...)
		}
	}
}

// GetLogger retrieves the request-scoped logger from the Gin context,
// or falls back to slog.Default() if not set.
func GetLogger(ginCtx *gin.Context) *slog.Logger {
	if val, ok := ginCtx.Get(loggerKey); ok {
		if logger, isLogger := val.(*slog.Logger); isLogger {
			return logger
		}
	}
	return slog.Default()
}
