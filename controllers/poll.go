package controllers

import (
	"net/http"

	"github.com/devhulk/go-htmx-sse/views"
)

func PollController(w http.ResponseWriter, r *http.Request) {
	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return just the main content
		component := views.PollContent()
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Return full page with layout
		component := views.PollExample()
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
