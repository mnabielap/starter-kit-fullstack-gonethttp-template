package middleware

import (
	"net/http"

	"starter-kit-fullstack-gonethttp-template/internal/services"
	"starter-kit-fullstack-gonethttp-template/pkg/response"

	"github.com/google/uuid"
)

// RequireAdmin ensures the authenticated user has the 'admin' role.
func RequireAdmin(service services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get UserID from context (set by AuthJWT)
			userIDStr, ok := r.Context().Value(UserIDKey).(string)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// 2. Fetch User from DB
			id, _ := uuid.Parse(userIDStr)
			user, err := service.GetUserByID(id)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "User not found")
				return
			}

			// 3. Check Role
			if user.Role != "admin" {
				response.Error(w, http.StatusForbidden, "Forbidden: Admins only")
				return
			}

			// 4. Proceed
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdminOrSelf ensures the user is admin OR accessing their own resource.
func RequireAdminOrSelf(service services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get UserID from context
			userIDStr, ok := r.Context().Value(UserIDKey).(string)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// 2. Get Target ID from URL
			targetID := r.PathValue("id")

			// 3. Check if Self
			if targetID == userIDStr {
				next.ServeHTTP(w, r)
				return
			}

			// 4. If not self, Check if Admin
			id, _ := uuid.Parse(userIDStr)
			user, err := service.GetUserByID(id)
			if err != nil || user.Role != "admin" {
				response.Error(w, http.StatusForbidden, "Forbidden: Access denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}