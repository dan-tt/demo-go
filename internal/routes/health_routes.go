package routes

import (
	"demo-go/internal/handler"

	"github.com/gorilla/mux"
)

// HealthRoutes handles health check routes
type HealthRoutes struct {
	userHandler *handler.UserHandler
}

// NewHealthRoutes creates a new health routes instance
func NewHealthRoutes(userHandler *handler.UserHandler) *HealthRoutes {
	return &HealthRoutes{
		userHandler: userHandler,
	}
}

// SetupRoutes configures health check routes (public)
func (hr *HealthRoutes) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/health", hr.userHandler.Health).Methods("GET")
}

// GetRoutes returns a list of health routes
func (hr *HealthRoutes) GetRoutes() []string {
	return []string{
		"GET /health - Health check",
	}
}
