package redispersistence

import (
	"context"
	"errors"
	"log"
	"time"

	"go_short/domain/urlshortener/repository"

	"github.com/redis/go-redis/v9"
)

// cacheRepository 是 CacheRepository 的 Redis 實現
type cacheRepository struct {
	client *redis.Client
}

// NewRedisCacheRepository 創建一個新的 Redis Cache 儲存庫實例
// 返回的是接口類型
func NewRedisCacheRepository(client *redis.Client) repository.CacheRepository {
	return &cacheRepository{
		client: client,
	}
}

// Get 從緩存中獲取 URL 映射
func (r *cacheRepository) Get(ctx context.Context, shortURL string) (string, bool) {
	if r.client == nil {
		return "", false
	}

	originalURL, err := r.client.Get(ctx, shortURL).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false // Key 不存在
		}
		// 可以考慮記錄其他錯誤
		return "", false
	}
	return originalURL, true
}

// Set 將 URL 映射保存到緩存
func (r *cacheRepository) Set(ctx context.Context, shortURL string, originalURL string, expiration time.Duration) error {
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
func (r *cacheRepository) Delete(ctx context.Context, shortURL string) error {
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
