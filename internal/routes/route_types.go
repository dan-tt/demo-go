package routes

import "github.com/gorilla/mux"

// RouteGroup defines the interface that all route groups must implement
type RouteGroup interface {
	SetupRoutes(router *mux.Router)
	GetRoutes() []string
}

// RouteConfig holds configuration for route setup
type RouteConfig struct {
	Prefix      string
	Middleware  []mux.MiddlewareFunc
	Description string
}

// RouteInfo contains information about a single route
type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	Description string
	Protected   bool
	AdminOnly   bool
}

// GetAllRouteInfo returns detailed information about all routes
func (r *Router) GetAllRouteInfo() []RouteInfo {
	var routes []RouteInfo

	routes = append(routes, r.getHealthRoutes()...)
	routes = append(routes, r.getAuthRoutes()...)
	routes = append(routes, r.getUserRoutes()...)
	routes = append(routes, r.getAdminRoutes()...)

	return routes
}

// getHealthRoutes returns health check route information
func (r *Router) getHealthRoutes() []RouteInfo {
	return []RouteInfo{
		{
			Method:      "GET",
			Path:        "/health",
			Handler:     "userHandler.Health",
			Description: "Health check endpoint",
			Protected:   false,
			AdminOnly:   false,
		},
	}
}

// getAuthRoutes returns authentication route information
func (r *Router) getAuthRoutes() []RouteInfo {
	return []RouteInfo{
		{
			Method:      "POST",
			Path:        "/auth/register",
			Handler:     "userHandler.Register",
			Description: "User registration",
			Protected:   false,
			AdminOnly:   false,
		},
		{
			Method:      "POST",
			Path:        "/auth/login",
			Handler:     "userHandler.Login",
			Description: "User login",
			Protected:   false,
			AdminOnly:   false,
		},
		{
			Method:      "POST",
			Path:        "/auth/refresh",
			Handler:     "userHandler.RefreshToken",
			Description: "Refresh JWT token",
			Protected:   false,
			AdminOnly:   false,
		},
	}
}

// getUserRoutes returns user API route information
func (r *Router) getUserRoutes() []RouteInfo {
	return []RouteInfo{
		{
			Method:      "GET",
			Path:        "/api/v1/profile",
			Handler:     "userHandler.GetProfile",
			Description: "Get user profile",
			Protected:   true,
			AdminOnly:   false,
		},
		{
			Method:      "PUT",
			Path:        "/api/v1/profile",
			Handler:     "userHandler.UpdateProfile",
			Description: "Update user profile",
			Protected:   true,
			AdminOnly:   false,
		},
	}
}

// getAdminRoutes returns admin API route information
func (r *Router) getAdminRoutes() []RouteInfo {
	return []RouteInfo{
		{
			Method:      "GET",
			Path:        "/api/v1/admin/users",
			Handler:     "userHandler.GetUsers",
			Description: "List all users",
			Protected:   true,
			AdminOnly:   true,
		},
		{
			Method:      "GET",
			Path:        "/api/v1/admin/users/{id}",
			Handler:     "userHandler.GetUserByID",
			Description: "Get user by ID",
			Protected:   true,
			AdminOnly:   true,
		},
		{
			Method:      "DELETE",
			Path:        "/api/v1/admin/users/{id}",
			Handler:     "userHandler.DeleteUser",
			Description: "Delete user",
			Protected:   true,
			AdminOnly:   true,
		},
	}
}
