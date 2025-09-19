.PHONY: all build run test clean help deps fmt lint dev
.PHONY: docker-build docker-run docker-compose-up docker-compose-down
.PHONY: test-unit test-integration test-e2e coverage
.PHONY: setup-tools install-tools generate-graphql

# Default target
all: help

# Help command
help:
	@echo "🚀 Demo Go Backend - Clean Architecture"
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make run           Run the server locally"
	@echo "  make build         Build the application"
	@echo "  make test          Run all tests"
	@echo "  make dev           Start development environment"
	@echo "  make fmt           Format code"
	@echo "  make lint          Run linter"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  Build Docker image"
	@echo "  make docker-run    Run in Docker container"
	@echo "  make docker-up     Start with Docker Compose"
	@echo "  make docker-down   Stop Docker Compose"
	@echo ""
	@echo "Testing:"
	@echo "  make test-unit     Run unit tests"
	@echo "  make test-integration  Run integration tests"
	@echo "  make test-e2e      Run end-to-end tests"
	@echo "  make coverage      Generate test coverage report"
	@echo ""
	@echo "Tools:"
	@echo "  make install-tools Install development tools"
	@echo "  make generate      Generate GraphQL code"
	@echo "  make deps          Download dependencies"
	@echo "  make clean         Clean build artifacts"

# Build the application
build:
	@echo "🔨 Building the application..."
	go build -o bin/server cmd/server/main.go

# Run the server
run:
	@echo "🚀 Starting the server..."
	go run cmd/server/main.go

# Development mode with hot reload
dev:
	@echo "🔧 Starting development environment..."
	docker-compose up -d redis mongodb
	@echo "💡 Dependencies started. Run 'make run' to start the server"

# Install dependencies
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Run all tests
test:
	@echo "🧪 Running all tests..."
	go test ./...

# Run unit tests
test-unit:
	@echo "🧪 Running unit tests..."
	go test -short ./...

# Run integration tests
test-integration:
	@echo "🧪 Running integration tests..."
	go test -tags=integration ./tests/integration/...

# Run end-to-end tests
test-e2e:
	@echo "🧪 Running end-to-end tests..."
	go test -tags=e2e ./tests/e2e/...

# Generate test coverage
coverage:
	@echo "📊 Generating test coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install development tools
install-tools:
	@echo "🛠️ Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/99designs/gqlgen@latest

# Generate GraphQL code
generate-graphql:
	@echo "🔄 Generating GraphQL code..."
	go generate ./...

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t demo-go:latest .

docker-run:
	@echo "🐳 Running Docker container..."
	docker run --rm -p 8080:8080 --env-file .env demo-go:latest

docker-up:
	@echo "🐳 Starting Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "🐳 Stopping Docker Compose..."
	docker-compose down

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Security scan
security-scan:
	@echo "🔒 Running security scan..."
	gosec ./...

# Performance test with k6
load-test:
	@echo "⚡ Running load tests..."
	k6 run tests/load-test.js

# Database commands
db-reset:
	@echo "🗄️ Resetting database..."
	docker-compose exec mongodb mongo --eval "db.dropDatabase()" demo_go

# Generate all code
generate: generate-graphql

# Verify everything works
verify: deps fmt lint test
	@echo "✅ All verifications passed!"
