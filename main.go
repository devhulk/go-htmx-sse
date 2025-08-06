package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devhulk/go-htmx-sse/controllers"
	"github.com/devhulk/go-htmx-sse/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Static files
	fileServer := http.FileServer(http.Dir("./assets"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Routes
	mux.Handle("/", middleware.LoggingMiddleware(http.HandlerFunc(controllers.HomeController)))
	
	// SSE endpoint
	mux.Handle("/events", middleware.LoggingMiddleware(http.HandlerFunc(controllers.SSEController)))
	
	// Polling endpoint
	mux.Handle("/poll", middleware.LoggingMiddleware(http.HandlerFunc(controllers.PollController)))
	
	// Status endpoint for polling example
	mux.Handle("/status", middleware.LoggingMiddleware(http.HandlerFunc(controllers.StatusController)))
	
	// Alternative SSE demo
	mux.Handle("/sse-alt", middleware.LoggingMiddleware(http.HandlerFunc(controllers.SSEAlternativeController)))
	
	// SSE Debug page
	mux.Handle("/sse-debug", middleware.LoggingMiddleware(http.HandlerFunc(controllers.SSEDebugController)))
	
	// Multi-event SSE demo
	mux.Handle("/multi-events", middleware.LoggingMiddleware(http.HandlerFunc(controllers.SSEMultiEventController)))
	mux.Handle("/sse-multi", middleware.LoggingMiddleware(http.HandlerFunc(controllers.SSEMultiEventPageController)))

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}