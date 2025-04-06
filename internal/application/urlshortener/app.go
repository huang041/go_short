package urlshortenerapp

import (
	"context"
	"log"
	"time"

	"go_short/domain/urlshortener/repository"
	"go_short/domain/urlshortener/service"
	gormpersistence "go_short/infra/persistence/gorm"
	redispersistence "go_short/infra/persistence/redis"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// App 是 URL 縮短服務的應用層
type App struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	URLService  service.URLShortenerService
	URLRepo     repository.URLRepository
	CacheRepo   repository.CacheRepository
}

// NewApp 創建一個新的 URL 縮短服務應用
func NewApp(db *gorm.DB, redisClient *redis.Client) *App {
	// 創建儲存庫
	urlRepo := gormpersistence.NewGormURLRepository(db)
	cacheRepo := redispersistence.NewRedisCacheRepository(redisClient)

	// 創建服務
	urlService := service.NewURLService(urlRepo, cacheRepo, 24*time.Hour)

	return &App{
		DB:          db,
		RedisClient: redisClient,
		URLService:  urlService,
		URLRepo:     urlRepo,
		CacheRepo:   cacheRepo,
	}
}

// GetURLService 返回 URL 服務，供外部 Handler 使用
func (app *App) GetURLService() service.URLShortenerService {
	return app.URLService
}

// InitDatabase 檢查數據庫連接（不再執行遷移）
func (app *App) InitDatabase() error {
	// 可以添加一個簡單的 ping 檢查，確保 DB 連接正常
	sqlDB, err := app.DB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return err
	}
	log.Println("Database connection verified.")
	return nil
}

// StartCleanupTask 啟動定期清理過期 URL 的任務
func (app *App) StartCleanupTask(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	log.Println("Starting background cleanup task...")
	go func() {
		defer log.Println("Background cleanup task stopped.")
		for {
			select {
			case <-ticker.C:
				log.Println("Running expired URLs cleanup...")
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
