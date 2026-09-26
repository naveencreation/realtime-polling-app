package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"polling-backend/pkg/response"
)

// RequireAuth checks the presence and validity of the authentication JWT in cookie or Authorization header.
func RequireAuth(service Service, cookieName string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		authToken, err := ginCtx.Cookie(cookieName)
		if err != nil || authToken == "" {
			authHeader := ginCtx.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				authToken = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
				err = nil
			}
		}
		if err != nil || authToken == "" {
			response.Error(ginCtx, http.StatusUnauthorized, "unauthorized", "authentication required")
			return
		}
		userID, err := service.VerifyToken(authToken)
		if err != nil {
			response.Error(ginCtx, http.StatusUnauthorized, "unauthorized", "authentication required")
			return
		}
		ginCtx.Set("userId", userID)
		ginCtx.Next()
	}
}
