# GraphQL Implementation with gqlgen

This project now includes a comprehensive GraphQL API implementation using gqlgen for User management operations.

## 🚀 Features

- **Type-safe GraphQL**: Automatically generated resolvers and models
- **User Operations**: Complete CRUD operations for user management
- **Real-time Subscriptions**: WebSocket support for live updates
- **GraphQL Playground**: Interactive query interface for development
- **Structured Logging**: Comprehensive logging for all GraphQL operations
- **Authentication**: JWT token support for secured operations
- **Caching Integration**: Works with existing Redis caching layer

## 📋 Prerequisites

- **Go 1.21+**
- **gqlgen**: GraphQL code generator for Go

### Installing gqlgen

```bash
# Install gqlgen CLI tool
go install github.com/99designs/gqlgen@latest

# Add gqlgen dependencies to your project
go get github.com/99designs/gqlgen/graphql
go get github.com/99designs/gqlgen/graphql/handler
go get github.com/99designs/gqlgen/graphql/playground
```

## 🏗️ Generated Files

Run gqlgen to generate the required GraphQL code:

```bash
# Generate GraphQL resolvers and models
gqlgen generate

# Or use the shorthand
gqlgen
```

This will generate:
- `internal/graphql/generated.go` - Generated GraphQL server code
- `internal/graphql/models_gen.go` - Generated model structs
- Updates to `internal/graphql/resolver.go` - Resolver interfaces

## 📊 GraphQL Schema

### User Type
```graphql
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: Time!
  updatedAt: Time!
}
```

### Queries
```graphql
type Query {
  # Get a single user by ID
  getUser(id: ID!): User
  
  # Get all users with pagination
  getUsers(limit: Int, offset: Int): [User!]!
  
  # Search users by name or email
  searchUsers(query: String!): [User!]!
  
  # Get current authenticated user
  me: User
}
```

### Mutations
```graphql
type Mutation {
  # Create a new user
  createUser(input: CreateUserInput!): User!
  
  # Update an existing user
  updateUser(id: ID!, input: UpdateUserInput!): User!
  
  # Delete a user
  deleteUser(id: ID!): Boolean!
}
```

### Subscriptions
```graphql
type Subscription {
  # Subscribe to user creation events
  userCreated: User!
  
  # Subscribe to user updates
  userUpdated: User!
  
  # Subscribe to user deletion events
  userDeleted: ID!
}
```

## 🔧 Integration with Server

Add GraphQL routes to your main server in `cmd/server/main.go`:

```go
// Add to setupRoutes function
func setupRoutes(userHandler *handler.UserHandler, cacheHandler *handler.CacheHandler, jwtMiddleware *middleware.JWTMiddleware, baseLogger *logger.Logger) *mux.Router {
    router := mux.NewRouter()
    
    // ... existing routes ...
    
    // GraphQL endpoints
    graphqlHandler := handler.NewGraphQLHandler(userService)
    
    // GraphQL API endpoint
    router.Handle("/graphql", graphqlHandler.Handler()).Methods("GET", "POST")
    
    // GraphQL Playground (development only)
    if os.Getenv("ENVIRONMENT") != "production" {
        router.Handle("/playground", graphqlHandler.PlaygroundHandler()).Methods("GET")
    }
    
    return router
}
```

## 🎯 Example Queries

### Get User by ID
```graphql
query GetUser($id: ID!) {
  getUser(id: $id) {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

**Variables:**
```json
{
  "id": "user-uuid-here"
}
```

### Get All Users with Pagination
```graphql
query GetUsers($limit: Int, $offset: Int) {
  getUsers(limit: $limit, offset: $offset) {
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
  "limit": 10,
  "offset": 0
}
```

### Search Users
```graphql
query SearchUsers($query: String!) {
  searchUsers(query: $query) {
    id
    name
    email
  }
}
```

**Variables:**
```json
{
  "query": "john"
}
```

### Get Current User
```graphql
query Me {
  me {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

## 🔄 Example Mutations

### Create User
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
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### Update User
```graphql
mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) {
    id
    name
    email
    updatedAt
  }
}
```

**Variables:**
```json
{
  "id": "user-uuid-here",
  "input": {
    "name": "John Smith"
  }
}
```

### Delete User
```graphql
mutation DeleteUser($id: ID!) {
  deleteUser(id: $id)
}
```

**Variables:**
```json
{
  "id": "user-uuid-here"
}
```

## 🔗 Subscriptions

### Subscribe to User Creation
```graphql
subscription UserCreated {
  userCreated {
    id
    name
    email
    createdAt
  }
}
```

### Subscribe to User Updates
```graphql
subscription UserUpdated {
  userUpdated {
    id
    name
    email
    updatedAt
  }
}
```

### Subscribe to User Deletions
```graphql
subscription UserDeleted {
  userDeleted
}
```

## 🖥️ GraphQL Playground

Access the interactive GraphQL Playground at:
```
http://localhost:8080/playground
```

The playground provides:
- **Query Editor**: Write and test GraphQL queries
- **Schema Explorer**: Browse available types and operations
- **Variables Panel**: Input variables for parameterized queries
- **Response Viewer**: View query results and errors
- **Documentation**: Auto-generated API documentation

## 🔐 Authentication

GraphQL endpoints support JWT authentication:

```bash
# Include JWT token in Authorization header
curl -X POST http://localhost:8080/graphql \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query Me { me { id name email } }"
  }'
```

## 📊 Performance Features

### Caching Integration
- **Automatic Caching**: Queries utilize existing Redis cache layer
- **Cache Invalidation**: Mutations automatically invalidate related cache entries
- **Performance Metrics**: Cache hit/miss rates logged for optimization

### Query Complexity Limiting
- **Complexity Analysis**: Prevents expensive queries from overloading the server
- **Configurable Limits**: Set maximum query complexity per request
- **Depth Limiting**: Control maximum query nesting depth

### Connection Pooling
- **Database Connections**: Reuse existing connection pools
- **Cache Connections**: Leverage Redis connection pooling
- **Resource Management**: Efficient resource utilization

## 🔍 Monitoring and Logging

### Request Logging
All GraphQL operations are logged with:
- **Operation Name**: Query/mutation/subscription name
- **Query Content**: Full GraphQL query text
- **Variables**: Input variables (sanitized)
- **Execution Time**: Query processing duration
- **Error Details**: Comprehensive error information

### Metrics Integration
- **Request Count**: Total GraphQL requests
- **Response Time**: Average query execution time
- **Error Rate**: Percentage of failed operations
- **Cache Hit Rate**: Percentage of cached vs database queries

## 🚀 Production Considerations

### Security
- **Input Validation**: Automatic validation of all input types
- **Rate Limiting**: Prevent abuse with request throttling
- **Query Whitelisting**: Allow only pre-approved queries in production
- **Introspection Disabled**: Disable schema introspection in production

### Performance
- **Query Batching**: Support multiple queries in single request
- **Persistent Queries**: Cache frequently used queries
- **Response Caching**: Cache query results at the GraphQL layer
- **Field-level Caching**: Cache individual field results

### Monitoring
- **APM Integration**: Connect to New Relic, DataDog, or similar
- **Custom Metrics**: Track business-specific GraphQL metrics
- **Error Tracking**: Comprehensive error monitoring and alerting
- **Performance Dashboards**: Real-time GraphQL performance visualization

## 🧪 Testing

### Unit Tests
```bash
# Test GraphQL resolvers
go test ./internal/graphql/...

# Test with coverage
go test -cover ./internal/graphql/...
```

### Integration Tests
```bash
# Test GraphQL endpoints
go test -tags=integration ./internal/handler/...
```

### Example Test Queries
```bash
# Test query with curl
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { getUsers(limit: 5) { id name email } }"
  }'
```

## 📚 Additional Resources

- [gqlgen Documentation](https://gqlgen.com/)
- [GraphQL Specification](https://graphql.org/learn/)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [Apollo GraphQL Guide](https://www.apollographql.com/docs/)

---

🎉 **Your Go web server now has a full-featured GraphQL API!** Use the playground to explore and test your GraphQL operations.
