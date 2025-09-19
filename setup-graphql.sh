#!/bin/bash

# GraphQL Setup Script
# This script sets up the GraphQL implementation using gqlgen

set -e

echo "🚀 GraphQL Setup with gqlgen"
echo "============================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo "Please install Go first: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed"

# Install gqlgen CLI tool
echo "📦 Installing gqlgen CLI tool..."
if ! command -v gqlgen &> /dev/null; then
    go install github.com/99designs/gqlgen@latest
    echo "✅ gqlgen CLI installed"
else
    echo "✅ gqlgen CLI already installed"
fi

# Add gqlgen to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin

# Verify gqlgen installation
if command -v gqlgen &> /dev/null; then
    echo "✅ gqlgen is available in PATH"
    gqlgen version
else
    echo "⚠️  gqlgen not found in PATH"
    echo "Add $(go env GOPATH)/bin to your PATH"
fi

echo ""
echo "📁 Project Structure:"
echo "   schema.graphql - GraphQL schema definition"
echo "   gqlgen.yml - gqlgen configuration"
echo "   internal/graphql/resolver.go - Resolver implementation"
echo "   internal/graphql/interfaces.go - Interface definitions"
echo "   internal/handler/graphql_handler.go - HTTP handler"

echo ""
echo "🔨 Generate GraphQL Code:"
echo "   Run: gqlgen generate"
echo "   This will create:"
echo "   • internal/graphql/generated.go"
echo "   • internal/graphql/models_gen.go"

echo ""
echo "🌐 Server Integration:"
echo "   Add GraphQL routes to cmd/server/main.go:"
echo ""
echo "   // GraphQL endpoints"
echo "   graphqlHandler := handler.NewGraphQLHandler(userService)"
echo "   router.Handle(\"/graphql\", graphqlHandler.Handler()).Methods(\"GET\", \"POST\")"
echo "   router.Handle(\"/playground\", graphqlHandler.PlaygroundHandler()).Methods(\"GET\")"

echo ""
echo "🎯 Available Endpoints:"
echo "   • POST /graphql - GraphQL API endpoint"
echo "   • GET /playground - GraphQL Playground (development)"

echo ""
echo "📊 Example Queries:"
echo ""
echo "   # Get User"
echo "   query {"
echo "     getUser(id: \"user-id\") {"
echo "       id name email"
echo "     }"
echo "   }"
echo ""
echo "   # Create User"
echo "   mutation {"
echo "     createUser(input: {name: \"John\", email: \"john@example.com\"}) {"
echo "       id name email"
echo "     }"
echo "   }"

echo ""
echo "🔧 Next Steps:"
echo "   1. Run: gqlgen generate"
echo "   2. Update server routes in cmd/server/main.go"
echo "   3. Start server: go run cmd/server/main.go"
echo "   4. Visit: http://localhost:8080/playground"
echo "   5. Test GraphQL queries and mutations"

echo ""
echo "📚 Documentation:"
echo "   • GraphQL API Guide: GRAPHQL_API.md"
echo "   • gqlgen Documentation: https://gqlgen.com/"
echo "   • GraphQL Spec: https://graphql.org/learn/"

echo ""
echo "🎉 GraphQL setup complete!"
