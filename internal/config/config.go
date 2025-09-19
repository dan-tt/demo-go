package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	JWT      JWTConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	MongoDB MongoDBConfig
}

// MongoDBConfig holds MongoDB-specific configuration
type MongoDBConfig struct {
	URI         string
	Database    string
	Timeout     time.Duration
	MaxPoolSize int
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Redis RedisConfig
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Address      string
	Password     string
	DB           int
	MaxRetries   int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	TTL          time.Duration
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	SecretKey  string
	Expiration time.Duration
}

// Load creates and returns a new Config with values from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			MongoDB: MongoDBConfig{
				URI:         getEnv("MONGODB_URI", "mongodb://localhost:27017"),
				Database:    getEnv("MONGODB_DATABASE", "demo_clean"),
				Timeout:     getDurationEnv("MONGODB_TIMEOUT", 10*time.Second),
				MaxPoolSize: getIntEnv("MONGODB_MAX_POOL_SIZE", 100),
			},
		},
		Cache: CacheConfig{
			Redis: RedisConfig{
				Address:      getEnv("REDIS_ADDRESS", "localhost:6379"),
				Password:     getEnv("REDIS_PASSWORD", ""),
				DB:           getIntEnv("REDIS_DB", 0),
				MaxRetries:   getIntEnv("REDIS_MAX_RETRIES", 3),
				PoolSize:     getIntEnv("REDIS_POOL_SIZE", 10),
				MinIdleConns: getIntEnv("REDIS_MIN_IDLE_CONNS", 2),
				DialTimeout:  getDurationEnv("REDIS_DIAL_TIMEOUT", 5*time.Second),
				ReadTimeout:  getDurationEnv("REDIS_READ_TIMEOUT", 3*time.Second),
				WriteTimeout: getDurationEnv("REDIS_WRITE_TIMEOUT", 3*time.Second),
				IdleTimeout:  getDurationEnv("REDIS_IDLE_TIMEOUT", 5*time.Minute),
				TTL:          getDurationEnv("REDIS_TTL", 1*time.Hour),
			},
		},
		JWT: JWTConfig{
			SecretKey:  getEnv("JWT_SECRET_KEY", "default-secret-key-change-in-production"),
			Expiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets an environment variable as int or returns a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv gets an environment variable as duration or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}