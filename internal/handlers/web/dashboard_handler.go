package web

import (
	"net/http"
	"starter-kit-fullstack-gonethttp-template/pkg/view"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "dashboard/index", map[string]interface{}{
		"Title": "Dashboard",
	}, "main")
}