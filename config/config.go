package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Name string
		Env  string
		Port string
		URL  string
	}
	DB struct {
		Driver   string // sqlite or postgres
		Name     string
		Host     string
		Port     string
		User     string
		Password string
		SSLMode  string
	}
	JWT struct {
		Secret                   string
		AccessExpirationMinutes  int
		RefreshExpirationDays    int
		ResetPasswordExpiration  int // Minutes
		VerifyEmailExpiration    int // Minutes
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}
}

// LoadConfig loads the environment variables into the Config struct
func LoadConfig() *Config {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{}

	// App
	cfg.App.Name = getEnv("APP_NAME", "StarterKit")
	cfg.App.Env = getEnv("APP_ENV", "development")
	cfg.App.Port = getEnv("PORT", "8080")
	cfg.App.URL = getEnv("APP_URL", "http://localhost:8080")

	// Database
	cfg.DB.Driver = getEnv("DB_DRIVER", "sqlite")
	cfg.DB.Name = getEnv("DB_NAME", "starter_kit_db.sqlite")
	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "root")
	cfg.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	// JWT
	cfg.JWT.Secret = getEnv("JWT_SECRET", "default_secret_please_change")
	cfg.JWT.AccessExpirationMinutes, _ = strconv.Atoi(getEnv("JWT_ACCESS_EXPIRATION_MINUTES", "30"))
	cfg.JWT.RefreshExpirationDays, _ = strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION_DAYS", "30"))
	cfg.JWT.ResetPasswordExpiration, _ = strconv.Atoi(getEnv("JWT_RESET_PASSWORD_EXPIRATION_MINUTES", "10"))
	cfg.JWT.VerifyEmailExpiration, _ = strconv.Atoi(getEnv("JWT_VERIFY_EMAIL_EXPIRATION_MINUTES", "10"))

	// SMTP
	cfg.SMTP.Host = getEnv("SMTP_HOST", "")
	cfg.SMTP.Port, _ = strconv.Atoi(getEnv("SMTP_PORT", "587"))
	cfg.SMTP.Username = getEnv("SMTP_USERNAME", "")
	cfg.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	cfg.SMTP.From = getEnv("EMAIL_FROM", "noreply@example.com")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}