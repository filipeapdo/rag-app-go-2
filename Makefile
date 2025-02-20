# Makefile for the Go project

# PHONY targets ensure that make always runs these commands
.PHONY: all build test fmt lint clean

# 'all' runs formatting, linting, tests, and then builds the binary.
all: fmt lint test build

# Build the binary from the main package.
build:
	@echo "Building the binary..."
	go build -o bin/ingestion-server ./cmd/ingestion-server

# Run all tests with verbose output.
test:
	@echo "Running tests..."
	go test -v ./...

# Format the code using go fmt.
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint the code.
# Uncomment the lint line if you have golangci-lint installed.
lint:
	@echo "Running lint checks..."
	# golangci-lint run

# Clean up build artifacts.
clean:
	@echo "Cleaning up..."
	rm -rf bin
