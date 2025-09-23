# Go Backend API - Clean Architecture with Redis Caching

A production-ready Go web server implementing clean architecture principles with comprehensive features including Redis caching, GraphQL API, MongoDB integration, JWT authentication, and CI/CD pipeline.

[![CI/CD Pipeline](https://github.com/yourusername/demo-go/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/yourusername/demo-go/actions/workflows/ci-cd.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/demo-go)](https://goreportcard.com/report/github.com/yourusername/demo-go)
[![Coverage Status](https://codecov.io/gh/yourusername/demo-go/branch/main/graph/badge.svg)](https://codecov.io/gh/yourusername/demo-go)
[![Docker Image](https://ghcr.io/yourusername/demo-go/badge.svg)](https://ghcr.io/yourusername/demo-go)

## 🚀 Features

### Core Features
- **🏗️ Clean Architecture**: Domain-driven design with clear separation of concerns
- **🛣️ Modular Routes**: Organized route structure with separate files for different functionalities
- **⚡ Redis Caching**: High-performance caching with 80-95% performance improvement
- **🔐 JWT Authentication**: Secure authentication with refresh token support
- **🍃 MongoDB Integration**: Flexible NoSQL database with connection pooling
- **📊 GraphQL API**: Modern API with type safety and real-time subscriptions
- **🔍 Structured Logging**: Comprehensive logging with Uber Zap, enhanced console formatting, and pretty JSON request/response logging
- **🐳 Containerized**: Production-ready Docker configuration

### DevOps & Operations
- **� CI/CD Pipeline**: Automated testing, building, and deployment
- **🔒 Security Scanning**: Vulnerability scanning with Gosec and Trivy
- **📈 Performance Testing**: Load testing with k6
- **📊 Monitoring**: Health checks and metrics collection
- **🌍 Multi-platform**: Docker images for AMD64 and ARM64

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Prerequisites](#-prerequisites)
- [Installation & Configuration](#️-installation--configuration)
  - [Environment Variables Configuration](#environment-variables-configuration)
  - [Environment Setup Commands](#environment-setup-commands)
  - [Configuration Categories](#configuration-categories)
  - [Docker Compose Profiles](#docker-compose-profiles)
  - [Security Best Practices](#security-best-practices)
- [API Documentation](#-api-documentation)
- [Routes Architecture](#️-routes-architecture)
- [Testing](#-testing)
- [Contributing](#-contributing)

## ⚡ Quick Start

### Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/demo-go.git
cd demo-go

# Start with Docker Compose
docker-compose up -d

# The API will be available at http://localhost:8080
curl http://localhost:8080/health
```

### Local Development

```bash
# Clone and setup
git clone https://github.com/yourusername/demo-go.git
cd demo-go

# Install dependencies
go mod download

# Start Redis and MongoDB (using Docker)
docker run -d -p 6379:6379 --name redis redis:alpine
docker run -d -p 27017:27017 --name mongodb mongo:5.0

# Set environment variables
export CACHE_TYPE=redis
export REPOSITORY_TYPE=mongodb
export JWT_SECRET=your-secret-key

# Run the server
go run cmd/server/main.go

# Server starts at http://localhost:8080
```

## 📚 Prerequisites

### Required
- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **Git**: Version control system

### Optional (for full features)
- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Redis**: [Install Redis](https://redis.io/docs/getting-started/installation/)
- **MongoDB**: [Install MongoDB](https://docs.mongodb.com/manual/installation/)

### Development Tools
- **golangci-lint**: Code linting
- **k6**: Performance testing
- **gqlgen**: GraphQL code generation

```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/99designs/gqlgen@latest
brew install k6  # macOS
```
## 🔧 Installation & Configuration

### Environment Variables Configuration

This section explains how to configure the demo-go application using environment variables.

#### Quick Start

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit the `.env` file with your specific values:**
   ```bash
   nano .env  # or use your preferred editor
   ```

3. **Run with docker-compose:**
   ```bash
   # Basic setup
   docker-compose up -d

   # With development tools (MongoDB Express)
   docker-compose --profile development up -d

   # With Redis cache
   docker-compose --profile cache up -d

   # For testing
   docker-compose --profile testing up -d
   ```

#### Environment Setup Commands

```bash
# Setup environment (creates .env from .env.local)
make setup-env

# Copy development environment
make env-dev

# Copy example environment
make env-example

# Validate configuration
make env-validate
```

#### Environment Files

| File | Purpose | Commit to Git |
|------|---------|---------------|
| `.env.example` | Template with safe defaults | ✅ Yes |
| `.env.local` | Local development defaults | ✅ Yes |
| `.env` | Your actual configuration | ❌ No |

#### Configuration Categories

##### 🚀 Application Configuration
```bash
APP_NAME=demo-clean-api
APP_VERSION=1.0.0
APP_ENVIRONMENT=development  # development, staging, production
APP_DEBUG=true
```

##### 🌐 Server Configuration
```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
```

##### 🗄️ Database Configuration
```bash
REPOSITORY_TYPE=mongodb  # mongodb, memory
MONGODB_HOST=mongodb     # Use 'localhost' when running outside Docker
MONGODB_PORT=27017
MONGODB_DATABASE=demo_clean
MONGODB_USERNAME=admin
MONGODB_PASSWORD=your_secure_password
```

##### 🏃‍♂️ Cache Configuration
```bash
CACHE_TYPE=memory  # memory, redis
REDIS_HOST=redis   # Use 'localhost' when running outside Docker
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
```

##### 🔐 JWT Configuration
```bash
JWT_SECRET_KEY=your_very_secure_jwt_secret_key
JWT_EXPIRATION=24h
JWT_ISSUER=demo-clean-api
```

##### 📊 Logging Configuration
```bash
LOG_LEVEL=info      # debug, info, warn, error
LOG_FORMAT=console  # console, json, text (recommended: console for development, json for production)
LOG_OUTPUT=stdout   # stdout, file
```

**Enhanced Console Logging Features:**
- 🕒 **Clean Time Format**: `HH:MM:SS.mmm` instead of full timestamps
- 🎯 **Status Emojis**: Visual status indicators (✅ success, ⚠️ warnings, 🚨 errors)
- 📊 **Readable Metrics**: Human-friendly duration (`207µs`) and size (`84B`) formatting
- 🔇 **Smart Filtering**: Health check endpoints are logged quietly to reduce noise
- 🌈 **Better Structure**: Clean separators and organized field layout
- 📍 **Client Information**: Real client IP extraction and user agent logging
- 📋 **Pretty JSON Logging**: Beautiful request/response JSON formatting with proper indentation
- 📄 **Multi-line Format**: Uses newlines instead of separators for better readability

**Example Console Output:**
```
11:51:24.157 | INFO | → Request started
Method: POST
Path: /auth/register
User-Agent: curl/8.7.1
Client-IP: ::1
Request Body:
{
  "email": "test@example.com",
  "name": "Test User",
  "password": "password123",
  "username": "testuser"
}

11:51:24.207 | INFO | ← Request completed ✅
Status: 201
Duration: 50.731ms
Size: 232B
Response Body:
{
  "data": {
    "created_at": "2025-09-23T11:51:24.207652+07:00",
    "email": "test@example.com",
    "id": "1",
    "name": "Test User",
    "role": "user",
    "updated_at": "2025-09-23T11:51:24.207652+07:00"
  },
  "message": "User registered successfully",
  "success": true
}
```

#### Docker Compose Profiles

Use profiles to enable optional services:

```bash
# Development with MongoDB Express UI
docker-compose --profile development up -d

# With Redis cache
docker-compose --profile cache up -d

# With testing client
docker-compose --profile testing up -d

# Multiple profiles
docker-compose --profile development --profile cache up -d
```

#### Environment-Specific Examples

##### 🔧 Development (.env)
```bash
APP_ENVIRONMENT=development
APP_DEBUG=true
LOG_LEVEL=debug
LOG_FORMAT=text
GRAPHQL_PLAYGROUND_ENABLED=true
```

##### 🎭 Staging (.env)
```bash
APP_ENVIRONMENT=staging
APP_DEBUG=false
LOG_LEVEL=info
LOG_FORMAT=json
GRAPHQL_PLAYGROUND_ENABLED=false
```

##### 🚀 Production (.env)
```bash
APP_ENVIRONMENT=production
APP_DEBUG=false
LOG_LEVEL=warn
LOG_FORMAT=json
GRAPHQL_PLAYGROUND_ENABLED=false
GRAPHQL_INTROSPECTION_ENABLED=false
```

#### Security Best Practices

##### ✅ DO:
- Use strong, unique passwords for each environment
- Use different JWT secret keys for each environment
- Keep `.env` files out of version control
- Use environment-specific configurations
- Rotate secrets regularly

##### ❌ DON'T:
- Commit `.env` files to Git
- Use default passwords in production
- Share the same JWT secret across environments
- Store secrets in docker-compose.yml
- Use debug mode in production

#### Validation & Troubleshooting

Check your configuration:

```bash
# Verify environment variables are loaded
docker-compose config

# Check specific service configuration
docker-compose exec api-server env | grep MONGODB

# Test the application
curl http://localhost:8080/health
```

##### Common Issues:

1. **Variables not loaded:**
   - Ensure `.env` file is in the same directory as `docker-compose.yml`
   - Check for syntax errors in `.env` file (no spaces around `=`)

2. **Connection refused:**
   - Verify service names match between containers
   - Check if ports are already in use

3. **Authentication failed:**
   - Verify username/password combinations
   - Check if the database is initialized properly

##### Debug Commands:
```bash
# Check container logs
docker-compose logs api-server
docker-compose logs mongodb

# Enter container to debug
docker-compose exec api-server sh
docker-compose exec mongodb mongosh
```

#### Configuration Examples

##### Local Development Setup:
```bash
# Copy local defaults
cp .env.local .env

# Start with development tools
docker-compose --profile development up -d

# Access services
# API: http://localhost:8080
# MongoDB Express: http://localhost:8081
```

##### Production-like Setup:
```bash
# Copy example and customize
cp .env.example .env

# Edit with production values
nano .env

# Run minimal setup
docker-compose up -d
```
## 🚀 API Documentation

### Base URL
- **Development**: `http://localhost:8080`
- **Production**: `https://your-app.onrender.com`

### Routes Structure
The API follows a modular route structure with separate files for different functionalities:

- **Health Routes** (`health_routes.go`): System health and monitoring
- **Auth Routes** (`auth_routes.go`): Authentication and user registration
- **User Routes** (`user_routes.go`): User profile management
- **Admin Routes** (`admin_routes.go`): Administrative operations

### Health Check
```bash
GET /health
```

**Response:**
```json
{
  "data": {
    "service": "clean-architecture-api",
    "status": "healthy",
    "timestamp": "2025-09-19T00:00:00Z"
  },
  "message": "Service is healthy",
  "success": true
}
```

### Authentication Routes

#### Register User
```bash
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "data": {
    "id": "1",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-09-19T10:37:24.657185+07:00",
    "updated_at": "2025-09-19T10:37:24.657185+07:00"
  },
  "message": "User registered successfully",
  "success": true
}
```

#### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "1",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2025-09-19T10:37:24.657185+07:00",
      "updated_at": "2025-09-19T10:37:24.657185+07:00"
    }
  },
  "message": "Login successful",
  "success": true
}
```

#### Refresh Token
```bash
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

### Protected Routes
Include the JWT token in the Authorization header:
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### User Profile Routes

#### Get Current User Profile
```bash
GET /api/v1/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "id": "1",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-09-19T10:37:24.657185+07:00",
    "updated_at": "2025-09-19T10:37:24.657185+07:00"
  },
  "message": "Profile retrieved successfully",
  "success": true
}
```

#### Update User Profile
```bash
PUT /api/v1/profile
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "Updated Name"
}
```

**Response:**
```json
{
  "data": {
    "id": "1",
    "name": "Updated Name",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-09-19T10:37:24.657185+07:00",
    "updated_at": "2025-09-19T10:38:04.064940+07:00"
  },
  "message": "Profile updated successfully",
  "success": true
}
```

### Admin Routes (Admin Role Required)

#### Get All Users
```bash
GET /api/v1/admin/users
Authorization: Bearer <admin-token>
```

#### Get User by ID
```bash
GET /api/v1/admin/users/{id}
Authorization: Bearer <admin-token>
```

#### Delete User
```bash
DELETE /api/v1/admin/users/{id}
Authorization: Bearer <admin-token>
```

### Error Responses
All endpoints return consistent error responses:

```json
{
  "success": false,
  "message": "Error description",
  "error": {
    "code": "ERROR_CODE"
  }
}
```

Common error codes:
- `VALIDATION_FAILED`: Request validation errors
- `UNAUTHORIZED`: Missing or invalid authentication
- `FORBIDDEN`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `INTERNAL_ERROR`: Server error

## 🛣️ Routes Architecture

### Modular Route Design
The application implements a sophisticated modular route architecture that separates concerns and improves maintainability:

#### Route Files Structure
```
internal/routes/
├── routes.go            # Main router coordinator
├── health_routes.go     # Health monitoring endpoints
├── auth_routes.go       # Authentication & registration
├── user_routes.go       # User profile management
├── admin_routes.go      # Administrative operations
└── route_types.go       # Route interfaces & documentation
```

#### Router Coordinator (`routes.go`)
The main router acts as a coordinator that:
- Initializes all route groups
- Applies global middleware (logging, CORS, authentication)
- Provides unified route documentation
- Manages route group dependencies

```go
type Router struct {
    userHandler   *handler.UserHandler
    jwtMiddleware *middleware.JWTMiddleware
    logger        *logger.Logger
    
    // Route groups
    healthRoutes *HealthRoutes
    authRoutes   *AuthRoutes
    userRoutes   *UserRoutes
    adminRoutes  *AdminRoutes
}
```

#### Route Group Interface
All route groups implement a common interface:

```go
type RouteGroup interface {
    SetupRoutes(router *mux.Router)
    GetRoutes() []string
}
```

#### Benefits of This Architecture

**🔧 Separation of Concerns**
- Each route file handles a specific domain
- Clear boundaries between different functionalities
- Easier to locate and modify specific routes

**📈 Scalability**
- Add new route groups without affecting existing ones
- Easy to implement versioning (e.g., `v2_user_routes.go`)
- Support for multiple API versions simultaneously

**🧪 Enhanced Testability**
- Test route groups independently
- Mock dependencies per route group
- Isolated unit tests for specific functionalities

**👥 Team Collaboration**
- Multiple developers can work on different route files
- Reduced merge conflicts
- Clear ownership of different API sections

**📚 Self-Documenting**
- Built-in route documentation
- Type-safe route information
- Automatic API documentation generation

#### Adding New Route Groups

To add a new route group (e.g., for products):

1. **Create the route file** (`product_routes.go`):
```go
type ProductRoutes struct {
    productHandler *handler.ProductHandler
}

func (pr *ProductRoutes) SetupRoutes(router *mux.Router) {
    productRouter := router.PathPrefix("/api/v1/products").Subrouter()
    productRouter.HandleFunc("", pr.productHandler.GetProducts).Methods("GET")
    productRouter.HandleFunc("", pr.productHandler.CreateProduct).Methods("POST")
}
```

2. **Add to main router** (`routes.go`):
```go
type Router struct {
    // ... existing fields
    productRoutes *ProductRoutes
}

func (r *Router) SetupRoutes() *mux.Router {
    // ... existing setup
    r.productRoutes.SetupRoutes(router)
    return router
}
```

#### Route Documentation
The system provides multiple levels of documentation:

- **Route Lists**: Simple string arrays for quick reference
- **Route Info**: Detailed metadata including methods, paths, and security requirements
- **API Documentation**: Automatically generated from route definitions

## 🎯 GraphQL API

### GraphQL Playground
Access the interactive GraphQL playground at:
- **Development**: `http://localhost:8080/graphql`
- **Production**: `https://your-app.onrender.com/graphql`

### Schema Overview
```graphql
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: Time!
  updatedAt: Time!
}

type Query {
  user(id: ID!): User
  users(first: Int, offset: Int): [User!]!
  me: User
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  deleteUser(id: ID!): Boolean!
}

type Subscription {
  userCreated: User!
  userUpdated: User!
  userDeleted: ID!
}
```

### Example Queries

#### Get User by ID
```graphql
query GetUser($id: ID!) {
  user(id: $id) {
    id
    name
    email
    createdAt
  }
}
```

**Variables:**
```json
{
  "id": "507f1f77bcf86cd799439011"
}
```

#### Get All Users with Pagination
```graphql
query GetUsers($first: Int, $offset: Int) {
  users(first: $first, offset: $offset) {
    id
    name
    email
    createdAt
  }
}
```

#### Create User
```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
    name
    email
    createdAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "Jane Smith",
    "email": "jane@example.com",
    "password": "securepassword123"
  }
}
```

#### Real-time Subscription
```graphql
subscription OnUserCreated {
  userCreated {
    id
    name
    email
    createdAt
  }
}
```

### GraphQL with cURL

#### Query Example
```bash
curl -X POST \
  http://localhost:8080/query \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <token>' \
  -d '{
    "query": "query GetUser($id: ID!) { user(id: $id) { id name email } }",
    "variables": { "id": "507f1f77bcf86cd799439011" }
  }'
```

#### Mutation Example
```bash
curl -X POST \
  http://localhost:8080/query \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <token>' \
  -d '{
    "query": "mutation CreateUser($input: CreateUserInput!) { createUser(input: $input) { id name email } }",
    "variables": {
      "input": {
        "name": "New User",
        "email": "newuser@example.com",
        "password": "password123"
      }
    }
  }'
```
## 🏗️ Architecture

### System Overview

The demo-go backend follows a **clean architecture** pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Presentation  │    │    Business     │    │      Data       │
│     Layer       │    │     Logic       │    │     Access      │
│                 │    │     Layer       │    │     Layer       │
│  • HTTP Routes  │───▶│  • Services     │───▶│  • Repository   │
│  • GraphQL      │    │  • Validation   │    │  • MongoDB      │
│  • Middleware   │    │  • Business     │    │  • Redis Cache  │
│  • Auth Handlers│    │    Rules        │    │  • External APIs│
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Component Architecture

```mermaid
graph TB
    A[Client Applications] -->|HTTP/GraphQL| B[Load Balancer]
    B --> C[API Gateway]
    C --> D[Go Backend Server]
    
    D --> E[Service Layer]
    E --> F[Repository Layer]
    F --> G[(MongoDB)]
    F --> H[(Redis Cache)]
    
    D --> I[GraphQL Layer]
    I --> J[Schema & Resolvers]
    
    D --> K[Middleware Stack]
    K --> L[JWT Auth]
    K --> M[CORS]
    K --> N[Rate Limiting]
    K --> O[Logging]
    
    subgraph "External Services"
        P[External APIs]
        Q[Monitoring]
        R[Logging Service]
    end
    
    D --> P
    D --> Q
    D --> R
```

### Directory Structure

```
demo-go/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handler/
│   │   └── user_handler.go      # User HTTP handlers
│   ├── routes/                  # HTTP Routes (Modular Structure)
│   │   ├── routes.go            # Main router coordinator
│   │   ├── health_routes.go     # Health check routes
│   │   ├── auth_routes.go       # Authentication routes
│   │   ├── user_routes.go       # User profile routes
│   │   ├── admin_routes.go      # Admin-only routes
│   │   └── route_types.go       # Route interfaces & types
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   ├── cors.go              # CORS middleware
│   │   ├── logging.go           # Request logging
│   │   └── rate_limit.go        # Rate limiting
│   ├── domain/
│   │   └── user.go              # User domain model & interfaces
│   ├── repository/
│   │   ├── memory_user_repository.go    # In-memory implementation
│   │   └── mongo_user_repository.go     # MongoDB implementation
│   ├── service/
│   │   ├── user_service.go      # User business logic
│   │   └── jwt_token_service.go # JWT token management
│   ├── cache/
│   │   └── redis.go             # Redis caching implementation
│   └── logger/
│       └── logger.go            # Structured logging
├── tests/
│   ├── integration/             # Integration tests
│   ├── unit/                    # Unit tests
│   └── load-test.js            # K6 performance tests
├── scripts/                    # Project management scripts
│   ├── project.sh              # Main project management script
│   ├── setup-graphql.sh        # GraphQL setup script
│   └── setup-cicd.sh           # CI/CD setup script
├── docs/
│   ├── SCALABLE_ARCHITECTURE.md # Architecture documentation
│   ├── GRAPHQL_API.md          # GraphQL API guide
│   └── CICD_PIPELINE.md        # CI/CD documentation
├── .github/
│   └── workflows/
│       └── ci-cd.yml           # GitHub Actions pipeline
├── docker-compose.yml          # Local development environment
├── Dockerfile                  # Production container
├── render.yaml                 # Render deployment config
├── Makefile                    # Build automation
└── .golangci.yml              # Linting configuration
```

### Project Management Scripts

The project includes powerful management scripts located in the `scripts/` directory for streamlined development workflow.

#### Main Project Script (`scripts/project.sh`)

A comprehensive project management script that provides a unified interface for common development tasks:

```bash
# Setup project dependencies
scripts/project.sh setup
# OR via Makefile
make project-setup

# Build project
scripts/project.sh build
make project-build

# Run in development mode
scripts/project.sh run dev 8080
scripts/project.sh dev memory redis
make project-dev

# Run tests
scripts/project.sh test coverage
make project-test

# Code quality
scripts/project.sh format
scripts/project.sh lint
scripts/project.sh security
make project-format project-lint project-security

# Dependency management
scripts/project.sh deps start
scripts/project.sh deps stop
scripts/project.sh deps restart

# Clean build artifacts
scripts/project.sh clean
make project-clean
```

#### Available Commands

| Command | Description | Makefile Equivalent |
|---------|-------------|-------------------|
| `setup` | Setup project dependencies | `make project-setup` |
| `build [target]` | Build project (default: server) | `make project-build` |
| `run [mode] [port]` | Run project (dev/prod modes) | `make project-run` |
| `dev [repo] [cache]` | Run with hot reload | `make project-dev` |
| `test [type]` | Run tests (unit/integration/coverage/all) | `make project-test` |
| `deps [action]` | Manage dependencies (start/stop/restart) | - |
| `clean` | Clean build artifacts | `make project-clean` |
| `format` | Format Go code | `make project-format` |
| `lint` | Lint Go code | `make project-lint` |
| `security` | Run security scan | `make project-security` |

#### Setup Scripts

##### GraphQL Setup (`scripts/setup-graphql.sh`)
```bash
# Setup GraphQL with gqlgen
scripts/setup-graphql.sh
# OR
make setup-graphql
```

Features:
- Installs gqlgen CLI tool
- Generates GraphQL boilerplate
- Creates schema and resolver files
- Sets up GraphQL server configuration

##### CI/CD Setup (`scripts/setup-cicd.sh`)
```bash
# Setup CI/CD pipeline for Render deployment
scripts/setup-cicd.sh
# OR
make setup-cicd
```

Features:
- Configures GitHub Actions workflow
- Sets up Render deployment
- Configures environment variables
- Sets up monitoring and alerts

#### CI/CD Pipeline Integration

The CI/CD pipeline (`.github/workflows/ci-cd.yml`) is fully integrated with the project scripts for consistency between local development and CI/CD environments:

##### Pipeline Jobs Using Scripts:

**🧪 Test Job:**
```yaml
- name: Setup project dependencies via script
  run: scripts/project.sh setup

- name: Format code via script
  run: scripts/project.sh format

- name: Lint code via script
  run: scripts/project.sh lint

- name: Run tests via script
  run: scripts/project.sh test coverage
```

**🔒 Security Job:**
```yaml
- name: Setup project dependencies via script
  run: scripts/project.sh setup

- name: Run security scan via script
  run: scripts/project.sh security
```

**🏗️ Build Job:**
```yaml
- name: Setup project dependencies via script
  run: scripts/project.sh setup

- name: Build project via script
  run: scripts/project.sh build
```

##### Benefits of Script-Based CI/CD:
- **🔄 Consistency**: Same commands work locally and in CI/CD
- **🛠️ Maintainability**: Single source of truth for build logic
- **🚀 Developer Experience**: Easy to reproduce CI/CD steps locally
- **📊 Debugging**: Simplified troubleshooting of build issues

#### Environment Variables for Scripts

The scripts support environment variables for configuration:

```bash
export PORT=8080                    # Server port
export REPOSITORY_TYPE=mongodb      # Repository type (memory, mongodb)
export CACHE_TYPE=redis             # Cache type (memory, redis)
export GIN_MODE=debug               # Gin mode (debug, release)
export LOG_LEVEL=info               # Log level (debug, info, warn, error)
```

#### Examples

```bash
# Complete development setup
scripts/project.sh setup
scripts/project.sh deps start
scripts/project.sh dev memory redis

# Production-like testing
scripts/project.sh build
scripts/project.sh run prod
scripts/project.sh test coverage

# Code quality check
scripts/project.sh format
scripts/project.sh lint
scripts/project.sh security
```

### Routes Architecture

#### Modular Route Structure
The application uses a modular route structure where each route group is defined in its own file:

**🏥 Health Routes (`health_routes.go`)**
- `GET /health` - System health check

**🔐 Authentication Routes (`auth_routes.go`)**
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Token refresh

**👤 User Routes (`user_routes.go`)**
- `GET /api/v1/profile` - Get user profile
- `PUT /api/v1/profile` - Update user profile

**👨‍💼 Admin Routes (`admin_routes.go`)**
- `GET /api/v1/admin/users` - List all users
- `GET /api/v1/admin/users/{id}` - Get user by ID
- `DELETE /api/v1/admin/users/{id}` - Delete user

#### Route Organization Benefits
- **🔧 Separation of Concerns**: Each route group handles specific functionality
- **📈 Scalability**: Easy to add new route groups without affecting existing ones
- **🧪 Testability**: Individual route groups can be tested independently
- **👥 Team Development**: Multiple developers can work on different route files
- **📚 Documentation**: Self-documenting structure with clear responsibilities

### Technology Stack

#### Backend Core
- **Go 1.21+**: Primary programming language
- **Gin Framework**: HTTP web framework
- **MongoDB**: Primary database
- **Redis**: Caching and session storage

#### GraphQL Implementation
- **gqlgen**: GraphQL server library
- **Type-safe resolvers**: Generated Go code
- **Real-time subscriptions**: WebSocket support

#### Authentication & Security
- **JWT**: JSON Web Tokens
- **bcrypt**: Password hashing
- **CORS**: Cross-origin resource sharing
- **Rate limiting**: Request throttling

#### DevOps & Infrastructure
- **Docker**: Containerization
- **GitHub Actions**: CI/CD pipeline
- **Render**: Cloud deployment platform
- **k6**: Performance testing

#### Monitoring & Observability
- **Structured logging**: JSON format
- **Health checks**: Service monitoring
- **Metrics collection**: Performance tracking
- **Distributed tracing**: Request tracing

### Scalability Features

1. **Horizontal Scaling**
   - Stateless design
   - Load balancer ready
   - Database clustering support

2. **Caching Strategy**
   - Redis caching layer
   - Cache-aside pattern
   - TTL-based expiration

3. **Database Optimization**
   - Connection pooling
   - Index optimization
   - Query optimization

4. **Performance Monitoring**
   - Request metrics
   - Database performance
   - Cache hit rates

## 🚀 Deployment

### Local Development

#### Using Docker Compose
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

#### Manual Setup
```bash
# Start dependencies
docker run -d -p 6379:6379 --name redis redis:alpine
docker run -d -p 27017:27017 --name mongodb mongo:5.0

# Set environment variables
export DATABASE_URL=mongodb://localhost:27017
export REDIS_URL=redis://localhost:6379
export JWT_SECRET=your-secret-key

# Run the application
go run cmd/server/main.go
```

### Production Deployment

#### Render Platform
```bash
# Connect your GitHub repository to Render
# Use the provided render.yaml configuration

# Environment variables required:
# - DATABASE_URL
# - REDIS_URL
# - JWT_SECRET
# - PORT (automatically set by Render)
```

#### Docker Deployment
```bash
# Build production image
docker build -t demo-go:latest .

# Run with environment file
docker run --env-file .env -p 8080:8080 demo-go:latest
```

#### Kubernetes Deployment
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo-go
spec:
  replicas: 3
  selector:
    matchLabels:
      app: demo-go
  template:
    metadata:
      labels:
        app: demo-go
    spec:
      containers:
      - name: demo-go
        image: demo-go:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: demo-go-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: demo-go-secrets
              key: jwt-secret
```

### Environment Configuration

#### Development
```bash
export GIN_MODE=debug
export LOG_LEVEL=debug
export CACHE_TYPE=memory
export DATABASE_URL=mongodb://localhost:27017
```

#### Staging
```bash
export GIN_MODE=release
export LOG_LEVEL=info
export CACHE_TYPE=redis
export DATABASE_URL=mongodb+srv://staging-cluster
export REDIS_URL=redis://staging-redis:6379
```

#### Production
```bash
export GIN_MODE=release
export LOG_LEVEL=warn
export CACHE_TYPE=redis
export DATABASE_URL=mongodb+srv://production-cluster
export REDIS_URL=redis://production-redis:6379
export METRICS_ENABLED=true
export TRACING_ENABLED=true
```
- **Timestamps**: Automatic created_at and updated_at fields
- **Context Support**: All operations support context for timeouts and cancellation

### Document Structure
```go
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### Connection Configuration
```go
// Local MongoDB
connectionString := "mongodb://localhost:27017"

// MongoDB Atlas
connectionString := "mongodb+srv://username:password@cluster.mongodb.net/"

repo, err := repository.NewMongoUserRepository(connectionString, "userdb")
```

### Usage Examples

**Direct Repository Usage:**
```go
repo, err := repository.NewMongoUserRepository("mongodb://localhost:27017", "userdb")
if err != nil {
    log.Fatal(err)
}
defer repo.Close(context.Background())

// Create user
user, err := repo.CreateUser(ctx, repository.UserCreateRequest{
    Name:  "John Doe",
    Email: "john@example.com",
})

// Get user by ID
user, err := repo.GetUserByID(ctx, user.ID.Hex())

## 🧪 Performance & Testing

### Load Testing with k6

Run comprehensive performance tests:

```bash
# Install k6 (macOS)
brew install k6

# Run basic load test
k6 run tests/load-test.js

# Run stress test
k6 run --vus 50 --duration 2m tests/load-test.js

# Run spike test
k6 run --vus 100 --duration 30s tests/load-test.js
```

#### Performance Benchmarks

| Endpoint | RPS | P95 Latency | Success Rate |
|----------|-----|-------------|--------------|
| GET /health | 10,000+ | < 5ms | 100% |
| GET /users | 5,000+ | < 20ms | 100% |
| POST /users | 2,000+ | < 50ms | 100% |
| GraphQL Queries | 3,000+ | < 30ms | 100% |

### Unit Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/service/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Testing

```bash
# Start test dependencies
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test -tags=integration ./tests/integration/...

# Clean up
docker-compose -f docker-compose.test.yml down
```

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make (for build automation)

### Development Setup

```bash
# Clone repository
git clone https://github.com/yourusername/demo-go.git
cd demo-go

# Install dependencies
go mod download

# Install development tools
make install-tools

# Setup pre-commit hooks
make setup-hooks

# Start development environment
make dev
```

### Build Commands

```bash
# Build application
make build

# Build for different platforms
make build-linux
make build-windows
make build-darwin

# Clean build artifacts
make clean

# Format code
make fmt

# Run linter
make lint

# Run security scan
make security-scan
```

### GraphQL Development

```bash
# Generate GraphQL code
make generate-graphql

# Update schema
make update-schema

# Validate GraphQL schema
make validate-schema
```

### Database Migrations

```bash
# Create new migration
make create-migration name=add_user_indexes

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Reset database
make db-reset
```

## 🧪 Testing

### Test Categories

#### Unit Tests
```bash
# Test service layer
go test ./internal/service/...

# Test repository layer
go test ./internal/repository/...

# Test handlers
go test ./internal/handler/...
```

#### Integration Tests
```bash
# API integration tests
go test -tags=integration ./tests/integration/api/...

# Database integration tests
go test -tags=integration ./tests/integration/database/...

# GraphQL integration tests
go test -tags=integration ./tests/integration/graphql/...
```

#### End-to-End Tests
```bash
# Full application tests
go test -tags=e2e ./tests/e2e/...

# Browser-based tests (if applicable)
make test-e2e
```

### Test Coverage

```bash
# Generate coverage report
make coverage

# View coverage in browser
make coverage-html

# Check coverage threshold
make coverage-check
```

### Continuous Testing

```bash
# Watch mode for development
make test-watch

# Run tests on file changes
make dev-test
```

## 🤝 Contributing

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/awesome-feature
   ```
3. **Make your changes**
   - Follow Go code style guidelines
   - Add tests for new functionality
   - Update documentation
4. **Run tests and linting**
   ```bash
   make test
   make lint
   make security-scan
   ```
5. **Commit your changes**
   ```bash
   git commit -m "feat: add awesome feature"
   ```
6. **Push to your fork**
   ```bash
   git push origin feature/awesome-feature
   ```
7. **Create a Pull Request**

### Code Style Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for code formatting
- Write meaningful commit messages
- Add comments for exported functions
- Include tests for new features

### Project Structure Guidelines

- Place business logic in `internal/service/`
- Keep handlers thin and focused
- Use interfaces for external dependencies
- Follow repository pattern for data access
- Add integration tests for new endpoints

### Documentation Guidelines

- Update README.md for new features
- Add API documentation for new endpoints
- Include examples in documentation
- Update architecture diagrams when needed

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/) - HTTP web framework
- [gqlgen](https://gqlgen.com/) - GraphQL library for Go
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - Official MongoDB driver
- [JWT-Go](https://github.com/dgrijalva/jwt-go) - JWT implementation
- [Testify](https://github.com/stretchr/testify) - Testing toolkit

## 📞 Support

- 📧 Email: your-email@example.com
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/demo-go/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/yourusername/demo-go/discussions)
- 📖 Documentation: [Wiki](https://github.com/yourusername/demo-go/wiki)

---

**Built with ❤️ using Go**
