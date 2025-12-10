.PHONY: build test clean install run deps lint

# Build the application
build:
	go build -o bin/megobasterd cmd/megobasterd/main.go

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/megobasterd-linux-amd64 cmd/megobasterd/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/megobasterd-windows-amd64.exe cmd/megobasterd/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/megobasterd-darwin-amd64 cmd/megobasterd/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/megobasterd-darwin-arm64 cmd/megobasterd/main.go

# Run the application
run:
	go run cmd/megobasterd/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run linter
lint:
	golangci-lint run ./...

# Install the application
install:
	go install cmd/megobasterd/main.go

# Package with Fyne (requires fyne CLI)
package-fyne:
	fyne package -os linux -icon assets/icon.png
	fyne package -os windows -icon assets/icon.png
	fyne package -os darwin -icon assets/icon.png

# Development build with race detector
dev:
	go build -race -o bin/megobasterd-dev cmd/megobasterd/main.go
	./bin/megobasterd-dev
