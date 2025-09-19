package routes

import (
	"demo-go/internal/handler"

	"github.com/gorilla/mux"
)

// AuthRoutes handles authentication routes
type AuthRoutes struct {
	userHandler *handler.UserHandler
}

// NewAuthRoutes creates a new auth routes instance
func NewAuthRoutes(userHandler *handler.UserHandler) *AuthRoutes {
	return &AuthRoutes{
		userHandler: userHandler,
	}
}

// SetupRoutes configures authentication routes (public)
func (ar *AuthRoutes) SetupRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", ar.userHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", ar.userHandler.Login).Methods("POST")
	authRouter.HandleFunc("/refresh", ar.userHandler.RefreshToken).Methods("POST")
}

// GetRoutes returns a list of auth routes
func (ar *AuthRoutes) GetRoutes() []string {
	return []string{
		"POST /auth/register - User registration",
		"POST /auth/login - User login",
		"POST /auth/refresh - Refresh JWT token",
	}
}
