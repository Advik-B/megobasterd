.PHONY: build test clean install run deps lint dev

# Development - run with hot reload
dev:
	wails dev

# Build the application
build:
	wails build

# Build for production (optimized)
build-prod:
	wails build -clean -upx

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf build/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy
	cd frontend && npm install

# Run linter
lint:
	golangci-lint run ./...

# Install Wails
install-wails:
	go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Generate Wails bindings
generate:
	wails generate module

# Doctor - check Wails setup
doctor:
	wails doctor

