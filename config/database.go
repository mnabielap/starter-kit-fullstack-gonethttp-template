package config

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB initializes the database connection
func ConnectDB(cfg *Config) {
	var err error
	var dsn string
	var dialector gorm.Dialector

	if cfg.DB.Driver == "sqlite" {
		dsn = cfg.DB.Name // For SQLite, DB_NAME is the file path
		dialector = sqlite.Open(dsn)
		log.Printf("Connecting to SQLite: %s", dsn)
	} else {
		// Postgres DSN
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port, cfg.DB.SSLMode)
		dialector = postgres.Open(dsn)
		log.Printf("Connecting to Postgres at %s", cfg.DB.Host)
	}

	// Configure GORM Logger
	logLevel := logger.Error
	if cfg.App.Env == "development" {
		logLevel = logger.Info
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully")
}