package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	RTSP     RTSPConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type RTSPConfig struct {
	Host          string
	Port          string
	APIURL        string
	PublicBaseURL string
	Username      string
	Password      string
}

type CORSConfig struct {
	AllowedOrigins string
}

// Load membaca konfigurasi dari environment variables
func Load() (*Config, error) {
	// Load .env file jika ada
	godotenv.Load()

	// Parse JWT expiration
	jwtExp, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		jwtExp = 24 * time.Hour
	}

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "CCTV Monitoring API"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "cctv_monitoring"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: jwtExp,
		},
		RTSP: RTSPConfig{
			Host:          getEnv("RTSP_TO_WEB_HOST", "localhost"),
			Port:          getEnv("RTSP_TO_WEB_PORT", "8083"),
			APIURL:        getEnv("RTSP_TO_WEB_API_URL", "http://localhost:8083"),
			PublicBaseURL: getEnv("RTSP_TO_WEB_PUBLIC_URL", "http://localhost:8083"), // ← Pastikan ini
			Username:      getEnv("RTSP_TO_WEB_USERNAME", ""),
			Password:      getEnv("RTSP_TO_WEB_PASSWORD", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},
	}

	return config, nil
}

// GetDSN mengembalikan connection string untuk PostgreSQL
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// getEnv membaca environment variable dengan fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
