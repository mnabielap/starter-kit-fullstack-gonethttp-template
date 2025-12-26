package routes

import (
	"net/http"

	"starter-kit-fullstack-gonethttp-template/config"
	apiHandlers "starter-kit-fullstack-gonethttp-template/internal/handlers/api"
	webHandlers "starter-kit-fullstack-gonethttp-template/internal/handlers/web"
	"starter-kit-fullstack-gonethttp-template/internal/middleware"
	"starter-kit-fullstack-gonethttp-template/internal/services"

	// Swagger Docs dependency
	_ "starter-kit-fullstack-gonethttp-template/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers struct {
	APIAuth *apiHandlers.AuthHandler
	APIUser *apiHandlers.UserHandler
	WebAuth *webHandlers.AuthHandler
	WebUser *webHandlers.UserHandler
	WebDash *webHandlers.DashboardHandler
}

func RegisterRoutes(cfg *config.Config, h Handlers, userService services.UserService) http.Handler {
	mux := http.NewServeMux()

	// Middleware Definitions
	logger := middleware.Logger
	security := middleware.SecurityHeaders
	rateLimit := middleware.RateLimit
	csrf := middleware.CSRF

	// Auth Middleware
	authJWT := middleware.AuthJWT(cfg, []string{})
	
	// Role Middleware
	requireAdmin := middleware.RequireAdmin(userService)
	requireAdminOrSelf := middleware.RequireAdminOrSelf(userService)

	// ---------------------------
	// 1. Static Files
	// ---------------------------
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	// ---------------------------
	// 2. Swagger Documentation
	// ---------------------------
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL(cfg.App.URL+"/swagger/doc.json"),
	))

	// ---------------------------
	// 3. Web Routes (HTML)
	// ---------------------------
	mux.HandleFunc("GET /login", h.WebAuth.ViewLogin)
	mux.HandleFunc("GET /register", h.WebAuth.ViewRegister)
	mux.HandleFunc("GET /forgot-password", h.WebAuth.ViewForgotPassword)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		h.WebDash.Index(w, r)
	})

	// Web User Management (View Only - API handles logic)
	mux.HandleFunc("GET /users", h.WebUser.Index)
	mux.HandleFunc("GET /users/create", h.WebUser.CreateView)
	mux.HandleFunc("GET /users/edit", h.WebUser.EditView)

	// ---------------------------
	// 4. API Routes (JSON)
	// ---------------------------

	// Public API
	mux.HandleFunc("POST /v1/auth/register", h.APIAuth.Register)
	mux.HandleFunc("POST /v1/auth/login", h.APIAuth.Login)
	mux.HandleFunc("POST /v1/auth/logout", h.APIAuth.Logout)
	mux.HandleFunc("POST /v1/auth/refresh-tokens", h.APIAuth.RefreshTokens)
	mux.HandleFunc("POST /v1/auth/forgot-password", h.APIAuth.ForgotPassword)
	mux.HandleFunc("POST /v1/auth/reset-password", h.APIAuth.ResetPassword)

	// Protected API (Requires Bearer Token)
	
	// GET /users -> Admin Only (List all users)
	mux.Handle("GET /v1/users", authJWT(requireAdmin(http.HandlerFunc(h.APIUser.GetUsers))))
	
	// POST /users -> Admin Only (Create user manually)
	mux.Handle("POST /v1/users", authJWT(requireAdmin(http.HandlerFunc(h.APIUser.CreateUser))))
	
	// GET /users/{id} -> Admin OR Self
	mux.Handle("GET /v1/users/{id}", authJWT(requireAdminOrSelf(http.HandlerFunc(h.APIUser.GetUser))))
	
	// PATCH /users/{id} -> Admin Only
	mux.Handle("PATCH /v1/users/{id}", authJWT(requireAdmin(http.HandlerFunc(h.APIUser.UpdateUser))))
	
	// DELETE /users/{id} -> Admin Only
	mux.Handle("DELETE /v1/users/{id}", authJWT(requireAdmin(http.HandlerFunc(h.APIUser.DeleteUser))))

	// ---------------------------
	// Global Middleware Chain
	// ---------------------------
	handler := security(mux)
	handler = csrf(handler)
	handler = logger(handler)
	
	if cfg.App.Env == "production" {
		handler = rateLimit(handler)
	}

	return handler
}