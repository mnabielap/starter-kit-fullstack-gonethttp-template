package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

// In-memory store for CSRF tokens
// Note: In a real production app with multiple replicas, use Redis.
var (
	csrfTokens = make(map[string]string) // Map SessionID -> Token
	csrfMu     sync.Mutex
)

const csrfTokenCtxKey contextKey = "csrf_token"

// GenerateCSRFToken creates a new random token
func GenerateCSRFToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetCSRFToken retrieves the token from the request context.
func GetCSRFToken(r *http.Request) string {
	val, ok := r.Context().Value(csrfTokenCtxKey).(string)
	if ok {
		return val
	}
	return ""
}

// CSRFMiddleware enforces CSRF checks on unsafe methods
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionID string
		var token string

		// 1. Check/Create Session Cookie
		cookie, err := r.Cookie("csrf_session")
		if err != nil {
			// Case A: No Cookie (New Session)
			sessionID = GenerateCSRFToken()
			cookie = &http.Cookie{
				Name:     "csrf_session",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			}
			http.SetCookie(w, cookie)

			// Generate New Token
			token = GenerateCSRFToken()
			csrfMu.Lock()
			csrfTokens[sessionID] = token
			csrfMu.Unlock()
		} else {
			// Case B: Existing Cookie
			sessionID = cookie.Value

			csrfMu.Lock()
			storedToken, exists := csrfTokens[sessionID]
			if !exists {
				// Server Restart Case: Cookie exists, but memory is wiped.
				// Regenerate token for this session ID.
				token = GenerateCSRFToken()
				csrfTokens[sessionID] = token
			} else {
				// Normal Case
				token = storedToken
			}
			csrfMu.Unlock()
		}

		// 2. Inject Token into Context (Crucial for View Rendering)
		ctx := context.WithValue(r.Context(), csrfTokenCtxKey, token)
		r = r.WithContext(ctx)

		// 3. Validate on Unsafe Methods
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
			clientToken := r.Header.Get("X-CSRF-TOKEN")
			if clientToken == "" {
				clientToken = r.FormValue("csrf_token")
			}

			if token == "" || clientToken != token {
				http.Error(w, "Invalid CSRF Token", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}