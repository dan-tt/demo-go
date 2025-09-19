package routes

import (
	"demo-go/internal/handler"
	"demo-go/internal/logger"
	"demo-go/internal/middleware"

	"github.com/gorilla/mux"
)

// Router holds the dependencies needed for route setup
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

// NewRouter creates a new router instance with dependencies
func NewRouter(userHandler *handler.UserHandler, jwtMiddleware *middleware.JWTMiddleware, logger *logger.Logger) *Router {
	return &Router{
		userHandler:   userHandler,
		jwtMiddleware: jwtMiddleware,
		logger:        logger,

		// Initialize route groups
		healthRoutes: NewHealthRoutes(userHandler),
		authRoutes:   NewAuthRoutes(userHandler),
		userRoutes:   NewUserRoutes(userHandler),
		adminRoutes:  NewAdminRoutes(userHandler, jwtMiddleware),
	}
}

// SetupRoutes configures all HTTP routes and returns the configured router
func (r *Router) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Add global middleware
	router.Use(middleware.LoggingMiddleware(r.logger))
	router.Use(middleware.CORSMiddleware)
	router.Use(r.jwtMiddleware.Authenticate)

	// Setup all route groups
	r.healthRoutes.SetupRoutes(router)
	r.authRoutes.SetupRoutes(router)
	r.userRoutes.SetupRoutes(router)
	r.adminRoutes.SetupRoutes(router)

	return router
}

// GetRoutesSummary returns a summary of all available routes
func (r *Router) GetRoutesSummary() map[string][]string {
	return map[string][]string{
		"Health Routes":         r.healthRoutes.GetRoutes(),
		"Authentication Routes": r.authRoutes.GetRoutes(),
		"User API Routes":       r.userRoutes.GetRoutes(),
		"Admin Routes":          r.adminRoutes.GetRoutes(),
	}
}
