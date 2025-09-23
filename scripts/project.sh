#!/bin/bash

# project.sh - Demo Go Project Management Script
# Usage: ./project.sh [command] [options]

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Project configuration
PROJECT_NAME="demo-go"
DEFAULT_PORT="8080"
DEFAULT_REPOSITORY_TYPE="memory"
DEFAULT_CACHE_TYPE="memory"

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo -e "${PURPLE}🚀 $1${NC}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    local missing_deps=()
    
    if ! command_exists go; then
        missing_deps+=("go")
    fi
    
    if ! command_exists git; then
        missing_deps+=("git")
    fi
    
    if [ ${#missing_deps[@]} -eq 0 ]; then
        print_success "All prerequisites are installed"
        return 0
    else
        print_error "Missing dependencies: ${missing_deps[*]}"
        print_info "Please install the missing dependencies and try again"
        return 1
    fi
}

# Function to install Go dependencies
install_deps() {
    print_header "Installing Go Dependencies"
    
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the project root."
        exit 1
    fi
    
    print_info "Downloading Go modules..."
    go mod download
    
    print_info "Tidying Go modules..."
    go mod tidy
    
    print_success "Dependencies installed successfully"
}

# Function to build the project
build_project() {
    print_header "Building Project"
    
    local target=${1:-"server"}
    local output_dir="bin"
    
    # Create bin directory if it doesn't exist
    mkdir -p "$output_dir"
    
    print_info "Building $target..."
    
    case $target in
        "server")
            go build -o "$output_dir/server" cmd/server/main.go
            ;;
        "all")
            go build -o "$output_dir/server" cmd/server/main.go
            ;;
        *)
            print_error "Unknown build target: $target"
            exit 1
            ;;
    esac
    
    print_success "Build completed successfully"
}

# Function to run tests
run_tests() {
    print_header "Running Tests"
    
    local test_type=${1:-"all"}
    
    case $test_type in
        "unit")
            print_info "Running unit tests..."
            go test ./internal/... -v
            ;;
        "integration")
            print_info "Running integration tests..."
            go test -tags=integration ./tests/integration/... -v
            ;;
        "coverage")
            print_info "Running tests with coverage..."
            go test -coverprofile=coverage.out ./...
            go tool cover -html=coverage.out -o coverage.html
            print_success "Coverage report generated: coverage.html"
            ;;
        "all")
            print_info "Running all tests..."
            go test ./... -v
            ;;
        *)
            print_error "Unknown test type: $test_type"
            exit 1
            ;;
    esac
    
    print_success "Tests completed"
}

# Function to start dependencies with Docker
start_dependencies() {
    print_header "Starting Dependencies"
    
    if ! command_exists docker; then
        print_warning "Docker not found. Skipping dependency startup."
        return 0
    fi
    
    local deps=${1:-"all"}
    
    case $deps in
        "redis")
            print_info "Starting Redis..."
            docker run -d --name redis-dev -p 6379:6379 redis:alpine || true
            ;;
        "mongodb")
            print_info "Starting MongoDB..."
            docker run -d --name mongodb-dev -p 27017:27017 mongo:5.0 || true
            ;;
        "all")
            print_info "Starting Redis and MongoDB..."
            docker run -d --name redis-dev -p 6379:6379 redis:alpine || true
            docker run -d --name mongodb-dev -p 27017:27017 mongo:5.0 || true
            ;;
        *)
            print_error "Unknown dependency: $deps"
            exit 1
            ;;
    esac
    
    print_success "Dependencies started"
}

# Function to stop dependencies
stop_dependencies() {
    print_header "Stopping Dependencies"
    
    if ! command_exists docker; then
        print_warning "Docker not found. Nothing to stop."
        return 0
    fi
    
    print_info "Stopping and removing containers..."
    docker stop redis-dev mongodb-dev 2>/dev/null || true
    docker rm redis-dev mongodb-dev 2>/dev/null || true
    
    print_success "Dependencies stopped"
}

# Function to run the project
run_project() {
    print_header "Running Project"
    
    local mode=${1:-"development"}
    local port=${2:-$DEFAULT_PORT}
    local repo_type=${3:-$DEFAULT_REPOSITORY_TYPE}
    local cache_type=${4:-$DEFAULT_CACHE_TYPE}
    
    # Set environment variables
    export PORT="$port"
    export REPOSITORY_TYPE="$repo_type"
    export CACHE_TYPE="$cache_type"
    
    case $mode in
        "development"|"dev")
            export GIN_MODE="debug"
            export LOG_LEVEL="debug"
            print_info "Running in development mode"
            print_info "Port: $port"
            print_info "Repository: $repo_type"
            print_info "Cache: $cache_type"
            ;;
        "production"|"prod")
            export GIN_MODE="release"
            export LOG_LEVEL="info"
            print_info "Running in production mode"
            ;;
        *)
            print_error "Unknown mode: $mode"
            exit 1
            ;;
    esac
    
    # Check if binary exists
    if [ ! -f "bin/server" ]; then
        print_info "Binary not found. Building project..."
        build_project
    fi
    
    print_success "Starting server on http://localhost:$port"
    print_info "Press Ctrl+C to stop the server"
    
    # Run the server
    ./bin/server
}

# Function to run with hot reload
run_dev() {
    print_header "Running Project with Hot Reload"
    
    if ! command_exists air; then
        print_info "Installing air for hot reload..."
        go install github.com/cosmtrek/air@latest
    fi
    
    # Set development environment
    export GIN_MODE="debug"
    export LOG_LEVEL="debug"
    export REPOSITORY_TYPE=${1:-$DEFAULT_REPOSITORY_TYPE}
    export CACHE_TYPE=${2:-$DEFAULT_CACHE_TYPE}
    
    print_info "Starting development server with hot reload..."
    print_success "Server will restart automatically on file changes"
    
    air
}

# Function to clean build artifacts
clean() {
    print_header "Cleaning Build Artifacts"
    
    print_info "Removing bin directory..."
    rm -rf bin/
    
    print_info "Removing coverage files..."
    rm -f coverage.out coverage.html
    
    print_info "Cleaning Go module cache..."
    go clean -modcache
    
    print_success "Cleanup completed"
}

# Function to format code
format_code() {
    print_header "Formatting Code"
    
    print_info "Running gofmt..."
    go fmt ./...
    
    if command_exists goimports; then
        print_info "Running goimports..."
        goimports -w .
    fi
    
    print_success "Code formatting completed"
}

# Function to lint code
lint_code() {
    print_header "Linting Code"
    
    if ! command_exists golangci-lint; then
        print_info "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi
    
    print_info "Running golangci-lint..."
    golangci-lint run
    
    print_success "Linting completed"
}

# Function to run security scan
security_scan() {
    print_header "Running Security Scan"
    
    if ! command_exists gosec; then
        print_info "Installing gosec..."
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    fi
    
    print_info "Running gosec security scan..."
    gosec ./...
    
    print_success "Security scan completed"
}

# Function to show project status
show_status() {
    print_header "Project Status"
    
    print_info "Go version: $(go version)"
    
    if [ -f "bin/server" ]; then
        print_success "Binary: bin/server exists"
    else
        print_warning "Binary: bin/server not found"
    fi
    
    if command_exists docker; then
        print_info "Docker containers:"
        docker ps --format "table {{.Names}}\t{{.Status}}" --filter "name=redis-dev" --filter "name=mongodb-dev" 2>/dev/null || print_warning "No development containers running"
    fi
    
    if [ -f "go.mod" ]; then
        print_info "Go module: $(grep '^module' go.mod | cut -d' ' -f2)"
    fi
}

# Function to show help
show_help() {
    cat << EOF
${PURPLE}🚀 Demo Go Project Management Script${NC}

${YELLOW}Usage:${NC}
  ./project.sh [command] [options]

${YELLOW}Commands:${NC}
  ${GREEN}setup${NC}                     Setup project dependencies
  ${GREEN}build${NC} [target]           Build project (default: server)
  ${GREEN}run${NC} [mode] [port]        Run project (modes: dev/development, prod/production)
  ${GREEN}dev${NC} [repo] [cache]       Run with hot reload (default: memory, memory)
  ${GREEN}test${NC} [type]              Run tests (types: unit, integration, coverage, all)
  ${GREEN}deps${NC} [action]            Manage dependencies (start, stop, restart)
  
  ${GREEN}clean${NC}                    Clean build artifacts
  ${GREEN}format${NC}                  Format Go code
  ${GREEN}lint${NC}                    Lint Go code
  ${GREEN}security${NC}                Run security scan
  ${GREEN}status${NC}                  Show project status
  
  ${GREEN}help${NC}                    Show this help message

${YELLOW}Examples:${NC}
  ./project.sh setup                    # Setup project
  ./project.sh run dev 8080            # Run in development mode on port 8080
  ./project.sh run prod                # Run in production mode
  ./project.sh dev memory redis        # Run with hot reload, memory repo, redis cache
  ./project.sh test coverage           # Run tests with coverage
  ./project.sh deps start              # Start Redis and MongoDB
  ./project.sh build                   # Build the server

${YELLOW}Environment Variables:${NC}
  PORT                 Server port (default: 8080)
  REPOSITORY_TYPE      Repository type (memory, mongodb)
  CACHE_TYPE          Cache type (memory, redis)
  GIN_MODE            Gin mode (debug, release)
  LOG_LEVEL           Log level (debug, info, warn, error)

EOF
}

# Main script logic
main() {
    case ${1:-help} in
        "setup")
            check_prerequisites
            install_deps
            ;;
        "build")
            build_project "${2:-server}"
            ;;
        "run")
            run_project "${2:-development}" "${3:-$DEFAULT_PORT}" "${4:-$DEFAULT_REPOSITORY_TYPE}" "${5:-$DEFAULT_CACHE_TYPE}"
            ;;
        "dev")
            run_dev "${2:-$DEFAULT_REPOSITORY_TYPE}" "${3:-$DEFAULT_CACHE_TYPE}"
            ;;
        "test")
            run_tests "${2:-all}"
            ;;
        "deps")
            case ${2:-start} in
                "start")
                    start_dependencies "${3:-all}"
                    ;;
                "stop")
                    stop_dependencies
                    ;;
                "restart")
                    stop_dependencies
                    start_dependencies "${3:-all}"
                    ;;
                *)
                    print_error "Unknown deps action: ${2}"
                    exit 1
                    ;;
            esac
            ;;
        "clean")
            clean
            ;;
        "format")
            format_code
            ;;
        "lint")
            lint_code
            ;;
        "security")
            security_scan
            ;;
        "status")
            show_status
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
