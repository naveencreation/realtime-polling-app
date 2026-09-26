package auth

import (
	"testing"
	"time"
)

func TestPasswordAndTokenRoundTrip(t *testing.T) {
	authService := Service{
		Secret: []byte("test-secret"),
		Expiry: time.Hour,
		Cost:   4,
	}

	passwordHash, err := authService.HashPassword("correct horse battery staple")
	if err != nil || !authService.ComparePassword(passwordHash, "correct horse battery staple") {
		t.Fatal("password round trip failed")
	}
	if authService.ComparePassword(passwordHash, "wrong-password") {
		t.Fatal("wrong password accepted")
	}

	authToken, err := authService.GenerateToken("user-1", "demo")
	if err != nil {
		t.Fatal(err)
	}

	parsedUserID, err := authService.VerifyToken(authToken)
	if err != nil || parsedUserID != "user-1" {
		t.Fatalf("token round trip failed: %q, err: %v", parsedUserID, err)
	}
}

func TestExpiredTokenRejected(t *testing.T) {
	authService := Service{
		Secret: []byte("test-secret"),
		Expiry: -time.Hour,
		Cost:   4,
	}

	authToken, err := authService.GenerateToken("user-1", "demo")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = authService.VerifyToken(authToken); err == nil {
		t.Fatal("expired token accepted")
	}
}
