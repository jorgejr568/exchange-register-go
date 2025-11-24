.PHONY: help build test test-unit test-integration test-coverage mocks clean install-tools migrate run run-sync docker-build docker-up docker-down

# Default target
help:
	@echo "Available targets:"
	@echo "  make build              - Build the application"
	@echo "  make test               - Run all unit tests"
	@echo "  make test-unit          - Run unit tests only"
	@echo "  make test-integration   - Run integration tests (requires database)"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make mocks              - Generate mocks using mockgen"
	@echo "  make install-tools      - Install required tools (mockgen)"
	@echo "  make migrate            - Run database migrations"
	@echo "  make run                - Run the service"
	@echo "  make run-sync           - Run the service with sync worker"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-up          - Start Docker Compose services"
	@echo "  make docker-down        - Stop Docker Compose services"
	@echo "  make clean              - Clean build artifacts"

# Build the application
build:
	@echo "Building application..."
	go build -o exchange.bin

# Install required tools
install-tools:
	@echo "Installing mockgen..."
	go install go.uber.org/mock/mockgen@latest
	@echo "Tools installed successfully!"

# Generate mocks
mocks:
	@echo "Generating mocks..."
	go generate ./...
	@echo "Mocks generated successfully!"

# Run all unit tests
test: test-unit

# Run unit tests only (skip integration tests)
test-unit:
	@echo "Running unit tests..."
	go test ./... -v -short

# Run integration tests (requires database)
test-integration:
	@echo "Running integration tests..."
	@echo "Make sure database is running and TEST_DATABASE_URL is set"
	go test -tags=integration ./internal/exchange/... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out -short
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run database migrations
migrate:
	@echo "Running migrations..."
	go run main.go migrate

# Run the service
run:
	@echo "Starting service..."
	go run main.go service --port 8080

# Run the service with sync worker
run-sync:
	@echo "Starting service with sync worker..."
	go run main.go service --sync --port 8080

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t exchange-register-go .

docker-up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f exchange.bin
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Setup development environment
setup: install-tools docker-up
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Running migrations..."
	@make migrate
	@echo "Development environment ready!"
