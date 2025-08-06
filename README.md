# go-htmx-sse

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

## Troubleshooting

### "api.swap is not a function" error
- You have mismatched HTMX and SSE extension versions
- Download both from the same HTMX release

### SSE not connecting
- Check browser console for errors
- Verify the `/events` endpoint is accessible
- Ensure middleware implements `http.Flusher`

### Styling not updating
- Make sure Tailwind watcher is running (`make live` includes this)
- Check that `assets/output.css` is being generated

### Templates not updating
- Ensure Templ watcher is running
- Manually run `templ generate` if needed
- Check for `*_templ.go` files being generated

## License

MIT
