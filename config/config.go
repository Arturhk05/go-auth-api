package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Security SecurityConfig
	CORS     CORSConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port        string
	Environment string
	Host        string
	FrontendURL string
	TLSCertFile string
	TLSKeyFile  string
}

type JWTConfig struct {
	Secret                 string
	RefreshSecret          string
	ExpirationHours        int
	RefreshExpirationHours int
}

type SecurityConfig struct {
	MaxLoginAttempts    int
	LockDurationMinutes int
	BcryptCost          int
}

type CORSConfig struct {
	AllowedOrigins []string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Environment: getEnv("SERVER_ENV", "development"),
			Host:        getEnv("SERVER_HOST", "localhost"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
			TLSCertFile: getEnv("TLS_CERT_FILE", ""),
			TLSKeyFile:  getEnv("TLS_KEY_FILE", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "auth_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", ""),
			RefreshSecret:          getEnv("JWT_REFRESH_SECRET", ""),
			ExpirationHours:        getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshExpirationHours: getEnvAsInt("JWT_REFRESH_EXPIRATION_HOURS", 720),
		},
		Security: SecurityConfig{
			MaxLoginAttempts:    getEnvAsInt("MAX_LOGIN_ATTEMPTS", 5),
			LockDurationMinutes: getEnvAsInt("LOCK_DURATION_MINUTES", 15),
			BcryptCost:          getEnvAsInt("BCRYPT_COST", 12),
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")},
		},
	}

	if cfg.JWT.Secret == "" || len(cfg.JWT.Secret) < 32 {
		fmt.Printf("Invalid JWT_SECRET: %q", cfg.JWT.Secret)
		return nil, fmt.Errorf("JWT_SECRET is invalid (minimum 32 characters required)")
	}

	if cfg.JWT.RefreshSecret == "" || len(cfg.JWT.RefreshSecret) < 32 {
		fmt.Printf("Invalid JWT_REFRESH_SECRET: %q", cfg.JWT.RefreshSecret)
		return nil, fmt.Errorf("JWT_REFRESH_SECRET is invalid (minimum 32 characters required)")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
