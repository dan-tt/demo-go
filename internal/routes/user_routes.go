package routes

import (
	"demo-go/internal/handler"

	"github.com/gorilla/mux"
)

// UserRoutes handles user-related API routes
type UserRoutes struct {
	userHandler *handler.UserHandler
}

// NewUserRoutes creates a new user routes instance
func NewUserRoutes(userHandler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		userHandler: userHandler,
	}
}

// SetupRoutes configures user API routes (authenticated)
func (ur *UserRoutes) SetupRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// User profile routes
	apiRouter.HandleFunc("/profile", ur.userHandler.GetProfile).Methods("GET")
	apiRouter.HandleFunc("/profile", ur.userHandler.UpdateProfile).Methods("PUT")
}

// GetRoutes returns a list of user routes
func (ur *UserRoutes) GetRoutes() []string {
	return []string{
		"GET /api/v1/profile - Get user profile",
		"PUT /api/v1/profile - Update user profile",
	}
}
