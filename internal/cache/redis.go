// Package cache provides caching functionality for the demo-go application
// using Redis as the backend cache store with configurable TTL and operations.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"demo-go/internal/config"
	"demo-go/internal/domain"
	"demo-go/internal/logger"

	"github.com/go-redis/redis/v8"
)

// Service defines the interface for cache operations
type Service interface {
	// User-specific cache operations
	GetUser(ctx context.Context, userID string) (*domain.UserResponse, error)
	SetUser(ctx context.Context, userID string, user *domain.UserResponse, ttl time.Duration) error
	DeleteUser(ctx context.Context, userID string) error

	// Generic cache operations
	Get(ctx context.Context, key string, result interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Batch operations
	DeleteByPattern(ctx context.Context, pattern string) error

	// Health check
	Ping(ctx context.Context) error

	// Close connection
	Close() error
}

// redisCache implements Service using Redis
type redisCache struct {
	client *redis.Client
	logger *logger.Logger
	config *config.RedisConfig
}

// NewRedisCache creates a new Redis cache service
func NewRedisCache(cfg *config.Config) (Service, error) {
	log := logger.GetGlobal().ForComponent("redis-cache")

	log.Info("Initializing Redis cache",
		"address", cfg.Cache.Redis.Address,
		"db", cfg.Cache.Redis.DB,
		"pool_size", cfg.Cache.Redis.PoolSize,
	)

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Cache.Redis.Address,
		Password:     cfg.Cache.Redis.Password,
		DB:           cfg.Cache.Redis.DB,
		MaxRetries:   cfg.Cache.Redis.MaxRetries,
		PoolSize:     cfg.Cache.Redis.PoolSize,
		MinIdleConns: cfg.Cache.Redis.MinIdleConns,
		DialTimeout:  cfg.Cache.Redis.DialTimeout,
		ReadTimeout:  cfg.Cache.Redis.ReadTimeout,
		WriteTimeout: cfg.Cache.Redis.WriteTimeout,
		IdleTimeout:  cfg.Cache.Redis.IdleTimeout,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Error("Failed to connect to Redis", "error", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Successfully connected to Redis")

	return &redisCache{
		client: client,
		logger: log,
		config: &cfg.Cache.Redis,
	}, nil
}

// GetUser retrieves a user from cache
func (c *redisCache) GetUser(ctx context.Context, userID string) (*domain.UserResponse, error) {
	key := c.userCacheKey(userID)
	log := c.logger.WithField("user_id", userID).WithField("cache_key", key)

	log.Debug("Getting user from cache")

	var user domain.UserResponse
	err := c.Get(ctx, key, &user)
	if err != nil {
		if err == redis.Nil {
			log.Debug("User cache miss")
			return nil, domain.ErrUserNotFound
		}
		log.Error("Failed to get user from cache", "error", err)
		return nil, err
	}

	log.Debug("User cache hit")
	return &user, nil
}

// SetUser stores a user in cache
func (c *redisCache) SetUser(ctx context.Context, userID string, user *domain.UserResponse, ttl time.Duration) error {
	key := c.userCacheKey(userID)
	log := c.logger.WithField("user_id", userID).WithField("cache_key", key).WithField("ttl", ttl)

	log.Debug("Setting user in cache")

	err := c.Set(ctx, key, user, ttl)
	if err != nil {
		log.Error("Failed to set user in cache", "error", err)
		return err
	}

	log.Debug("User cached successfully")
	return nil
}

// DeleteUser removes a user from cache
func (c *redisCache) DeleteUser(ctx context.Context, userID string) error {
	key := c.userCacheKey(userID)
	log := c.logger.WithField("user_id", userID).WithField("cache_key", key)

	log.Debug("Deleting user from cache")

	err := c.Delete(ctx, key)
	if err != nil {
		log.Error("Failed to delete user from cache", "error", err)
		return err
	}

	log.Debug("User deleted from cache")
	return nil
}

// Get retrieves a value from cache and unmarshals it into result
func (c *redisCache) Get(ctx context.Context, key string, result interface{}) error {
	log := c.logger.WithField("cache_key", key)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Debug("Cache miss")
			return redis.Nil
		}
		log.Error("Redis GET failed", "error", err)
		return err
	}

	if err := json.Unmarshal([]byte(val), result); err != nil {
		log.Error("Failed to unmarshal cached value", "error", err)
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	log.Debug("Cache hit")
	return nil
}

// Set stores a value in cache with TTL
func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	log := c.logger.WithField("cache_key", key).WithField("ttl", ttl)

	// Use default TTL if not specified
	if ttl == 0 {
		ttl = c.config.TTL
	}

	data, err := json.Marshal(value)
	if err != nil {
		log.Error("Failed to marshal value for caching", "error", err)
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		log.Error("Redis SET failed", "error", err)
		return err
	}

	log.Debug("Value cached successfully")
	return nil
}

// Delete removes a key from cache
func (c *redisCache) Delete(ctx context.Context, key string) error {
	log := c.logger.WithField("cache_key", key)

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		log.Error("Redis DELETE failed", "error", err)
		return err
	}

	log.Debug("Key deleted from cache")
	return nil
}

// Exists checks if a key exists in cache
func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	log := c.logger.WithField("cache_key", key)

	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		log.Error("Redis EXISTS failed", "error", err)
		return false, err
	}

	exists := count > 0
	log.Debug("Key existence check", "exists", exists)
	return exists, nil
}

// DeleteByPattern deletes all keys matching a pattern
func (c *redisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	log := c.logger.WithField("pattern", pattern)

	log.Debug("Deleting keys by pattern")

	// Get all keys matching the pattern
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		log.Error("Failed to get keys by pattern", "error", err)
		return err
	}

	if len(keys) == 0 {
		log.Debug("No keys found matching pattern")
		return nil
	}

	// Delete all matching keys
	pipe := c.client.Pipeline()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error("Failed to delete keys by pattern", "error", err, "key_count", len(keys))
		return err
	}

	log.Info("Deleted keys by pattern", "key_count", len(keys))
	return nil
}

// Ping checks if Redis is reachable
func (c *redisCache) Ping(ctx context.Context) error {
	log := c.logger

	err := c.client.Ping(ctx).Err()
	if err != nil {
		log.Error("Redis ping failed", "error", err)
		return err
	}

	log.Debug("Redis ping successful")
	return nil
}

// Close closes the Redis connection
func (c *redisCache) Close() error {
	log := c.logger

	log.Info("Closing Redis connection")

	err := c.client.Close()
	if err != nil {
		log.Error("Failed to close Redis connection", "error", err)
		return err
	}

	log.Info("Redis connection closed")
	return nil
}

// userCacheKey generates a cache key for user data
func (c *redisCache) userCacheKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

// InvalidateUserCache invalidates all cache entries related to a user
func (c *redisCache) InvalidateUserCache(ctx context.Context, userID string) error {
	log := c.logger.WithField("user_id", userID)

	log.Debug("Invalidating user cache")

	// Delete user-specific cache entries
	userKey := c.userCacheKey(userID)
	if err := c.Delete(ctx, userKey); err != nil {
		log.Error("Failed to invalidate user cache", "error", err)
		return err
	}

	// Delete user list cache entries (if any)
	listPattern := "users:list:*"
	if err := c.DeleteByPattern(ctx, listPattern); err != nil {
		log.Warn("Failed to invalidate user list cache", "error", err)
		// Don't return error for list cache invalidation failures
	}

	log.Debug("User cache invalidated successfully")
	return nil
}

// Stats represents cache statistics
type Stats struct {
	Hits        int64   `json:"hits"`
	Misses      int64   `json:"misses"`
	Keys        int64   `json:"keys"`
	Memory      int64   `json:"memory_bytes"`
	Connections int     `json:"connections"`
	HitRate     float64 `json:"hit_rate"`
}

// GetStats returns cache statistics
func (c *redisCache) GetStats(ctx context.Context) (*Stats, error) {
	log := c.logger

	log.Debug("Getting cache statistics")

	_, err := c.client.Info(ctx, "stats", "memory", "clients").Result()
	if err != nil {
		log.Error("Failed to get Redis info", "error", err)
		return nil, err
	}

	// Parse Redis INFO output
	stats := &Stats{}

	// Get database size
	dbSize, err := c.client.DBSize(ctx).Result()
	if err != nil {
		log.Warn("Failed to get database size", "error", err)
	} else {
		stats.Keys = dbSize
	}

	// Note: Parsing Redis INFO for detailed stats would require more complex parsing
	// For now, we'll return basic stats
	stats.Connections = c.config.PoolSize

	log.Debug("Cache statistics retrieved", "keys", stats.Keys)
	return stats, nil
}
