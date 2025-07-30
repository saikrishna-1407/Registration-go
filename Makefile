# Variables
GO := go
BINARY_NAME := app
PORT := 8000

# Default target
all: run-dev

# Run the application in development mode
run-dev:
	@echo "Running the application in development mode..."
	$(GO) run cmd/*.go

# Build the application
build:
	@echo "Building the application..."
	$(GO) build -o $(APP_NAME) cmd/main.go
	@echo "Built binary: $(APP_NAME)"

# Run the binary (after building)
run: build
	@echo "Starting the application on port $(PORT)..."
	./$(BINARY_NAME)

# Clean up the binary
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy

# Run tests
test:
	@echo "Running tests..."
	$(GO) test ./... -v

.PHONY: all run-dev build run clean deps test