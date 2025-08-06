package controllers

import (
	"fmt"
	"net/http"
	"time"
)

func SSEController(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel to signal when the client disconnects
	clientGone := r.Context().Done()

	// Create a ticker to send events every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Send initial message event
	fmt.Fprintf(w, "event: message\n")
	fmt.Fprintf(w, "data: <div class=\"p-4 bg-green-100 rounded\">Connected to SSE stream</div>\n\n")
	w.(http.Flusher).Flush()

	counter := 0
	for {
		select {
		case <-clientGone:
			// Client disconnected
			return
		case <-ticker.C:
			// Send periodic update
			counter++
			timestamp := time.Now().Format("15:04:05")
			
			// Send as HTML that HTMX can swap
			fmt.Fprintf(w, "event: message\n")
			fmt.Fprintf(w, "data: <div class=\"p-4 bg-blue-100 rounded\">Update #%d at %s</div>\n\n", counter, timestamp)
			w.(http.Flusher).Flush()
		}
	}
}