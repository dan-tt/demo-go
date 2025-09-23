.PHONY: all build run test clean help deps fmt lint dev
.PHONY: docker-build docker-run docker-compose-up docker-compose-down
.PHONY: test-unit test-integration test-e2e coverage
.PHONY: setup-tools install-tools generate-graphql setup-env
.PHONY: project-setup project-build project-run project-dev project-test
.PHONY: project-clean project-format project-lint project-security
.PHONY: setup-graphql setup-cicd

# Default target
all: help

# Help command
help:
	@echo "🚀 Demo Go Backend - Clean Architecture"
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "Environment Setup:"
	@echo "  make setup-env     Setup environment variables"
	@echo "  make env-dev       Copy development environment"
	@echo "  make env-example   Copy example environment"
	@echo ""
	@echo "Project Scripts:"
	@echo "  make project-setup     Setup project dependencies"
	@echo "  make project-run       Run project with script"
	@echo "  make project-dev       Run with hot reload"
	@echo "  make project-test      Run tests via script"
	@echo "  make project-clean     Clean build artifacts"
	@echo "  make project-format    Format code via script"
	@echo "  make project-lint      Lint code via script"
	@echo "  make project-security  Run security scan"
	@echo ""
	@echo "Setup Scripts:"
	@echo "  make setup-graphql     Setup GraphQL with gqlgen"
	@echo "  make setup-cicd        Setup CI/CD pipeline"
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

# Environment setup commands
setup-env: env-check
	@echo "🔧 Environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Edit .env file with your configuration"
	@echo "2. Run 'make dev' to start development environment"
	@echo "3. Visit http://localhost:8080/health to verify"

env-check:
	@if [ ! -f .env ]; then \
		echo "📋 Creating .env file from .env.local..."; \
		cp .env.local .env; \
		echo "✅ .env file created with development defaults"; \
	else \
		echo "ℹ️ .env file already exists"; \
	fi

env-dev:
	@echo "📋 Copying development environment..."
	@cp .env.local .env
	@echo "✅ Development environment configured"

env-example:
	@echo "📋 Copying example environment..."
	@cp .env.example .env
	@echo "⚠️ Please edit .env file with your configuration"

env-validate:
	@echo "🔍 Validating environment configuration..."
	@docker-compose config > /dev/null && echo "✅ Docker Compose configuration is valid" || echo "❌ Docker Compose configuration has errors"

env-show:
	@echo "📋 Current environment configuration:"
	@docker-compose config | grep -A 50 environment: || echo "No environment variables configured"

# =============================================================================
# Project Script Commands
# =============================================================================

# Project setup via script
project-setup:
	@echo "🔧 Setting up project dependencies..."
	@scripts/project.sh setup

# Project build via script
project-build:
	@echo "🏗️ Building project via script..."
	@scripts/project.sh build

# Project run via script
project-run:
	@echo "🚀 Running project via script..."
	@scripts/project.sh run $(MODE) $(PORT)

# Project dev mode via script
project-dev:
	@echo "🛠️ Starting development environment..."
	@scripts/project.sh dev $(REPO) $(CACHE)

# Project test via script
project-test:
	@echo "🧪 Running tests via script..."
	@scripts/project.sh test $(TYPE)

# Project clean via script
project-clean:
	@echo "🧹 Cleaning project via script..."
	@scripts/project.sh clean

# Project format via script
project-format:
	@echo "📝 Formatting code via script..."
	@scripts/project.sh format

# Project lint via script
project-lint:
	@echo "🔍 Linting code via script..."
	@scripts/project.sh lint

# Project security scan via script
project-security:
	@echo "🔒 Running security scan via script..."
	@scripts/project.sh security

# =============================================================================
# Setup Script Commands
# =============================================================================

# Setup GraphQL
setup-graphql:
	@echo "📊 Setting up GraphQL with gqlgen..."
	@scripts/setup-graphql.sh

# Setup CI/CD
setup-cicd:
	@echo "🚀 Setting up CI/CD pipeline..."
	@scripts/setup-cicd.sh

# Verify everything works
verify: deps fmt lint test
	@echo "✅ All verifications passed!"
