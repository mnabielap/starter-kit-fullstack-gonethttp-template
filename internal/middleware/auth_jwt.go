package middleware

import (
	"context"
	"net/http"
	"strings"

	"starter-kit-fullstack-gonethttp-template/config"
	"starter-kit-fullstack-gonethttp-template/pkg/response"
	"starter-kit-fullstack-gonethttp-template/pkg/utils"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	UserKey   contextKey = "user"
)

func AuthJWT(cfg *config.Config, requiredRights []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "Please authenticate")
				return
			}

			// Format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "Invalid token format")
				return
			}

			tokenString := parts[1]
			claims, err := utils.ValidateToken(tokenString, cfg.JWT.Secret)
			if err != nil || claims.Type != "access" {
				response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.Sub)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}