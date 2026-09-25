package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrDuplicateUser = errors.New("duplicate user")

type Service struct {
	Secret []byte
	Expiry time.Duration
	Cost   int
}
type claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s Service) HashPassword(v string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(v), s.Cost)
	return string(b), err
}
func (s Service) ComparePassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
func (s Service) GenerateToken(id, username string) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims{UserID: id, Username: username, RegisteredClaims: jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(s.Expiry))}}).SignedString(s.Secret)
}
func (s Service) VerifyToken(raw string) (string, error) {
	var c claims
	t, err := jwt.ParseWithClaims(raw, &c, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.Secret, nil
	})
	if err != nil || !t.Valid || c.UserID == "" {
		return "", ErrInvalidToken
	}
	return c.UserID, nil
}
