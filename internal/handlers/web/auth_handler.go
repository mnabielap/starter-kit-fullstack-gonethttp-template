package web

import (
	"net/http"
	"starter-kit-fullstack-gonethttp-template/pkg/view"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) ViewLogin(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "auth/login", map[string]interface{}{
		"Title": "Login",
	}, "auth")
}

func (h *AuthHandler) ViewRegister(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "auth/register", map[string]interface{}{
		"Title": "Register",
	}, "auth")
}

func (h *AuthHandler) ViewForgotPassword(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "auth/forgot-password", map[string]interface{}{
		"Title": "Forgot Password",
	}, "auth")
}