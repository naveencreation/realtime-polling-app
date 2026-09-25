package auth

import (
	"testing"
	"time"
)

func TestPasswordAndTokenRoundTrip(t *testing.T) {
	s := Service{Secret: []byte("test-secret"), Expiry: time.Hour, Cost: 4}
	hash, err := s.HashPassword("correct horse battery staple")
	if err != nil || !s.ComparePassword(hash, "correct horse battery staple") {
		t.Fatal("password round trip failed")
	}
	if s.ComparePassword(hash, "wrong") {
		t.Fatal("wrong password accepted")
	}
	token, err := s.GenerateToken("user-1", "demo")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := s.VerifyToken(token); err != nil || got != "user-1" {
		t.Fatalf("token round trip failed: %q %v", got, err)
	}
}

func TestExpiredTokenRejected(t *testing.T) {
	s := Service{Secret: []byte("test-secret"), Expiry: -time.Hour, Cost: 4}
	token, err := s.GenerateToken("user-1", "demo")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.VerifyToken(token); err == nil {
		t.Fatal("expired token accepted")
	}
}
