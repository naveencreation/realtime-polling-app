package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestLoggerAttachesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	router := gin.New()
	router.Use(RequestLogger(logger))
	router.GET("/test", func(c *gin.Context) {
		reqLogger := GetLogger(c)
		if reqLogger == nil {
			t.Error("expected logger in context, got nil")
		}
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	requestID := resp.Header().Get(RequestIDHeader)
	if requestID == "" {
		t.Error("expected X-Request-ID response header, got empty string")
	}

	logOutput := logBuffer.String()
	if len(logOutput) == 0 {
		t.Error("expected log output, got empty buffer")
	}
}

func TestRequestLoggerPreservesIncomingRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	router := gin.New()
	router.Use(RequestLogger(logger))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	customID := "client-custom-req-12345"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(RequestIDHeader, customID)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if got := resp.Header().Get(RequestIDHeader); got != customID {
		t.Errorf("expected preserved request ID %q, got %q", customID, got)
	}
}
