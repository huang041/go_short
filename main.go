package main

import (
	"context"
	"log"
	"net/http" // 引入 net/http 以便使用 http.Server
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_short/internal/bootstrap" // 引入新的 bootstrap 包
)

func main() {
	// 初始化依賴項
	deps, err := bootstrap.InitDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer deps.Close()

	// 創建一個可用於取消背景任務的 context
	appCtx, cancelAppCtx := context.WithCancel(context.Background())
	defer cancelAppCtx()

	// 啟動定期清理過期 URL 的任務 (確保 URLApp 實例被正確傳遞)
	deps.URLApp.StartCleanupTask(appCtx) // 使用 Bootstrap 返回的 URLApp 實例

	// --- 配置和啟動 HTTP 伺服器 ---
	server := &http.Server{
		Addr:    ":8080",        // 應從 deps.Config 讀取
		Handler: deps.GinEngine, // 使用 bootstrap 返回的 gin Engine
	}

	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	// --- HTTP 伺服器啟動結束 ---

	// --- 優雅關閉 ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 給伺服器一點時間處理剩餘請求
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second) // 增加關閉超時
	defer cancelShutdown()

	// 觸發背景任務的取消
	cancelAppCtx()

	// 關閉 HTTP 伺服器
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
	// --- 優雅關閉結束 ---
}
