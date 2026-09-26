package auth

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"polling-backend/internal/middleware"
	"polling-backend/internal/models"
	"polling-backend/pkg/response"
)

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// Handler exposes authentication HTTP endpoints for user registration, session management, and logout.
type Handler struct {
	Repo       Repository
	Service    Service
	AuthCookie string
	Secure     bool
}

type credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler Handler) setAuthCookie(ginCtx *gin.Context, authToken string) {
	sameSite := http.SameSiteLaxMode
	if handler.Secure {
		sameSite = http.SameSiteNoneMode
	}
	maxAgeSeconds := int(handler.Service.Expiry.Seconds())
	http.SetCookie(ginCtx.Writer, &http.Cookie{
		Name:     handler.AuthCookie,
		Value:    authToken,
		MaxAge:   maxAgeSeconds,
		HttpOnly: true,
		Secure:   handler.Secure,
		SameSite: sameSite,
		Path:     "/",
	})
}

// Signup handles POST /api/auth/signup. Validates input credentials, hashes the password,
// creates the user record, and issues an HTTP-only auth cookie and JSON response.
func (handler Handler) Signup(ginCtx *gin.Context) {
	logger := middleware.GetLogger(ginCtx)

	var req credentials
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "invalid request body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "username, email and password are required")
		return
	}
	if !emailRegex.MatchString(req.Email) {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "invalid email format")
		return
	}
	if len([]rune(req.Password)) < 8 {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "password must be at least 8 characters")
		return
	}

	existingUser, err := handler.Repo.FindByEmailOrUsername(ginCtx, req.Email, req.Username)
	if err != nil {
		logger.Error("database error looking up existing user", "email", req.Email, "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}
	if existingUser != nil {
		response.Error(ginCtx, http.StatusConflict, "conflict", "username or email already in use")
		return
	}

	passwordHash, err := handler.Service.HashPassword(req.Password)
	if err != nil {
		logger.Error("bcrypt error hashing password", "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}

	newUser := &models.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}

	if err = handler.Repo.Create(ginCtx, newUser); err != nil {
		if errors.Is(err, ErrDuplicateUser) {
			response.Error(ginCtx, http.StatusConflict, "conflict", "username or email already in use")
		} else {
			logger.Error("database error creating user", "user_id", newUser.ID.Hex(), "error", err)
			response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		}
		return
	}

	authToken, err := handler.Service.GenerateToken(newUser.ID.Hex(), newUser.Username)
	if err != nil {
		logger.Error("jwt error generating auth token", "user_id", newUser.ID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "failed to generate authentication token")
		return
	}

	if authToken != "" {
		handler.setAuthCookie(ginCtx, authToken)
	}

	response.JSON(ginCtx, http.StatusCreated, gin.H{
		"id":       newUser.ID.Hex(),
		"username": newUser.Username,
		"token":    authToken,
	})
}

// Login handles POST /api/auth/login. Authenticates email and password, returning a signed JWT.
func (handler Handler) Login(ginCtx *gin.Context) {
	logger := middleware.GetLogger(ginCtx)

	var req credentials
	if err := ginCtx.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Email) == "" || req.Password == "" {
		response.Error(ginCtx, http.StatusBadRequest, "validation_error", "email and password are required")
		return
	}

	user, err := handler.Repo.FindByEmail(ginCtx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		logger.Error("database error looking up user during login", "email", req.Email, "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}
	if user == nil || !handler.Service.ComparePassword(user.PasswordHash, req.Password) {
		response.Error(ginCtx, http.StatusUnauthorized, "unauthorized", "invalid email or password")
		return
	}

	authToken, err := handler.Service.GenerateToken(user.ID.Hex(), user.Username)
	if err != nil {
		logger.Error("jwt error generating token during login", "user_id", user.ID.Hex(), "error", err)
		response.Error(ginCtx, http.StatusInternalServerError, "server_error", "something went wrong")
		return
	}

	handler.setAuthCookie(ginCtx, authToken)
	response.JSON(ginCtx, http.StatusOK, gin.H{
		"id":       user.ID.Hex(),
		"username": user.Username,
		"token":    authToken,
	})
}

// Logout handles POST /api/auth/logout. Clears the authentication session cookie.
func (handler Handler) Logout(ginCtx *gin.Context) {
	sameSite := http.SameSiteLaxMode
	if handler.Secure {
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(ginCtx.Writer, &http.Cookie{
		Name:     handler.AuthCookie,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		HttpOnly: true,
		Secure:   handler.Secure,
		SameSite: sameSite,
		Path:     "/",
	})
	ginCtx.Status(http.StatusOK)
}
