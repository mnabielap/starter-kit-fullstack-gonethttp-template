package main

import (
	"log"
	"net/http"
	"time"

	"starter-kit-fullstack-gonethttp-template/config"
	apiHandlers "starter-kit-fullstack-gonethttp-template/internal/handlers/api"
	webHandlers "starter-kit-fullstack-gonethttp-template/internal/handlers/web"
	"starter-kit-fullstack-gonethttp-template/internal/models"
	"starter-kit-fullstack-gonethttp-template/internal/repository"
	"starter-kit-fullstack-gonethttp-template/internal/routes"
	"starter-kit-fullstack-gonethttp-template/internal/services"
	"starter-kit-fullstack-gonethttp-template/pkg/view"
)

// @title Starter Kit Fullstack Go Native
// @version 1.0
// @description A Fullstack starter kit using Go (net/http), GORM, and JWT. Matches PHP Native features.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Template Engine
	view.Init(cfg)

	// 3. Connect Database
	config.ConnectDB(cfg)

	// 4. Auto Migration
	log.Println("Running Database Migrations...")
	err := config.DB.AutoMigrate(&models.User{}, &models.Token{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 5. Setup Dependency Injection
	userRepo := repository.NewUserRepository(config.DB)
	tokenRepo := repository.NewTokenRepository(config.DB)

	tokenService := services.NewTokenService(tokenRepo, cfg)
	emailService := services.NewEmailService(cfg)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenRepo, tokenService, emailService, cfg)

	handlers := routes.Handlers{
		APIAuth: apiHandlers.NewAuthHandler(authService),
		APIUser: apiHandlers.NewUserHandler(userService),
		WebAuth: webHandlers.NewAuthHandler(),
		WebUser: webHandlers.NewUserHandler(),
		WebDash: webHandlers.NewDashboardHandler(),
	}

	// 6. Setup Router
	// Pass userService here for Middleware Roles
	router := routes.RegisterRoutes(cfg, handlers, userService)

	// 7. Start Server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", cfg.App.Port)
	log.Printf("Swagger Docs available at %s/swagger/index.html", cfg.App.URL)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}