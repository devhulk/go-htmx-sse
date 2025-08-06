package controllers

import (
	"net/http"

	"github.com/devhulk/go-htmx-sse/views"
)

func PollController(w http.ResponseWriter, r *http.Request) {
	component := views.PollExample()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}