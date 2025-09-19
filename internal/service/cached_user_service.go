package service

import (
	"context"
	"time"

	"demo-go/internal/cache"
	"demo-go/internal/domain"
	"demo-go/internal/logger"
)

// cachedUserService wraps a UserService with caching capabilities
type cachedUserService struct {
	userService domain.UserService
	cache       cache.CacheService
	logger      *logger.Logger
	cacheTTL    time.Duration
}

// NewCachedUserService creates a new cached user service wrapper
func NewCachedUserService(userService domain.UserService, cacheService cache.CacheService, cacheTTL time.Duration) domain.UserService {
	return &cachedUserService{
		userService: userService,
		cache:       cacheService,
		logger:      logger.GetGlobal().ForComponent("cached-user-service"),
		cacheTTL:    cacheTTL,
	}
}

// Register creates a new user account (no caching needed for write operations)
func (s *cachedUserService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	log := s.logger.ForService("user", "register").WithField("email", req.Email)
	
	log.Debug("Registering new user (bypassing cache)")
	
	user, err := s.userService.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// Cache the newly created user
	if cacheErr := s.cache.SetUser(ctx, user.ID, user, s.cacheTTL); cacheErr != nil {
		log.Warn("Failed to cache newly registered user", "user_id", user.ID, "error", cacheErr)
		// Don't fail the operation if caching fails
	} else {
		log.Debug("Cached newly registered user", "user_id", user.ID)
	}
	
	return user, nil
}

// Login authenticates a user and returns a JWT token (no caching needed for authentication)
func (s *cachedUserService) Login(ctx context.Context, req *domain.LoginRequest) (string, *domain.UserResponse, error) {
	log := s.logger.ForService("user", "login").WithField("email", req.Email)
	
	log.Debug("User login (bypassing cache for authentication)")
	
	token, user, err := s.userService.Login(ctx, req)
	if err != nil {
		return "", nil, err
	}
	
	// Cache the user data after successful login
	if cacheErr := s.cache.SetUser(ctx, user.ID, user, s.cacheTTL); cacheErr != nil {
		log.Warn("Failed to cache user after login", "user_id", user.ID, "error", cacheErr)
		// Don't fail the operation if caching fails
	} else {
		log.Debug("Cached user after login", "user_id", user.ID)
	}
	
	return token, user, nil
}

// GetProfile retrieves a user profile (cache-enabled)
func (s *cachedUserService) GetProfile(ctx context.Context, userID string) (*domain.UserResponse, error) {
	log := s.logger.ForService("user", "get-profile").WithField("user_id", userID)
	
	log.Debug("Getting user profile")
	
	// Try to get from cache first
	user, err := s.cache.GetUser(ctx, userID)
	if err == nil {
		log.Debug("User profile cache hit")
		return user, nil
	}
	
	// Cache miss or error - check if it's a real miss vs error
	if err != domain.ErrUserNotFound {
		log.Warn("Cache error when getting user profile", "error", err)
	} else {
		log.Debug("User profile cache miss")
	}
	
	// Get from underlying service
	user, err = s.userService.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	if cacheErr := s.cache.SetUser(ctx, userID, user, s.cacheTTL); cacheErr != nil {
		log.Warn("Failed to cache user profile", "user_id", userID, "error", cacheErr)
		// Don't fail the operation if caching fails
	} else {
		log.Debug("Cached user profile", "user_id", userID)
	}
	
	return user, nil
}

// UpdateProfile updates a user profile and invalidates cache
func (s *cachedUserService) UpdateProfile(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	log := s.logger.ForService("user", "update-profile").WithField("user_id", userID)
	
	log.Debug("Updating user profile")
	
	// Update in underlying service
	user, err := s.userService.UpdateProfile(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	
	// Invalidate cache for this user
	if cacheErr := s.cache.DeleteUser(ctx, userID); cacheErr != nil {
		log.Warn("Failed to invalidate user cache after update", "user_id", userID, "error", cacheErr)
	} else {
		log.Debug("Invalidated user cache after update", "user_id", userID)
	}
	
	// Cache the updated user
	if cacheErr := s.cache.SetUser(ctx, userID, user, s.cacheTTL); cacheErr != nil {
		log.Warn("Failed to cache updated user", "user_id", userID, "error", cacheErr)
	} else {
		log.Debug("Cached updated user", "user_id", userID)
	}
	
	return user, nil
}

// GetUsers retrieves a list of users (cache-enabled with list caching strategy)
func (s *cachedUserService) GetUsers(ctx context.Context, limit, offset int) ([]*domain.UserResponse, int64, error) {
	log := s.logger.ForService("user", "get-users").WithFields(map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	})
	
	log.Debug("Getting users list")
	
	// For list operations, we could implement more complex caching strategies
	// For now, we'll bypass cache for list operations and delegate to underlying service
	// This avoids complex cache invalidation scenarios for list data
	
	users, total, err := s.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	// Opportunistically cache individual users from the list
	go func() {
		// Use background context to avoid cancellation
		bgCtx := context.Background()
		for _, user := range users {
			if cacheErr := s.cache.SetUser(bgCtx, user.ID, user, s.cacheTTL); cacheErr != nil {
				log.Debug("Failed to cache user from list", "user_id", user.ID, "error", cacheErr)
			}
		}
		log.Debug("Opportunistically cached users from list", "count", len(users))
	}()
	
	return users, total, nil
}

// GetUserByID retrieves a user by ID (cache-enabled)
func (s *cachedUserService) GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	log := s.logger.ForService("user", "get-by-id").WithField("user_id", id)
	
	log.Debug("Getting user by ID")
	
	// Try to get from cache first
	user, err := s.cache.GetUser(ctx, id)
	if err == nil {
		log.Debug("User cache hit")
		return user, nil
	}
	
	// Cache miss or error - check if it's a real miss vs error
	if err != domain.ErrUserNotFound {
		log.Warn("Cache error when getting user by ID", "error", err)
	} else {
		log.Debug("User cache miss")
	}
	
	// Get from underlying service
	user, err = s.userService.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	if cacheErr := s.cache.SetUser(ctx, id, user, s.cacheTTL); cacheErr != nil {
		log.Warn("Failed to cache user", "user_id", id, "error", cacheErr)
		// Don't fail the operation if caching fails
	} else {
		log.Debug("Cached user", "user_id", id)
	}
	
	return user, nil
}

// DeleteUser deletes a user and invalidates cache
func (s *cachedUserService) DeleteUser(ctx context.Context, id string) error {
	log := s.logger.ForService("user", "delete").WithField("user_id", id)
	
	log.Debug("Deleting user")
	
	// Delete from underlying service
	err := s.userService.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	
	// Invalidate cache for this user
	if cacheErr := s.cache.DeleteUser(ctx, id); cacheErr != nil {
		log.Warn("Failed to invalidate user cache after deletion", "user_id", id, "error", cacheErr)
		// Don't fail the operation if cache invalidation fails
	} else {
		log.Debug("Invalidated user cache after deletion", "user_id", id)
	}
	
	return nil
}

// RefreshToken generates a new token for the user (cache-enabled for user lookup)
func (s *cachedUserService) RefreshToken(ctx context.Context, userID string) (string, error) {
	log := s.logger.ForService("user", "refresh-token").WithField("user_id", userID)
	
	log.Debug("Refreshing user token")
	
	// For token refresh, we need fresh user data from the database
	// to ensure the user is still active and valid
	// So we bypass cache for this operation
	token, err := s.userService.RefreshToken(ctx, userID)
	if err != nil {
		return "", err
	}
	
	log.Debug("Token refreshed successfully")
	return token, nil
}

// CacheHealthCheck checks the health of the cache service
func (s *cachedUserService) CacheHealthCheck(ctx context.Context) error {
	return s.cache.Ping(ctx)
}

// GetCacheStats returns cache statistics if the underlying cache supports it
func (s *cachedUserService) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	log := s.logger.WithField("operation", "get-cache-stats")
	
	// For now, return basic cache health info
	// More detailed stats implementation would require extending the cache interface
	if err := s.cache.Ping(ctx); err != nil {
		log.Error("Cache health check failed", "error", err)
		return map[string]interface{}{
			"healthy": false,
			"error":   err.Error(),
		}, err
	}
	
	return map[string]interface{}{
		"healthy": true,
		"message": "Cache is operational",
	}, nil
}

// InvalidateAllUserCache invalidates all user-related cache entries
func (s *cachedUserService) InvalidateAllUserCache(ctx context.Context) error {
	log := s.logger.WithField("operation", "invalidate-all-cache")
	
	log.Info("Invalidating all user cache")
	
	// Delete all user cache entries
	err := s.cache.DeleteByPattern(ctx, "user:*")
	if err != nil {
		log.Error("Failed to invalidate all user cache", "error", err)
		return err
	}
	
	log.Info("All user cache invalidated successfully")
	return nil
}
