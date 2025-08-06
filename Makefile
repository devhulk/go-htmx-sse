# run templ generation in watch mode to detect all .templ files and
# re-create _templ.go files on change, then send reload event to browser.
# Default url: http://localhost:7331
live/templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# run tailwindcss to generate the output.css bundle in watch mode.
live/tailwind:
	npx tailwindcss -i ./assets/css/input.css -o ./assets/output.css --watch

# run postcss in watch mode
live/postcss:
	npx postcss ./assets/css/input.css -o ./assets/output.css -w

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	go run github.com/air-verse/air@v1.62.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "echo 'assets changed'" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

# start all 4 watch processes in parallel.
live:
	make -j4 live/templ live/server live/tailwind live/sync_assets

# Clean build artifacts
clean:
	rm -rf tmp/
	rm -f main
	rm -f assets/output.css
	find . -name "*_templ.go" -type f -delete

# ===================== Build ===================== #
# ================== (Production) ================= #

# Generate templ files for production
build/templ:
	templ generate -v

# Build the Go server for production
build/server:
	go build -o main .

# Generate the output.css bundle for production
build/tailwind:
	npx tailwindcss -i ./assets/css/input.css -o ./assets/output.css --minify

# Generate postcss for production
build/postcss:
	npx postcss -i ./assets/css/input.css -o ./assets/output.css

# Run all build processes sequentially for production
build: build/templ build/server build/tailwind

# ===================== Dev Setup =================== #

# Install all dependencies
setup:
	go mod download
	go install github.com/a-h/templ/cmd/templ@latest
	npm install

# Run tests
test:
	go test ./...

# Run tests with coverage
test/coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	go fmt ./...
	templ fmt .

# Lint code
lint:
	golangci-lint run

.PHONY: live live/templ live/server live/tailwind live/postcss live/sync_assets \
	clean build build/templ build/server build/tailwind build/postcss \
	setup test test/coverage fmt lint