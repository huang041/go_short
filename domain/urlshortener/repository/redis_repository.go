package repository

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCacheRepository 是 CacheRepository 的 Redis 實現
type RedisCacheRepository struct {
	client *redis.Client
}

// NewRedisCacheRepository 創建一個新的 Redis 緩存儲存庫
func NewRedisCacheRepository(client *redis.Client) *RedisCacheRepository {
	return &RedisCacheRepository{
		client: client,
	}
}

// Get 從緩存中獲取 URL 映射
func (r *RedisCacheRepository) Get(ctx context.Context, shortURL string) (string, bool) {
	if r.client == nil {
		return "", false
	}

	originalURL, err := r.client.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		// Key 不存在
		return "", false
	} else if err != nil {
		// 發生錯誤
		log.Printf("Redis error: %v", err)
		return "", false
	}

	return originalURL, true
}

// Set 將 URL 映射保存到緩存
func (r *RedisCacheRepository) Set(ctx context.Context, shortURL string, originalURL string, expiration time.Duration) error {
	if r.client == nil {
		return nil
	}

	err := r.client.Set(ctx, shortURL, originalURL, expiration).Err()
	if err != nil {
		log.Printf("Failed to cache URL: %v", err)
		return err
	}
	return nil
}

// Delete 從緩存中刪除 URL 映射
func (r *RedisCacheRepository) Delete(ctx context.Context, shortURL string) error {
	if r.client == nil {
		return nil
	}

	err := r.client.Del(ctx, shortURL).Err()
	if err != nil {
		log.Printf("Failed to delete URL from cache: %v", err)
		return err
	}
	return nil
}
