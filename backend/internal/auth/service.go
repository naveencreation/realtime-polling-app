package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateUser      = errors.New("duplicate user")
)

// Service provides password hashing and JWT token issuance and validation.
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

// HashPassword hashes a plaintext password using bcrypt at the configured cost.
func (service Service) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), service.Cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword checks if a plaintext password matches a stored bcrypt hash.
func (service Service) ComparePassword(passwordHash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}

// GenerateToken creates an HMAC-SHA256 signed JWT containing the user ID and username claims.
func (service Service) GenerateToken(userID, username string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(service.Expiry)),
		},
	})
	return token.SignedString(service.Secret)
}

// VerifyToken validates the JWT signature and expiration, returning the subject's userID.
func (service Service) VerifyToken(tokenString string) (string, error) {
	var tokenClaims claims
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return service.Secret, nil
	})
	if err != nil || !token.Valid || tokenClaims.UserID == "" {
		return "", ErrInvalidToken
	}
	return tokenClaims.UserID, nil
}
