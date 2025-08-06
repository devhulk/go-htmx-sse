package controllers

import (
	"net/http"

	"github.com/devhulk/go-htmx-sse/views"
)

func HomeController(w http.ResponseWriter, r *http.Request) {
	component := views.Home()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}