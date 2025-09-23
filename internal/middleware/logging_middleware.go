package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"demo-go/internal/logger"

	"github.com/google/uuid"
)

// Context key types to avoid collisions
type loggingContextKey string

const (
	requestIDKey loggingContextKey = "request_id"
)

// LoggingMiddleware provides request logging with structured output
func LoggingMiddleware(baseLogger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := generateRequestID(r)

			// Create logger for this request
			log := baseLogger.ForRequest(r.Method, r.URL.Path, requestID)

			// Add request ID to context for downstream use
			ctx := r.Context()
			ctx = requestIDContext(ctx, requestID)
			r = r.WithContext(ctx)

			// Add request ID header to response
			w.Header().Set("X-Request-ID", requestID)

			// Capture request body for JSON logging
			var requestBody []byte
			if r.Body != nil && shouldLogBody(r) {
				requestBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			// Create a response writer wrapper to capture status code, size, and body
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				size:           0,
				body:           &bytes.Buffer{},
			}

			// Log incoming request (only for non-health checks to reduce noise)
			if r.URL.Path != "/health" {
				logMessage := fmt.Sprintf("→ Request started\nMethod: %s\nPath: %s\nUser-Agent: %s\nClient-IP: %s",
					r.Method, r.URL.Path, r.UserAgent(), getClientIP(r))
				
				// Add pretty JSON request body if present
				if len(requestBody) > 0 {
					if prettyJSON := formatJSON(requestBody); prettyJSON != "" {
						logMessage += fmt.Sprintf("\nRequest Body:\n%s", prettyJSON)
					}
				}
				
				log.ConsoleInfo(logMessage)
			}

			// Call next handler
			next.ServeHTTP(wrapper, r)

			// Log completed request
			duration := time.Since(start)
			
			// Choose appropriate log level based on status code
			statusEmoji := getStatusEmoji(wrapper.statusCode)
			
			if r.URL.Path == "/health" {
				// Minimal logging for health checks
				log.ConsoleDebug(fmt.Sprintf("✓ Health check - Status: %d, Duration: %v", wrapper.statusCode, duration.Round(time.Microsecond)))
			} else {
				logMessage := fmt.Sprintf("← Request completed %s\nStatus: %d\nDuration: %v\nSize: %s",
					statusEmoji, wrapper.statusCode, duration.Round(time.Microsecond), formatBytes(wrapper.size))
				
				// Add pretty JSON response body if present
				if wrapper.body.Len() > 0 {
					if prettyJSON := formatJSON(wrapper.body.Bytes()); prettyJSON != "" {
						logMessage += fmt.Sprintf("\nResponse Body:\n%s", prettyJSON)
					}
				}
				
				log.ConsoleInfo(logMessage)
			}
		})
	}
}

// generateRequestID creates or extracts a request ID
func generateRequestID(r *http.Request) string {
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}
	return uuid.New().String()
}

// requestIDContext adds request ID to context
func requestIDContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code and response size
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int64
	body       *bytes.Buffer
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	// Capture response body
	if w.body != nil {
		w.body.Write(data)
	}
	
	size, err := w.ResponseWriter.Write(data)
	w.size += int64(size)
	return size, err
}

// CORSMiddleware provides CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the real client IP from various headers
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// getStatusEmoji returns an emoji based on HTTP status code
func getStatusEmoji(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "✅" // Success
	case statusCode >= 300 && statusCode < 400:
		return "↩️" // Redirect
	case statusCode >= 400 && statusCode < 500:
		return "⚠️" // Client error
	case statusCode >= 500:
		return "🚨" // Server error
	default:
		return "ℹ️" // Info
	}
}

// formatBytes formats byte count in human readable format
func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
	}
}

// shouldLogBody determines if we should capture and log the request body
func shouldLogBody(r *http.Request) bool {
	// Only log JSON content types
	contentType := r.Header.Get("Content-Type")
	return strings.Contains(contentType, "application/json") && 
		   r.ContentLength > 0 && r.ContentLength < 1024*10 // Max 10KB
}

// formatJSON formats JSON bytes into a pretty string
func formatJSON(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	
	// Try to parse and format as JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		// If not valid JSON, return as-is (truncated if too long)
		if len(data) > 500 {
			return string(data[:500]) + "..."
		}
		return string(data)
	}
	
	// Pretty print JSON with 2-space indentation
	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return string(data)
	}
	
	return string(prettyJSON)
}
