package view

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"starter-kit-fullstack-gonethttp-template/config"
	"starter-kit-fullstack-gonethttp-template/internal/middleware"
)

var (
	baseDir = "web/templates"
	cfg     *config.Config
	once    sync.Once
)

// Init sets the configuration for views
func Init(c *config.Config) {
	once.Do(func() {
		cfg = c
	})
}

// Render renders a template with a layout
func Render(w http.ResponseWriter, r *http.Request, viewPath string, data map[string]interface{}, layout string) {
	if data == nil {
		data = make(map[string]interface{})
	}

	// Add Global Data
	data["AppURL"] = cfg.App.URL
	data["AppName"] = cfg.App.Name
	data["CSRFToken"] = middleware.GetCSRFToken(r) // Inject CSRF token

	// Define standard functions for templates
	funcMap := template.FuncMap{
		"baseUrl": func(path string) string {
			if path == "/" {
				return cfg.App.URL
			}
			return fmt.Sprintf("%s/%s", cfg.App.URL, path)
		},
	}

	// Build path to files
	// Layouts are in web/templates/layouts/
	layoutFile := filepath.Join(baseDir, "layouts", layout+".html")
	// Views are relative to web/templates/
	viewFile := filepath.Join(baseDir, viewPath+".html")
	
	// Partials (Always include all partials for simplicity in this starter kit)
	partials, _ := filepath.Glob(filepath.Join(baseDir, "partials", "*.html"))

	files := append([]string{layoutFile, viewFile}, partials...)

	tmpl, err := template.New(filepath.Base(layoutFile)).Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template Execute Error: "+err.Error(), http.StatusInternalServerError)
	}
}