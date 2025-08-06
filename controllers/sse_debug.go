package controllers

import (
	"net/http"

	"github.com/devhulk/go-htmx-sse/views"
)

func SSEDebugController(w http.ResponseWriter, r *http.Request) {
	component := views.SSEDebug()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}