package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func StatusController(w http.ResponseWriter, r *http.Request) {
	// Simulate some dynamic status
	statuses := []string{"Online", "Processing", "Idle", "Active", "Busy"}
	status := statuses[rand.Intn(len(statuses))]
	
	timestamp := time.Now().Format("15:04:05")
	
	// Return HTML fragment for HTMX to swap
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="p-4 bg-green-100 rounded">
			<span class="font-bold">Status:</span> %s
			<span class="text-gray-600 ml-2">(%s)</span>
		</div>
	`, status, timestamp)
}