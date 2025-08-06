package controllers

import (
	"net/http"

	"github.com/devhulk/go-htmx-sse/views"
)

func HomeController(w http.ResponseWriter, r *http.Request) {
	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return just the main content
		component := views.HomeContent()
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Return full page with layout
		component := views.Home()
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}