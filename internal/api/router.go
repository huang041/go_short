package api

import (
	"go_short/conf"
	"go_short/domain/urlshortener/handler"

	"github.com/gin-gonic/gin"
)

// Router 負責集中管理所有 API 路由
type Router struct {
	engine     *gin.Engine
	urlHandler *handler.URLHandler
	config     *conf.Config
}

// NewRouter 建立一個新的路由管理器
func NewRouter(engine *gin.Engine, urlHandler *handler.URLHandler, config *conf.Config) *Router {
	return &Router{
		engine:     engine,
		urlHandler: urlHandler,
		config:     config,
	}
}

// SetupRoutes 設定所有的路由
func (r *Router) SetupRoutes() {
	r.setupMiddlewares()
	r.setupHealthCheckRoutes()
	r.setupURLShortenerRoutes()

	// 在未來可以增加更多其他領域的路由設定
	// r.setupUserRoutes()
	// r.setupAnalyticsRoutes()
	// 等等...
}

// setupMiddlewares 設定全域中間件
func (r *Router) setupMiddlewares() {
	// 設置中間件，將配置傳遞給處理器
	r.engine.Use(func(c *gin.Context) {
		c.Set("algorithm", r.config.ShortenerAlgorithm)
		c.Next()
	})

	// 可添加其他全域中間件如 CORS、認證、限流等
}

// setupHealthCheckRoutes 設定健康檢查路由
func (r *Router) setupHealthCheckRoutes() {
	r.engine.GET("/ping", r.urlHandler.HealthCheck)
}

// setupURLShortenerRoutes 設定短連結相關路由
func (r *Router) setupURLShortenerRoutes() {
	// URL 映射 API
	r.engine.GET("/url_mapping", r.urlHandler.GetAllURLMappings)
	r.engine.POST("/url_mapping", r.urlHandler.CreateShortURL)

	// 重定向 API
	r.engine.GET("/:shortURL", r.urlHandler.RedirectToOriginalURL)
}
