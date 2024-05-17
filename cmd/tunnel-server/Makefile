# Variables
PACKAGE_NAME := $(shell basename "$(PWD)")
BUILD_DIR=./build

# Default target
all: build

# Build the Go binary
build:
	@echo "Building the binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Run the Go application
run: build
	@echo "Running the application..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Install dependencies (optional)
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Format the code (optional)
fmt:
	@echo "Formatting the code..."
	@go fmt ./...

# Test the application (optional)
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: all build clean run deps fmt test

