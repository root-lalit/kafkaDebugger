.PHONY: build run clean test docker-up docker-down demo help

# Build the application
build:
	@echo "Building Kafka Debugger..."
	@go build -o kafkaDebugger .
	@echo "✓ Build complete: ./kafkaDebugger"

# Run the application
run: build
	@echo "Starting Kafka Debugger..."
	@./kafkaDebugger

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f kafkaDebugger
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "✓ Dependencies tidied"

# Start Docker Compose environment
docker-up:
	@echo "Starting Kafka environment..."
	@docker-compose up -d
	@echo "✓ Kafka is starting... wait 30 seconds for it to be ready"
	@echo "  You can check status with: docker-compose ps"

# Stop Docker Compose environment
docker-down:
	@echo "Stopping Kafka environment..."
	@docker-compose down
	@echo "✓ Kafka environment stopped"

# Run demo setup (requires running Kafka)
demo:
	@echo "Setting up demo data..."
	@./demo-setup.sh

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@echo "✓ Dependencies installed"

# Show help
help:
	@echo "Kafka Debugger - Makefile commands:"
	@echo ""
	@echo "  make build       - Build the application"
	@echo "  make run         - Build and run the application"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make tidy        - Tidy Go dependencies"
	@echo "  make deps        - Download dependencies"
	@echo "  make docker-up   - Start Docker Compose Kafka environment"
	@echo "  make docker-down - Stop Docker Compose environment"
	@echo "  make demo        - Set up demo data in Kafka"
	@echo "  make help        - Show this help message"
	@echo ""
	@echo "Quick start:"
	@echo "  make docker-up && sleep 30 && make demo && make run"

# Default target
.DEFAULT_GOAL := help
