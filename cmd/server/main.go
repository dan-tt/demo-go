package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"demo-go/internal/cache"
	"demo-go/internal/config"
	"demo-go/internal/domain"
	"demo-go/internal/handler"
	"demo-go/internal/logger"
	"demo-go/internal/middleware"
	"demo-go/internal/repository"
	"demo-go/internal/routes"
	"demo-go/internal/service"
)

// MongoDB disconnect timeout
const MongoDisconnectTimeout = 10 * time.Second

func main() {
	// Initialize logger first
	loggerConfig := logger.DefaultConfig()
	if err := logger.InitGlobal(loggerConfig); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := logger.GetGlobal().Sync(); err != nil {
			// Log sync failed, but we're exiting anyway
			fmt.Printf("Failed to sync logger: %v\n", err)
		}
	}()

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
		log.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	// Start server
	go func() {
		log.Info("Server listening",
			"address", fmt.Sprintf("http://%s:%s", cfg.Server.Host, cfg.Server.Port),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start", "error", err)
			os.Exit(1)
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

	// Initialize repository
	userRepo, cleanup, err := initializeRepository(cfg, log)
	if err != nil {
		return nil, nil, err
	}

	// Initialize services
	userService, cacheCleanup := initializeServices(cfg, userRepo, log)

	// Combine cleanup functions
	combinedCleanup := func() {
		cacheCleanup()
		cleanup()
	}

	// Initialize handlers and middleware
	userHandler := handler.NewUserHandler(userService)
	jwtMiddleware := middleware.NewJWTMiddleware(service.NewJWTTokenService(cfg))

	// Setup routes and server
	router := routes.NewRouter(userHandler, jwtMiddleware, baseLogger)
	httpRouter := router.SetupRoutes()

	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      httpRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return server, combinedCleanup, nil
}

// initializeRepository sets up the data repository based on configuration
func initializeRepository(cfg *config.Config, log *logger.Logger) (domain.UserRepository, func(), error) {
	repositoryType := os.Getenv("REPOSITORY_TYPE")

	if repositoryType == "memory" || repositoryType == "" {
		log.Info("Using in-memory repository")
		return repository.NewMemoryUserRepository(), func() {}, nil
	}

	if repositoryType == "mongodb" {
		log.Info("Using MongoDB repository")

		mongoClient, err := repository.NewMongoClient(cfg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
		}

		userRepo := repository.NewMongoUserRepository(mongoClient, cfg)

		cleanup := func() {
			log.Info("Disconnecting from MongoDB")
			ctx, cancel := context.WithTimeout(context.Background(), MongoDisconnectTimeout)
			defer cancel()
			if err := mongoClient.Disconnect(ctx); err != nil {
				log.Error("Error disconnecting from MongoDB", "error", err)
			} else {
				log.Info("Disconnected from MongoDB")
			}
		}

		return userRepo, cleanup, nil
	}

	return nil, nil, fmt.Errorf("unsupported repository type: %s", repositoryType)
}

// initializeServices sets up the business logic services with optional caching
func initializeServices(cfg *config.Config, userRepo domain.UserRepository, log *logger.Logger) (domain.UserService, func()) {
	tokenService := service.NewJWTTokenService(cfg)
	baseUserService := service.NewUserService(userRepo, tokenService)

	cacheType := os.Getenv("CACHE_TYPE")
	if cacheType != "redis" {
		log.Info("Cache disabled or not configured")
		return baseUserService, func() {}
	}

	log.Info("Initializing Redis cache")
	cacheService, err := cache.NewRedisCache(cfg)
	if err != nil {
		log.Warn("Failed to initialize Redis cache, using service without cache", "error", err)
		return baseUserService, func() {}
	}

	log.Info("Redis cache initialized successfully")
	userService := service.NewCachedUserService(baseUserService, cacheService, cfg.Cache.Redis.TTL)

	cleanup := func() {
		log.Info("Closing cache connection")
		if err := cacheService.Close(); err != nil {
			log.Error("Error closing cache connection", "error", err)
		} else {
			log.Info("Cache connection closed")
		}
	}

	return userService, cleanup
}
