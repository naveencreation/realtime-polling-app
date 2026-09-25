package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port, MongoURI, MongoDBName, RedisAddr, JWTSecret, FrontendBaseURL string
	JWTExpiryHours, BcryptCost                                         int
	CookieSecure                                                       bool
	AuthCookieName, VoterCookieName                                    string
}

func Load() (Config, error) {
	c := Config{Port: env("PORT", "8080"), MongoURI: os.Getenv("MONGO_URI"), MongoDBName: env("MONGO_DB_NAME", "polling_app"), RedisAddr: env("REDIS_ADDR", "localhost:6379"), JWTSecret: os.Getenv("JWT_SECRET"), FrontendBaseURL: env("FRONTEND_BASE_URL", "http://localhost:5173"), AuthCookieName: "auth_token", VoterCookieName: "voter_token", JWTExpiryHours: intEnv("JWT_EXPIRY_HOURS", 24), BcryptCost: intEnv("BCRYPT_COST", 11), CookieSecure: boolEnv("COOKIE_SECURE", false)}
	if c.MongoURI == "" {
		return c, fmt.Errorf("MONGO_URI is required")
	}
	if c.JWTSecret == "" {
		return c, fmt.Errorf("JWT_SECRET is required")
	}
	if c.JWTExpiryHours <= 0 || c.BcryptCost < 4 || c.BcryptCost > 31 {
		return c, fmt.Errorf("invalid JWT_EXPIRY_HOURS or BCRYPT_COST")
	}
	return c, nil
}
func env(k, d string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return d
}
func intEnv(k string, d int) int {
	v, err := strconv.Atoi(env(k, strconv.Itoa(d)))
	if err != nil {
		return d
	}
	return v
}
func boolEnv(k string, d bool) bool {
	v, err := strconv.ParseBool(env(k, strconv.FormatBool(d)))
	if err != nil {
		return d
	}
	return v
}
