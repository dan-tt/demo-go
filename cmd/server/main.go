package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"demo-go/internal/config"
	"demo-go/internal/cache"
	"demo-go/internal/domain"
	"demo-go/internal/handler"
	"demo-go/internal/logger"
	"demo-go/internal/middleware"
	"demo-go/internal/repository"
	"demo-go/internal/routes"
	"demo-go/internal/service"
)

func main() {
	// Initialize logger first
	loggerConfig := logger.DefaultConfig()
	if err := logger.InitGlobal(loggerConfig); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.GetGlobal().Sync()

	log := logger.GetGlobal().ForComponent("main")

	// Load configuration
	cfg := config.Load()
	
	log.Info("Starting Clean Architecture API server",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"environment", loggerConfig.Environment,
	)
	
	// Initialize dependencies
	server, cleanup, err := initializeServer(cfg, logger.GetGlobal())
	if err != nil {
		log.Fatal("Failed to initialize server", "error", err)
	}
	defer cleanup()
	
	// Start server
	go func() {
		log.Info("Server listening", 
			"address", fmt.Sprintf("http://%s:%s", cfg.Server.Host, cfg.Server.Port),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", "error", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Info("Shutting down server")
	
	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	
	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}
	
	log.Info("Server stopped gracefully")
}

// initializeServer sets up all dependencies and returns the HTTP server
func initializeServer(cfg *config.Config, baseLogger *logger.Logger) (*http.Server, func(), error) {
	log := baseLogger.ForComponent("server")
	
	// Choose repository implementation based on environment
	var userRepo domain.UserRepository
	var cleanup func() = func() {}
	
	repositoryType := os.Getenv("REPOSITORY_TYPE")
	if repositoryType == "memory" || repositoryType == "" {
		log.Info("Using in-memory repository")
		userRepo = repository.NewMemoryUserRepository()
	} else if repositoryType == "mongodb" {
		log.Info("Using MongoDB repository")
		
		// Initialize MongoDB client
		mongoClient, err := repository.NewMongoClient(cfg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
		}
		
		userRepo = repository.NewMongoUserRepository(mongoClient, cfg)
		
		// Setup cleanup function
		cleanup = func() {
			log.Info("Disconnecting from MongoDB")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := mongoClient.Disconnect(ctx); err != nil {
				log.Error("Error disconnecting from MongoDB", "error", err)
			} else {
				log.Info("Disconnected from MongoDB")
			}
		}
	} else {
		return nil, nil, fmt.Errorf("unsupported repository type: %s", repositoryType)
	}
	
	// Initialize services
	tokenService := service.NewJWTTokenService(cfg)
	baseUserService := service.NewUserService(userRepo, tokenService)
	
	// Initialize cache service
	var userService domain.UserService = baseUserService
	var cacheCleanup func() = func() {}
	
	cacheType := os.Getenv("CACHE_TYPE")
	if cacheType == "redis" {
		log.Info("Initializing Redis cache")
		cacheService, err := cache.NewRedisCache(cfg)
		if err != nil {
			log.Warn("Failed to initialize Redis cache, using service without cache", "error", err)
			// Continue without cache
		} else {
			log.Info("Redis cache initialized successfully")
			// Wrap the base service with caching
			userService = service.NewCachedUserService(baseUserService, cacheService, cfg.Cache.Redis.TTL)
			
			// Add cache cleanup
			cacheCleanup = func() {
				log.Info("Closing cache connection")
				if err := cacheService.Close(); err != nil {
					log.Error("Error closing cache connection", "error", err)
				} else {
					log.Info("Cache connection closed")
				}
			}
		}
	} else {
		log.Info("Cache disabled or not configured")
	}
	
	// Combine cleanup functions
	originalCleanup := cleanup
	cleanup = func() {
		cacheCleanup()
		originalCleanup()
	}
	
	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	
	// Initialize middleware
	jwtMiddleware := middleware.NewJWTMiddleware(tokenService)
	
	// Setup routes
	router := routes.NewRouter(userHandler, jwtMiddleware, baseLogger)
	httpRouter := router.SetupRoutes()
	
	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      httpRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
	
	return server, cleanup, nil
}
