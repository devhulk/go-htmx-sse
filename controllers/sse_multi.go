package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func SSEMultiEventController(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel to signal when the client disconnects
	clientGone := r.Context().Done()

	// Create different tickers for different event types
	messageTicker := time.NewTicker(2 * time.Second)
	alertTicker := time.NewTicker(5 * time.Second)
	statusTicker := time.NewTicker(3 * time.Second)

	defer messageTicker.Stop()
	defer alertTicker.Stop()
	defer statusTicker.Stop()

	// Send initial connection event
	fmt.Fprintf(w, "event: connected\n")
	fmt.Fprintf(w, "data: <div class=\"p-2 bg-green-100 rounded\">Connected to multi-event SSE stream</div>\n\n")
	w.(http.Flusher).Flush()

	messageCount := 0
	alertCount := 0
	statusCount := 0

	for {
		select {
		case <-clientGone:
			// Client disconnected
			return

		case <-messageTicker.C:
			// Send a regular message event
			messageCount++
			timestamp := time.Now().Format("15:04:05")
			fmt.Fprintf(w, "event: message\n")
			fmt.Fprintf(w, "data: <div class=\"p-2 bg-blue-100 rounded\">Message #%d at %s</div>\n\n", messageCount, timestamp)
			w.(http.Flusher).Flush()

		case <-alertTicker.C:
			// Send an alert event
			alertCount++
			alerts := []string{"Info", "Warning", "Success", "Update"}
			alertType := alerts[rand.Intn(len(alerts))]

			bgColor := map[string]string{
				"Info":    "bg-blue-200",
				"Warning": "bg-yellow-200",
				"Success": "bg-green-200",
				"Update":  "bg-purple-200",
			}[alertType]

			fmt.Fprintf(w, "event: alert\n")
			fmt.Fprintf(w, "data: <div class=\"p-2 %s rounded font-semibold\">🔔 %s Alert #%d</div>\n\n", bgColor, alertType, alertCount)
			w.(http.Flusher).Flush()

		case <-statusTicker.C:
			// Send a status update event
			statusCount++
			statuses := []string{"Online", "Processing", "Idle", "Active"}
			status := statuses[rand.Intn(len(statuses))]

			statusIcon := map[string]string{
				"Online":     "🟢",
				"Processing": "⚡",
				"Idle":       "💤",
				"Active":     "🔥",
			}[status]

			fmt.Fprintf(w, "event: status\n")
			fmt.Fprintf(w, "data: <span class=\"inline-flex items-center px-2 py-1 bg-gray-100 rounded\">%s %s</span>\n\n", statusIcon, status)
			w.(http.Flusher).Flush()
		}
	}
}
