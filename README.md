# go-htmx-sse (NOT for PRODUCTION just for learning)

A demonstration project showcasing Server-Sent Events (SSE) and polling patterns using Go, HTMX, and Templ.

## Overview

This project demonstrates real-time communication patterns in web applications:
- **Server-Sent Events (SSE)** - Server pushes updates to the client
- **Polling** - Client periodically requests updates from the server
- Comparison of both approaches with practical examples

## Tech Stack

- **Backend**: Go (1.22+)
- **Templating**: [Templ](https://templ.guide/) - Type-safe HTML templates
- **Frontend**: [HTMX](https://htmx.org/) 2.0 - High-level hypermedia library
- **Styling**: Tailwind CSS
- **Hot Reload**: Air for Go, Templ watch mode

## Prerequisites

- Go 1.22 or higher
- Node.js and npm (for Tailwind CSS)
- Templ CLI (`go install github.com/a-h/templ/cmd/templ@latest`)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/devhulk/go-htmx-sse.git
cd go-htmx-sse
```

2. Install Go dependencies:
```bash
go mod download
```

3. Install Node dependencies:
```bash
npm install
```

4. Install Templ CLI if not already installed:
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

5. (Optional) Set up OpenAI API key for the Real Example:
```bash
export OPENAI_API_KEY="your-openai-api-key-here"
```

## Running the Project

### Development Mode (with hot reload)

```bash
make live
```

This starts 4 concurrent processes:
- `templ generate --watch` - Watches and regenerates Templ files
- `air` - Go server with hot reload
- `tailwindcss --watch` - Watches and compiles CSS
- Asset sync for browser reload

The server will be available at http://localhost:8080

### Production Build

```bash
make build
./main
```

## Available Pages

- `/` - Home page with SSE and polling demos side by side
- `/poll` - Advanced polling examples with different intervals
- `/sse-alt` - Alternative SSE implementation using vanilla JavaScript
- `/sse-debug` - Debug page for testing different SSE configurations
- `/sse-multi` - Multiple SSE event types demonstration
- `/real-example` - **OpenAI Integration** demonstrating both polling and SSE approaches with real API calls

## Project Structure

```
.
├── main.go              # Entry point and route definitions
├── controllers/         # HTTP handlers
│   ├── home.go         # Home page controller
│   ├── sse.go          # SSE endpoint
│   ├── poll.go         # Polling examples
│   └── status.go       # Status endpoint for polling
├── views/              # Templ templates
│   ├── layout.templ    # Base layout
│   ├── home.templ      # Home page template
│   └── poll.templ      # Polling examples template
├── middleware/         # HTTP middleware
│   └── logging.go      # Request logging with Flusher support
├── assets/            # Static files
│   ├── css/           # Tailwind input files
│   ├── js/            # HTMX and extensions
│   └── output.css     # Generated Tailwind output
└── Makefile           # Build and development commands
```

## Important: HTMX SSE Gotchas

### Version Compatibility Issue

**Problem**: The SSE extension must be compatible with your HTMX version. Using mismatched versions will cause errors like:
```
Uncaught TypeError: api.swap is not a function
```

**Solution**: Always download the SSE extension from the same HTMX version:

✅ **Correct approach**:
```bash
# Download HTMX 2.0.6
curl -L -o assets/js/htmx.min.js https://unpkg.com/htmx.org@2.0.6/dist/htmx.min.js

# Download SSE extension FROM THE SAME VERSION
curl -L -o assets/js/htmx-ext-sse.js https://unpkg.com/htmx.org@2.0.6/dist/ext/sse.js
```

❌ **Wrong approach**:
```bash
# Don't mix versions or use third-party packages
curl -L -o assets/js/htmx.min.js https://unpkg.com/htmx.org@2.0.0/dist/htmx.min.js
curl -L -o assets/js/htmx-ext-sse.js https://unpkg.com/htmx-ext-sse@2.2.2/sse.js  # Different package!
```

### Middleware Flusher Support

**Problem**: SSE requires the `http.Flusher` interface for streaming responses. Middleware that wraps `http.ResponseWriter` must implement this interface.

**Solution**: Ensure your middleware properly implements `Flush()`:
```go
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

// Must implement Flush for SSE support
func (rw *responseWriter) Flush() {
    if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
        flusher.Flush()
    }
}
```

### SSE Event Format

**Problem**: HTMX SSE extension expects specific event names when using `sse-swap`.

**Solution**: Send events with the correct format:
```go
// Event name must match what you specify in sse-swap="message"
fmt.Fprintf(w, "event: message\n")
fmt.Fprintf(w, "data: <div>Your HTML content here</div>\n\n")
w.(http.Flusher).Flush()
```

### Browser Connection Limits

Browsers limit the number of SSE connections per domain (typically 6). Keep this in mind when designing applications with multiple SSE streams.

### Multiple Event Types with HTMX SSE

**Problem**: You want to handle different types of SSE events (not just "message" events) and route them to different parts of your UI.

**Solution**: Use custom event names and the `sse-swap` attribute:

```go
// Server sends different event types
fmt.Fprintf(w, "event: message\n")
fmt.Fprintf(w, "data: <div>Regular message</div>\n\n")

fmt.Fprintf(w, "event: alert\n")  
fmt.Fprintf(w, "data: <div>Alert notification</div>\n\n")

fmt.Fprintf(w, "event: status\n")
fmt.Fprintf(w, "data: <div>Status update</div>\n\n")
```

```html
<!-- Client listens for specific events -->
<div hx-ext="sse" sse-connect="/events">
    <!-- Only receives "message" events -->
    <div sse-swap="message">...</div>
    
    <!-- Only receives "alert" events -->
    <div sse-swap="alert">...</div>
    
    <!-- Receives both "message" and "alert" events -->
    <div sse-swap="message,alert">...</div>
    
    <!-- Triggers HTMX request when "alert" event arrives -->
    <div hx-trigger="sse:alert" hx-get="/handle-alert">...</div>
</div>
```

### Single Connection for Multiple Listeners

**Problem**: Creating multiple SSE connections (`sse-connect`) on the same page causes errors and hits browser connection limits.

**Solution**: Use a single parent element with `sse-connect` and multiple child elements with `sse-swap`:

❌ **Wrong - Multiple connections:**
```html
<div hx-ext="sse" sse-connect="/events">
    <div sse-swap="message">Messages here</div>
</div>
<div hx-ext="sse" sse-connect="/events">  <!-- Creates another connection! -->
    <div sse-swap="alert">Alerts here</div>
</div>
```

✅ **Correct - Single connection:**
```html
<div hx-ext="sse" sse-connect="/events">  <!-- Single connection -->
    <div sse-swap="message">Messages here</div>
    <div sse-swap="alert">Alerts here</div>
    <div sse-swap="status">Status here</div>
</div>
```

Benefits:
- More efficient - uses only one connection
- Avoids browser connection limits
- All listeners share the same EventSource
- Cleaner architecture

## Development Tips

1. **Check browser console** - SSE connection issues and HTMX events are logged there
2. **Use the debug page** (`/sse-debug`) to test different SSE configurations
3. **Test SSE endpoint directly**: 
   ```bash
   curl -N -H "Accept: text/event-stream" http://localhost:8080/events
   ```
4. **Ensure Templ files are regenerated** after changes:
   ```bash
   templ generate
   ```
5. **Monitor SSE connections** in browser DevTools:
   - Network tab → Filter by "EventStream"
   - See connection status, events received, and any errors
6. **Use custom event types** to organize your real-time updates:
   - `message` for general updates
   - `alert` for notifications
   - `status` for state changes
   - Create your own domain-specific events

## Troubleshooting

### "api.swap is not a function" error
- You have mismatched HTMX and SSE extension versions
- Download both from the same HTMX release

### SSE not connecting
- Check browser console for errors
- Verify the `/events` endpoint is accessible
- Ensure middleware implements `http.Flusher`
- Check if you have multiple `sse-connect` attributes (should only have one per connection)

### Multiple SSE connection errors
- Look for multiple `sse-connect` attributes on the same page
- Consolidate to a single parent element with `sse-connect`
- Child elements should only have `sse-swap` attributes

### Events not being received
- Verify event names match between server and client
- Server: `fmt.Fprintf(w, "event: myevent\n")`
- Client: `<div sse-swap="myevent">`
- Check browser DevTools Network tab for EventStream data

### Styling not updating
- Make sure Tailwind watcher is running (`make live` includes this)
- Check that `assets/output.css` is being generated

### Templates not updating
- Ensure Templ watcher is running
- Manually run `templ generate` if needed
- Check for `*_templ.go` files being generated

## SSE Workflow Deep Dive

### Complete SSE Flow (Real Example with OpenAI)

The OpenAI SSE implementation demonstrates a complete workflow from form submission to streaming completion. Here's the detailed step-by-step process:

#### 1. User Interaction
- User fills form and clicks "Generate with SSE" button
- Form submits with `hx-post="/openai-sse-start"` targeting `#sse-container`

#### 2. SSE Start Controller (`OpenAISSEStartController`)
Returns HTML with HTMX SSE extension setup:
```html
<div hx-ext="sse" sse-connect="/openai-sse?prompt=...">
    <div sse-swap="message"></div>     <!-- Initial status messages -->
    <div sse-swap="update"></div>      <!-- Streaming content updates -->
    <div sse-swap="complete" hx-swap="outerHTML" hx-target="#sse-container"></div>
    <div sse-swap="error"></div>       <!-- Error handling -->
</div>
```

**Key Elements**:
- `hx-ext="sse"` - Loads HTMX SSE extension
- `sse-connect` - Establishes EventSource connection to streaming endpoint
- `sse-swap` - Defines which elements listen for specific event types
- `hx-swap="outerHTML"` - Critical for cleanup (replaces entire container)

#### 3. HTMX SSE Extension Auto-Connection
- When HTML is inserted into DOM, HTMX automatically creates `EventSource`
- Connects to `/openai-sse?prompt=...`
- Sets up event listeners for different event types

#### 4. SSE Controller Stream Processing (`OpenAISSEController`)

**A. Initial Setup**:
```go
// Set SSE headers
w.Header().Set("Content-Type", "text/event-stream")
w.Header().Set("Cache-Control", "no-cache")
w.Header().Set("Connection", "keep-alive")

// Send initial status
fmt.Fprintf(w, "event: message\n")
fmt.Fprintf(w, "data: <div>Connecting to OpenAI...</div>\n\n")
flusher.Flush()
```

**B. OpenAI Streaming Loop**:
```go
stream, err := openaiClient.CreateChatCompletionStream(ctx, req)
var fullResponse strings.Builder

for {
    response, err := stream.Recv()
    if err == io.EOF { break }
    
    content := response.Choices[0].Delta.Content
    if content != "" {
        fullResponse.WriteString(content)
        
        // Send accumulated response so far
        fmt.Fprintf(w, "event: update\n")
        fmt.Fprintf(w, "data: <p>%s</p>\n\n", fullResponse.String())
        flusher.Flush()
    }
    
    if response.Choices[0].FinishReason != "" { break }
}
```

**C. Completion and Cleanup**:
```go
// Send final response with clean container structure
fmt.Fprintf(w, "event: complete\n")
fmt.Fprintf(w, "data: <div id=\"sse-container\" class=\"mt-4 min-h-[100px]\">")
fmt.Fprintf(w, "<div>%s</div>", finalResponse)
fmt.Fprintf(w, "</div>\n\n")
flusher.Flush()
// Function return closes HTTP connection
```

#### 5. HTMX Event Handling
Throughout the stream:
- **"message" events** → Update initial status div
- **"update" events** → Replace content showing progressive response
- **"complete" event** → `outerHTML` swap replaces entire `#sse-container`
- **"error" events** → Display error messages

#### 6. Automatic Connection Cleanup
When "complete" event fires:
1. `hx-swap="outerHTML"` replaces the entire SSE connector div
2. Removing `sse-connect` element from DOM triggers HTMX cleanup
3. EventSource connection automatically closes
4. New container is ready for next request

### Critical Success Factors

1. **Proper Event Structure**: Each SSE event must have `event:` and `data:` lines
2. **Immediate Flushing**: Call `flusher.Flush()` after each event
3. **Container Replacement**: Use `outerHTML` to replace SSE connector for cleanup
4. **Complete Response Inclusion**: Final event includes full response text
5. **Natural Connection Closure**: HTTP function return closes EventSource

### Common Pitfalls Avoided

- **Connection Reuse**: Each request gets fresh connection (no session conflicts)
- **Manual Cleanup**: No complex JavaScript or `sse-close` management needed
- **Response Preservation**: Complete response included in final event
- **Event Type Confusion**: Clear separation between update and complete events

## License

MIT
