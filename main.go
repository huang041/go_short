package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go_short/conf"
	"go_short/domain/urlshortener/app"
	"go_short/infra/database"
	"go_short/internal/api"
	"go_short/internal/api/handler"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 初始化配置
	config := conf.Conf()

	// 初始化數據庫連接
	db, err := database.InitDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化 Redis 客戶端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// 創建應用程序
	ctx := context.Background()
	urlApp := app.NewURLShortenerApp(db, redisClient)
	urlService := urlApp.GetURLService()

	// 創建 Handler，注入 Service
	urlHandler := handler.NewURLHandler(urlService)

	// 初始化數據庫結構
	if err := urlApp.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 啟動定期清理過期 URL 的任務
	urlApp.StartCleanupTask(ctx)

	// 設置 Gin 路由
	server := gin.Default()

	// 使用集中式的路由管理
	apiRouter := api.NewRouter(server, urlHandler, config)
	apiRouter.SetupRoutes()

	// 啟動 HTTP 服務器
	go func() {
		if err := server.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 優雅關閉
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 關閉 Redis 連接
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("Server exiting")
}
