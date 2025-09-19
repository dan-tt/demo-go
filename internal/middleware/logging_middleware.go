package middleware

import (
	"context"
	"net/http"
	"time"

	"demo-go/internal/logger"

	"github.com/google/uuid"
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

			// Create a response writer wrapper to capture status code and size
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				size:           0,
			}

			// Log incoming request
			log.Info("Request started",
				"user_agent", r.UserAgent(),
				"remote_addr", r.RemoteAddr,
				"content_length", r.ContentLength,
			)

			// Call next handler
			next.ServeHTTP(wrapper, r)

			// Log completed request
			duration := time.Since(start)
			log.Info("Request completed",
				"status_code", wrapper.statusCode,
				"duration_ms", duration.Milliseconds(),
				"response_size", wrapper.size,
			)
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
	return context.WithValue(ctx, "request_id", requestID)
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code and response size
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
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
