package app

import (
	"context"
	"log"
	"time"

	"go_short/domain/urlshortener/entity"
	"go_short/domain/urlshortener/handler"
	"go_short/domain/urlshortener/repository"
	"go_short/domain/urlshortener/service"
	gormpersistence "go_short/infra/persistence/gorm"
	redispersistence "go_short/infra/persistence/redis"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// URLShortenerApp 是 URL 縮短服務的應用層
type URLShortenerApp struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	URLHandler  *handler.URLHandler
	URLService  service.URLShortenerService
	URLRepo     repository.URLRepository
	CacheRepo   repository.CacheRepository
}

// NewURLShortenerApp 創建一個新的 URL 縮短服務應用
func NewURLShortenerApp(db *gorm.DB, redisClient *redis.Client) *URLShortenerApp {
	// 創建儲存庫
	urlRepo := gormpersistence.NewGormURLRepository(db)
	cacheRepo := redispersistence.NewRedisCacheRepository(redisClient)

	// 創建服務
	urlService := service.NewURLService(urlRepo, cacheRepo, 24*time.Hour)

	// 創建處理器
	urlHandler := handler.NewURLHandler(urlService)

	return &URLShortenerApp{
		DB:          db,
		RedisClient: redisClient,
		URLHandler:  urlHandler,
		URLService:  urlService,
		URLRepo:     urlRepo,
		CacheRepo:   cacheRepo,
	}
}

// GetURLHandler 返回 URL 處理器，供外部路由系統使用
func (app *URLShortenerApp) GetURLHandler() *handler.URLHandler {
	return app.URLHandler
}

// InitDatabase 初始化數據庫
func (app *URLShortenerApp) InitDatabase() error {
	// 自動遷移數據庫結構
	return app.DB.AutoMigrate(&entity.URLMapping{})
}

// StartCleanupTask 啟動定期清理過期 URL 的任務
func (app *URLShortenerApp) StartCleanupTask(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := app.URLService.CleanupExpiredURLs(ctx); err != nil {
					log.Printf("Failed to cleanup expired URLs: %v", err)
				} else {
					log.Println("Expired URLs cleanup completed successfully")
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
