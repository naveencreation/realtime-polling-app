package config

import (
	"os"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	// Clear relevant env vars
	os.Unsetenv("MONGO_URI")
	os.Unsetenv("JWT_SECRET")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when MONGO_URI is missing, got nil")
	}

	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	defer os.Unsetenv("MONGO_URI")

	_, err = Load()
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is missing, got nil")
	}

	os.Setenv("JWT_SECRET", "super-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected successful config load, got: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.BcryptCost != 11 {
		t.Errorf("expected default bcrypt cost 11, got %d", cfg.BcryptCost)
	}
}

func TestConfigInvalidBcryptCost(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("JWT_SECRET", "super-secret-key")
	os.Setenv("BCRYPT_COST", "2") // Valid range is 4..31
	defer os.Unsetenv("MONGO_URI")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("BCRYPT_COST")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error with invalid bcrypt cost 2, got nil")
	}
}
