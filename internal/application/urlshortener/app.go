package urlshortenerapp

import (
	"context"
	"log"
	"time"

	"go_short/domain/urlshortener/service"
)

// App 是 URL 縮短服務的應用層
type App struct {
	URLService service.URLShortenerService // 依賴 Domain Service Interface
}

// NewApp 創建應用服務實例，接收 Service Interface 作為依賴
func NewApp(urlService service.URLShortenerService /*, urlRepo repository.URLRepository, cacheRepo repository.CacheRepository*/) *App {
	return &App{
		URLService: urlService,
	}
}

// GetURLService 返回 URL 服務 (現在只是返回注入的 service)
func (app *App) GetURLService() service.URLShortenerService {
	return app.URLService
}

// InitDatabase 檢查數據庫連接 - 現在應由 Domain Service 或 Infra 處理
// 這個方法可能不再屬於 Application 層的職責
func (app *App) InitDatabase() error {
	// 如果 Service 需要 DB 連接，可以透過 Service 的方法檢查
	// 或者直接移除此方法，讓 bootstrap 或 infra 處理連接檢查
	log.Println("Database initialization/check responsibility moved.")
	return nil
}

// StartCleanupTask 啟動背景任務
func (app *App) StartCleanupTask(ctx context.Context) {
	// 這個邏輯可以保留在 App 層，因為它協調了 Service 的操作
	ticker := time.NewTicker(1 * time.Hour)
	log.Println("Starting background cleanup task...")
	go func() {
		defer log.Println("Background cleanup task stopped.")
		for {
			select {
			case <-ticker.C:
				log.Println("Running expired URLs cleanup...")
				// 呼叫注入的 Service
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
