package services

import (
	"context"
	"fmt"
	"go_short/conf"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedisClient initializes the Redis client
func InitRedisClient() {
	config := conf.Conf()
	
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	// Test the connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		// We don't want to fail the application if Redis is not available
		// Just log the error and continue
	} else {
		log.Println("Connected to Redis successfully")
	}
}

// GetCachedURL retrieves a URL from the cache
func GetCachedURL(shortURL string) (string, bool) {
	if RedisClient == nil {
		return "", false
	}

	originalURL, err := RedisClient.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		// Key does not exist
		return "", false
	} else if err != nil {
		// Error occurred
		log.Printf("Redis error: %v", err)
		return "", false
	}

	return originalURL, true
}

// CacheURL stores a URL in the cache with expiration
func CacheURL(shortURL, originalURL string, expiration time.Duration) {
	if RedisClient == nil {
		return
	}

	err := RedisClient.Set(ctx, shortURL, originalURL, expiration).Err()
	if err != nil {
		log.Printf("Failed to cache URL: %v", err)
	}
}

// InvalidateCache removes a URL from the cache
func InvalidateCache(shortURL string) {
	if RedisClient == nil {
		return
	}

	err := RedisClient.Del(ctx, shortURL).Err()
	if err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}
}
