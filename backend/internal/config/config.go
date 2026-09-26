package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all runtime configuration values loaded from environment variables or .env files.
type Config struct {
	Port            string
	MongoURI        string
	MongoDBName     string
	RedisAddr       string
	JWTSecret       string
	FrontendBaseURL string
	JWTExpiryHours  int
	BcryptCost      int
	CookieSecure    bool
	AuthCookieName  string
	VoterCookieName string
}

// Load populates configuration from active environment variables, falling back to local defaults.
// Validates required database URIs and security thresholds before returning.
func Load() (Config, error) {
	loadDotEnv()

	appConfig := Config{
		Port:            getEnv("PORT", "8080"),
		MongoURI:        os.Getenv("MONGO_URI"),
		MongoDBName:     getEnv("MONGO_DB_NAME", "polling_app"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		FrontendBaseURL: getEnv("FRONTEND_BASE_URL", "http://localhost:5173"),
		AuthCookieName:  "auth_token",
		VoterCookieName: "voter_token",
		JWTExpiryHours:  getIntEnv("JWT_EXPIRY_HOURS", 24),
		BcryptCost:      getIntEnv("BCRYPT_COST", 11),
		CookieSecure:    getBoolEnv("COOKIE_SECURE", false),
	}

	if appConfig.MongoURI == "" {
		return appConfig, fmt.Errorf("MONGO_URI is required")
	}
	if appConfig.JWTSecret == "" {
		return appConfig, fmt.Errorf("JWT_SECRET is required")
	}
	if appConfig.JWTExpiryHours <= 0 || appConfig.BcryptCost < 4 || appConfig.BcryptCost > 31 {
		return appConfig, fmt.Errorf("invalid JWT_EXPIRY_HOURS or BCRYPT_COST")
	}

	return appConfig, nil
}

func getEnv(key, defaultValue string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	valStr := getEnv(key, strconv.Itoa(defaultValue))
	parsedVal, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return parsedVal
}

func getBoolEnv(key string, defaultValue bool) bool {
	valStr := getEnv(key, strconv.FormatBool(defaultValue))
	parsedVal, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultValue
	}
	return parsedVal
}

// loadDotEnv reads local .env files if present, without overwriting existing environment variables.
func loadDotEnv() {
	searchPaths := []string{".env", "../.env"}
	for _, envPath := range searchPaths {
		fileBytes, err := os.ReadFile(envPath)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(fileBytes), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			keyValueParts := strings.SplitN(line, "=", 2)
			if len(keyValueParts) == 2 {
				envKey := strings.TrimSpace(keyValueParts[0])
				envVal := strings.Trim(strings.TrimSpace(keyValueParts[1]), `"'`+"\r")
				if _, exists := os.LookupEnv(envKey); !exists {
					_ = os.Setenv(envKey, envVal)
				}
			}
		}
		break
	}
}
