// Package middleware provides HTTP middleware functions for the demo-go application,
// including JWT authentication, authorization, logging, and request processing.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"demo-go/internal/domain"
)

// Context key types to avoid collisions
type contextKey string

const (
	userIDKey    contextKey = "user_id"
	userEmailKey contextKey = "user_email"
	userRoleKey  contextKey = "user_role"
)

// Helper functions to safely retrieve context values

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// GetUserEmailFromContext extracts the user email from the request context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}

// GetUserRoleFromContext extracts the user role from the request context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(userRoleKey).(string)
	return role, ok
}

// JWTMiddleware provides JWT authentication middleware
type JWTMiddleware struct {
	tokenService domain.TokenService
	skipPaths    map[string]bool
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(tokenService domain.TokenService) *JWTMiddleware {
	// Define paths that should skip authentication
	skipPaths := map[string]bool{
		"/health":        true,
		"/auth/register": true,
		"/auth/login":    true,
	}

	return &JWTMiddleware{
		tokenService: tokenService,
		skipPaths:    skipPaths,
	}
}

// Authenticate is a middleware that validates JWT tokens
func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for certain paths
		if m.shouldSkipPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		tokenString := m.extractTokenFromHeader(r)
		if tokenString == "" {
			m.writeUnauthorizedResponse(w, "Missing or invalid Authorization header")
			return
		}

		// Validate token
		claims, err := m.tokenService.ValidateToken(tokenString)
		if err != nil {
			m.writeUnauthorizedResponse(w, "Invalid or expired token")
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, userEmailKey, claims.Email)
		ctx = context.WithValue(ctx, userRoleKey, claims.Role)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole is a middleware that checks if user has required role
func (m *JWTMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value(userRoleKey)
			if userRole == nil {
				m.writeForbiddenResponse(w, "User role not found in context")
				return
			}

			roleStr, ok := userRole.(string)
			if !ok || roleStr != role {
				m.writeForbiddenResponse(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin is a middleware that checks if user is admin
func (m *JWTMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireRole("admin")(next)
}

// Helper methods

func (m *JWTMiddleware) shouldSkipPath(path string) bool {
	return m.skipPaths[path]
}

func (m *JWTMiddleware) extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check if header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}

	// Extract token part
	return strings.TrimSpace(authHeader[len(bearerPrefix):])
}

func (m *JWTMiddleware) writeUnauthorizedResponse(w http.ResponseWriter, message string) {
	m.writeJSONError(w, http.StatusUnauthorized, message, "UNAUTHORIZED")
}

func (m *JWTMiddleware) writeForbiddenResponse(w http.ResponseWriter, message string) {
	m.writeJSONError(w, http.StatusForbidden, message, "FORBIDDEN")
}

func (m *JWTMiddleware) writeJSONError(w http.ResponseWriter, statusCode int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := `{
		"success": false,
		"message": "` + message + `",
		"error": {
			"code": "` + code + `"
		}
	}`

	if _, err := w.Write([]byte(response)); err != nil {
		// Log the error but there's not much we can do at this point
		// since we're already in an error handling path
	}
}
