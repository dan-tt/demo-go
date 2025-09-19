package routes

import (
	"demo-go/internal/handler"
	"demo-go/internal/middleware"

	"github.com/gorilla/mux"
)

// AdminRoutes handles admin-only routes
type AdminRoutes struct {
	userHandler   *handler.UserHandler
	jwtMiddleware *middleware.JWTMiddleware
}

// NewAdminRoutes creates a new admin routes instance
func NewAdminRoutes(userHandler *handler.UserHandler, jwtMiddleware *middleware.JWTMiddleware) *AdminRoutes {
	return &AdminRoutes{
		userHandler:   userHandler,
		jwtMiddleware: jwtMiddleware,
	}
}

// SetupRoutes configures admin routes (admin only)
func (ar *AdminRoutes) SetupRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	adminRouter := apiRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(ar.jwtMiddleware.RequireAdmin)

	adminRouter.HandleFunc("/users", ar.userHandler.GetUsers).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", ar.userHandler.GetUserByID).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", ar.userHandler.DeleteUser).Methods("DELETE")
}

// GetRoutes returns a list of admin routes
func (ar *AdminRoutes) GetRoutes() []string {
	return []string{
		"GET /api/v1/admin/users - List all users",
		"GET /api/v1/admin/users/{id} - Get user by ID",
		"DELETE /api/v1/admin/users/{id} - Delete user",
	}
}
