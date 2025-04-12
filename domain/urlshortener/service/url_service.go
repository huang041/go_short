package service

import (
	"context"
	"errors"
	"time"

	"go_short/domain/urlshortener/entity"
	"go_short/domain/urlshortener/repository"
)

// URLService 錯誤定義
var (
	ErrURLNotFound     = errors.New("URL not found")
	ErrURLExpired      = errors.New("URL has expired")
	ErrInvalidURL      = errors.New("invalid URL format")
	ErrDatabaseError   = errors.New("database operation failed")
	ErrCacheError      = errors.New("cache operation failed")
)

// URLShortenerService 定義了 URL 縮短服務的介面
type URLShortenerService interface {
	// CreateShortURL 創建一個新的短 URL
	CreateShortURL(ctx context.Context, originalURL string, algorithm string, expiresIn *time.Duration) (*entity.URLMapping, error)
	
	// GetOriginalURL 根據短 URL 獲取原始 URL
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	
	// GetAllURLMappings 獲取所有 URL 映射
	GetAllURLMappings(ctx context.Context) ([]*entity.URLMapping, error)
	
	// CleanupExpiredURLs 清理過期的 URL 映射
	CleanupExpiredURLs(ctx context.Context) error
}

// URLService 是 URLShortenerService 的實現
type URLService struct {
	urlRepo       repository.URLRepository
	cacheRepo     repository.CacheRepository
	cacheDuration time.Duration
}

// NewURLService 創建一個新的 URL 服務
func NewURLService(urlRepo repository.URLRepository, cacheRepo repository.CacheRepository, cacheDuration time.Duration) *URLService {
	return &URLService{
		urlRepo:       urlRepo,
		cacheRepo:     cacheRepo,
		cacheDuration: cacheDuration,
	}
}

// CreateShortURL 創建一個新的短 URL
func (s *URLService) CreateShortURL(ctx context.Context, originalURL string, algorithm string, expiresIn *time.Duration) (*entity.URLMapping, error) {
	// 檢查 URL 是否已存在
	existingMapping, err := s.urlRepo.FindByOriginalURL(ctx, originalURL)
	if err != nil {
		return nil, ErrDatabaseError
	}
	
	// 如果 URL 已存在，直接返回
	if existingMapping != nil {
		return existingMapping, nil
	}
	
	// 創建新的 URL 映射
	urlMapping := &entity.URLMapping{
		OriginalURL: originalURL,
		Algorithm:   algorithm,
	}
	
	// 設置過期時間（如果有）
	if expiresIn != nil {
		expiresAt := time.Now().Add(*expiresIn)
		urlMapping.ExpiresAt = &expiresAt
	}
	
	// 保存到數據庫以獲取 ID
	if err := s.urlRepo.Save(ctx, urlMapping); err != nil {
		return nil, ErrDatabaseError
	}
	
	// 根據算法生成短 URL
	id := int(urlMapping.ID)
	var shortener ShortenerStrategy
	
	switch algorithm {
	case "base64":
		shortener = &Base64Strategy{}
	case "md5":
		shortener = &MD5Strategy{}
	case "random":
		shortener = &RandomStrategy{}
	default: // base62 是默認值
		shortener = &Base62Strategy{}
	}
	
	// 生成短 URL
	urlMapping.ShortURL = shortener.Generate(originalURL, id)
	
	// 更新數據庫
	if err := s.urlRepo.Update(ctx, urlMapping); err != nil {
		return nil, ErrDatabaseError
	}
	
	// 緩存 URL 映射
	if urlMapping.ShortURL != nil {
		cacheExpiration := s.cacheDuration
		if urlMapping.ExpiresAt != nil {
			// 如果 URL 有過期時間，使用較短的緩存時間
			timeUntilExpiry := time.Until(*urlMapping.ExpiresAt)
			if timeUntilExpiry < cacheExpiration {
				cacheExpiration = timeUntilExpiry
			}
		}
		
		s.cacheRepo.Set(ctx, *urlMapping.ShortURL, urlMapping.OriginalURL, cacheExpiration)
	}
	
	return urlMapping, nil
}

// GetOriginalURL 根據短 URL 獲取原始 URL
func (s *URLService) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	// 先從緩存中查找
	if originalURL, found := s.cacheRepo.Get(ctx, shortURL); found {
		return originalURL, nil
	}
	
	// 如果緩存中沒有，從數據庫查找
	urlMapping, err := s.urlRepo.FindByShortURL(ctx, shortURL)
	if err != nil {
		return "", ErrDatabaseError
	}
	
	if urlMapping == nil {
		return "", ErrURLNotFound
	}
	
	// 檢查 URL 是否過期
	if urlMapping.IsExpired() {
		return "", ErrURLExpired
	}
	
	// 增加訪問計數
	urlMapping.IncrementVisits()
	if err := s.urlRepo.Update(ctx, urlMapping); err != nil {
		// 這裡我們只記錄錯誤，不阻止用戶訪問
		// 因為增加訪問計數不是關鍵操作
	}
	
	// 緩存結果
	cacheExpiration := s.cacheDuration
	if urlMapping.ExpiresAt != nil {
		// 如果 URL 有過期時間，使用較短的緩存時間
		timeUntilExpiry := time.Until(*urlMapping.ExpiresAt)
		if timeUntilExpiry < cacheExpiration {
			cacheExpiration = timeUntilExpiry
		}
	}
	
	s.cacheRepo.Set(ctx, shortURL, urlMapping.OriginalURL, cacheExpiration)
	
	return urlMapping.OriginalURL, nil
}

// GetAllURLMappings 獲取所有 URL 映射
func (s *URLService) GetAllURLMappings(ctx context.Context) ([]*entity.URLMapping, error) {
	return s.urlRepo.FindAll(ctx)
}

// CleanupExpiredURLs 清理過期的 URL 映射
func (s *URLService) CleanupExpiredURLs(ctx context.Context) error {
	return s.urlRepo.DeleteExpired(ctx)
}
