package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"polling-backend/internal/models"
	"polling-backend/pkg/response"
	"regexp"
	"strings"
	"time"
)

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

func (h Handler) Signup(c *gin.Context) {
	var in credentials
	if c.ShouldBindJSON(&in) != nil {
		response.Error(c, 400, "validation_error", "invalid request body")
		return
	}
	in.Username, in.Email = strings.TrimSpace(in.Username), strings.TrimSpace(strings.ToLower(in.Email))
	if in.Username == "" || in.Email == "" || in.Password == "" {
		response.Error(c, 400, "validation_error", "username, email and password are required")
		return
	}
	if !regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`).MatchString(in.Email) {
		response.Error(c, 400, "validation_error", "invalid email format")
		return
	}
	if len([]rune(in.Password)) < 8 {
		response.Error(c, 400, "validation_error", "password must be at least 8 characters")
		return
	}
	u, err := h.Repo.FindByEmailOrUsername(c, in.Email, in.Username)
	if err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	if u != nil {
		response.Error(c, 409, "conflict", "username or email already in use")
		return
	}
	hash, err := h.Service.HashPassword(in.Password)
	if err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	u = &models.User{ID: primitive.NewObjectID(), Username: in.Username, Email: in.Email, PasswordHash: hash, CreatedAt: time.Now().UTC()}
	if err = h.Repo.Create(c, u); err != nil {
		if err == ErrDuplicateUser {
			response.Error(c, 409, "conflict", "username or email already in use")
		} else {
			response.Error(c, 500, "server_error", "something went wrong")
		}
		return
	}
	response.JSON(c, 201, gin.H{"id": u.ID.Hex(), "username": u.Username})
}
func (h Handler) Login(c *gin.Context) {
	var in credentials
	if c.ShouldBindJSON(&in) != nil || strings.TrimSpace(in.Email) == "" || in.Password == "" {
		response.Error(c, 400, "validation_error", "email and password are required")
		return
	}
	u, err := h.Repo.FindByEmail(c, strings.ToLower(strings.TrimSpace(in.Email)))
	if err != nil || u == nil || !h.Service.ComparePassword(u.PasswordHash, in.Password) {
		response.Error(c, 401, "unauthorized", "invalid email or password")
		return
	}
	token, err := h.Service.GenerateToken(u.ID.Hex(), u.Username)
	if err != nil {
		response.Error(c, 500, "server_error", "something went wrong")
		return
	}
	maxAge := int(h.Service.Expiry.Seconds())
	http.SetCookie(c.Writer, &http.Cookie{Name: h.AuthCookie, Value: token, MaxAge: maxAge, HttpOnly: true, Secure: h.Secure, SameSite: http.SameSiteLaxMode, Path: "/"})
	response.JSON(c, 200, gin.H{"id": u.ID.Hex(), "username": u.Username})
}
func (h Handler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{Name: h.AuthCookie, MaxAge: -1, Expires: time.Unix(1, 0), HttpOnly: true, Secure: h.Secure, SameSite: http.SameSiteLaxMode, Path: "/"})
	c.Status(200)
}
func NewVoterToken() string { return uuid.NewString() }
