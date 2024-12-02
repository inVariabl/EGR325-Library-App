package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))
}

// PageData holds the data to be passed to templates
type PageData struct {
	Title string
	Data  interface{}
}

// renderTemplate executes the template with the given name and data
func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
