package handlers

import (
	"net/http"

	"github.com/eth0net/geckosite/templates"
)

// About returns the about us page.
func About(w http.ResponseWriter, r *http.Request) {
	pageData := struct{ Title, Path string }{"About Us", r.URL.Path}
	templates.Parse("about").ExecuteTemplate(w, "layout", pageData)
}
