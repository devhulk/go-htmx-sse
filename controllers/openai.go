package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/devhulk/go-htmx-sse/views"
	"github.com/sashabaranov/go-openai"
)

// OpenAI API client - using environment variable for API key
var openaiClient *openai.Client

func init() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("Warning: OPENAI_API_KEY not set. OpenAI features will not work.")
	}
	openaiClient = openai.NewClient(apiKey)
}

// OpenAIExampleController handles the main Real Example page
func OpenAIExampleController(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		component := views.OpenAIExampleContent()
		component.Render(r.Context(), w)
	} else {
		component := views.OpenAIExample()
		component.Render(r.Context(), w)
	}
}

// OpenAIPollController handles polling-based OpenAI requests
func OpenAIPollController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prompt := r.FormValue("prompt")
	if prompt == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="text-red-600">Please enter a prompt</div>`))
		return
	}

	// Show loading state immediately and poll every second
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div id="poll-result" hx-get="/openai-poll-status" hx-trigger="load delay:1s" hx-swap="outerHTML">
		<div class="bg-orange-50 border border-orange-200 rounded-lg p-4">
			<h4 class="font-semibold text-orange-800 mb-2">Response (Polling):</h4>
			<div class="animate-pulse">
				<div class="bg-orange-200 h-4 rounded mb-2"></div>
				<div class="bg-orange-200 h-4 rounded mb-2 w-3/4"></div>
				<div class="bg-orange-200 h-4 rounded w-1/2"></div>
			</div>
			<p class="text-sm text-orange-600 mt-2">Initializing request...</p>
		</div>
	</div>`))

	// Start background processing
	go processOpenAIRequest(prompt)
}

// Global variables for polling demo (in production, use proper storage)
var (
	pollResult    string
	pollStatus    = "idle" // idle, processing, completed, error
	pollError     string
	pollStartTime time.Time
	pollProgress  string
)

func processOpenAIRequest(prompt string) {
	pollStatus = "processing"
	pollResult = ""
	pollError = ""
	pollStartTime = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens: 150,
	}

	resp, err := openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		pollStatus = "error"
		pollError = err.Error()
		return
	}

	if len(resp.Choices) > 0 {
		pollResult = resp.Choices[0].Message.Content
		pollStatus = "completed"
	} else {
		pollStatus = "error"
		pollError = "No response generated"
	}
}

// OpenAIPollStatusController handles polling status checks
func OpenAIPollStatusController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	switch pollStatus {
	case "processing":
		// Calculate elapsed time to show different status messages
		elapsed := time.Since(pollStartTime).Seconds()
		var statusMessage string
		
		if elapsed < 2 {
			statusMessage = "🔄 Connecting to OpenAI..."
		} else if elapsed < 5 {
			statusMessage = "⚡ Processing your request..."
		} else if elapsed < 10 {
			statusMessage = "🧠 Generating response..."
		} else if elapsed < 15 {
			statusMessage = "✨ Finalizing response..."
		} else {
			statusMessage = "⏳ Almost ready..."
		}
		
		w.Write([]byte(fmt.Sprintf(`<div id="poll-result" hx-get="/openai-poll-status" hx-trigger="load delay:1s" hx-swap="outerHTML">
			<div class="bg-orange-50 border border-orange-200 rounded-lg p-4">
				<h4 class="font-semibold text-orange-800 mb-2">Response (Polling):</h4>
				<div class="animate-pulse">
					<div class="bg-orange-200 h-4 rounded mb-2"></div>
					<div class="bg-orange-200 h-4 rounded mb-2 w-3/4"></div>
					<div class="bg-orange-200 h-4 rounded w-1/2"></div>
				</div>
				<p class="text-sm text-orange-600 mt-2">%s</p>
				<div class="text-xs text-gray-500 mt-1">Elapsed: %.1fs</div>
			</div>
		</div>`, statusMessage, elapsed)))
	case "completed":
		w.Write([]byte(fmt.Sprintf(`<div id="poll-result" class="bg-green-50 border border-green-200 rounded-lg p-4">
			<h4 class="font-semibold text-green-800 mb-2">Response (Polling):</h4>
			<p class="text-gray-700">%s</p>
			<div class="text-sm text-green-600 mt-2">✓ Response complete</div>
		</div>`, strings.ReplaceAll(pollResult, "\n", "<br>"))))
		pollStatus = "idle" // Reset for next request
	case "error":
		w.Write([]byte(fmt.Sprintf(`<div id="poll-result" class="bg-red-50 border border-red-200 rounded-lg p-4">
			<h4 class="font-semibold text-red-800 mb-2">Error:</h4>
			<p class="text-red-700">%s</p>
		</div>`, pollError)))
		pollStatus = "idle" // Reset for next request
	default:
		w.Write([]byte(`<div id="poll-result"></div>`))
	}
}

// OpenAISSEStartController handles the form submission and returns SSE-enabled HTML
func OpenAISSEStartController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prompt := r.FormValue("prompt")
	if prompt == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="text-red-600">Please enter a prompt</div>`))
		return
	}

	// Return HTML that will connect to SSE stream
	w.Header().Set("Content-Type", "text/html")
	escapedPrompt := strings.ReplaceAll(prompt, `"`, `&quot;`)
	w.Write([]byte(fmt.Sprintf(`<div hx-ext="sse" sse-connect="/openai-sse?prompt=%s">
		<div sse-swap="message"></div>
		<div sse-swap="update"></div>
		<div sse-swap="complete" hx-swap="outerHTML" hx-target="#sse-container"></div>
		<div sse-swap="error"></div>
	</div>`, url.QueryEscape(escapedPrompt))))
}


// Global storage for SSE responses (in production, use proper session storage)
var sseResponses = make(map[string]string)

// OpenAISSEController handles Server-Sent Events for streaming OpenAI responses
func OpenAISSEController(w http.ResponseWriter, r *http.Request) {
	prompt := r.URL.Query().Get("prompt")
	sessionID := r.URL.Query().Get("session")

	if prompt == "" {
		http.Error(w, "Prompt required", http.StatusBadRequest)
		return
	}
	
	if sessionID == "" {
		sessionID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}


	// Send initial message
	fmt.Fprintf(w, "event: message\n")
	fmt.Fprintf(w, "data: <div class=\"bg-blue-50 border border-blue-200 rounded-lg p-4\"><h4 class=\"font-semibold text-blue-800 mb-2\">Response (SSE Streaming):</h4><div class=\"text-sm text-blue-600 mb-2\">Connecting to OpenAI...</div><div id=\"sse-response\"></div></div>\n\n")
	flusher.Flush()

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens: 150,
		Stream:    true,
	}

	stream, err := openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "event: error\n")
		fmt.Fprintf(w, "data: <div class=\"text-red-600\">Error: %s</div>\n\n", err.Error())
		flusher.Flush()
		return
	}
	defer stream.Close()

	var fullResponse strings.Builder
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "event: error\n")
			fmt.Fprintf(w, "data: <div class=\"text-red-600\">Stream error: %s</div>\n\n", err.Error())
			flusher.Flush()
			return
		}

		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content
			if content != "" {
				fullResponse.WriteString(content)

				// Send the accumulated response
				escapedContent := strings.ReplaceAll(fullResponse.String(), "\n", "<br>")
				fmt.Fprintf(w, "event: update\n")
				fmt.Fprintf(w, "data: <p class=\"text-gray-700\">%s</p>\n\n", escapedContent)
				flusher.Flush()
			}
			
			// Check if OpenAI indicates completion
			if response.Choices[0].FinishReason != "" {
				break
			}
		}
	}

	// Send completion event with the final response and reset container for next request
	finalResponse := fullResponse.String()
	escapedFinalContent := strings.ReplaceAll(finalResponse, "\n", "<br>")
	
	fmt.Fprintf(w, "event: complete\n")
	fmt.Fprintf(w, "data: <div id=\"sse-container\" class=\"mt-4 min-h-[100px]\"><div class=\"bg-blue-50 border border-blue-200 rounded-lg p-4\"><h4 class=\"font-semibold text-blue-800 mb-2\">Response (SSE Streaming):</h4><p class=\"text-gray-700\">%s</p><div class=\"text-sm text-green-600 mt-2\">✓ Response complete</div></div></div>\n\n", escapedFinalContent)
	flusher.Flush()
	
	// The connection will naturally close when this function returns
}

// OpenAICleanupController handles cleanup after SSE completion (kept for route compatibility)
func OpenAICleanupController(w http.ResponseWriter, r *http.Request) {
	// This controller is no longer needed but kept to avoid breaking existing routes
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div id="sse-container" class="mt-4 min-h-[100px]"></div>`))
}

