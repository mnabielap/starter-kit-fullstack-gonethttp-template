package web

import (
	"net/http"
	"starter-kit-fullstack-gonethttp-template/pkg/view"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "users/index", map[string]interface{}{
		"Title":     "User List",
		"PageTitle": "Users",
	}, "main")
}

func (h *UserHandler) CreateView(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "users/create", map[string]interface{}{
		"Title":     "Create User",
		"PageTitle": "Users",
	}, "main")
}

func (h *UserHandler) EditView(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "users/edit", map[string]interface{}{
		"Title":     "Edit User",
		"PageTitle": "Users",
	}, "main")
}