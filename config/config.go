package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Push     PushConfig
	Security SecurityConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port           string
	Environment    string
	AllowedOrigins []string
}

// JWTConfig holds JWT settings
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// PushConfig holds push notification settings
type PushConfig struct {
	FCMServerKey    string
	APNSKeyID       string
	APNSTeamID      string
	APNSBundleID    string
	APNSKeyPath     string
	APNSProduction  bool
}

// SecurityConfig holds security settings
type SecurityConfig struct {
	EncryptionKey string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Parse JWT expiry
	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		jwtExpiry = 24 * time.Hour
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "mms_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENV", "development"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret-change-me"),
			Expiry: jwtExpiry,
		},
		Push: PushConfig{
			FCMServerKey:   getEnv("FCM_SERVER_KEY", ""),
			APNSKeyID:      getEnv("APNS_KEY_ID", ""),
			APNSTeamID:     getEnv("APNS_TEAM_ID", ""),
			APNSBundleID:   getEnv("APNS_BUNDLE_ID", ""),
			APNSKeyPath:    getEnv("APNS_KEY_PATH", ""),
			APNSProduction: getEnv("APNS_PRODUCTION", "false") == "true",
		},
		Security: SecurityConfig{
			EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		},
	}

	AppConfig = config
	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
