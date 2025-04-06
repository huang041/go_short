package repository

import (
	"context"
	"time"

	"go_short/domain/urlshortener/entity"
)

// URLRepository 定義了 URL 映射的儲存庫介面
type URLRepository interface {
	// FindByShortURL 根據短 URL 查找映射
	FindByShortURL(ctx context.Context, shortURL string) (*entity.URLMapping, error)
	
	// FindByOriginalURL 根據原始 URL 查找映射
	FindByOriginalURL(ctx context.Context, originalURL string) (*entity.URLMapping, error)
	
	// Save 保存 URL 映射
	Save(ctx context.Context, mapping *entity.URLMapping) error
	
	// Update 更新 URL 映射
	Update(ctx context.Context, mapping *entity.URLMapping) error
	
	// FindAll 獲取所有 URL 映射
	FindAll(ctx context.Context) ([]*entity.URLMapping, error)
	
	// DeleteExpired 刪除所有過期的 URL 映射
	DeleteExpired(ctx context.Context) error
}

// CacheRepository 定義了 URL 映射的緩存儲存庫介面
type CacheRepository interface {
	// Get 從緩存中獲取 URL 映射
	Get(ctx context.Context, shortURL string) (string, bool)
	
	// Set 將 URL 映射保存到緩存
	Set(ctx context.Context, shortURL string, originalURL string, expiration time.Duration) error
	
	// Delete 從緩存中刪除 URL 映射
	Delete(ctx context.Context, shortURL string) error
}
