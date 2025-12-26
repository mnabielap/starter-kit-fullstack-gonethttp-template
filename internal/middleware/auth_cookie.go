package middleware

import (
	"net/http"
)

// AuthCookie is a placeholder middleware for server-side cookie authentication.
func AuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}