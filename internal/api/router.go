package api

import (
	"go_short/conf"
	"go_short/internal/api/handler"
	"go_short/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// Router 負責集中管理所有 API 路由
type Router struct {
	engine      *gin.Engine
	urlHandler  *handler.URLHandler
	userHandler *handler.UserHandler
	config      *conf.Config
	jwtSecret   []byte // 新增：保存 JWT 密鑰
}

func NewRouter(engine *gin.Engine, urlHandler *handler.URLHandler, userHandler *handler.UserHandler, config *conf.Config, jwtSecret []byte) *Router {
	return &Router{
		engine:      engine,
		urlHandler:  urlHandler,
		userHandler: userHandler,
		config:      config,
		jwtSecret:   jwtSecret, // 保存密鑰
	}
}

// SetupRoutes 設定所有的路由
func (r *Router) SetupRoutes() {
	r.setupMiddlewares()
	r.setupPublicRoutes()        // 公開訪問的路由
	r.setupAuthenticatedRoutes() // 需要認證的路由
}

// setupMiddlewares 設定全域中間件
func (r *Router) setupMiddlewares() {
	// 添加一些基礎中間件，例如 Logger, Recovery
	r.engine.Use(gin.Logger())
	r.engine.Use(gin.Recovery())

	// 可以添加 CORS 中間件等
	// ...

	// 將算法設置放入中間件，確保所有路由可用
	r.engine.Use(func(c *gin.Context) {
		c.Set("algorithm", r.config.ShortenerAlgorithm)
		c.Next()
	})
}

// setupPublicRoutes 設定無需認證即可訪問的路由
func (r *Router) setupPublicRoutes() {
	// 健康檢查
	r.engine.GET("/ping", r.urlHandler.HealthCheck)

	// 用戶認證相關 (註冊和登入本身不需要先登入)
	authGroup := r.engine.Group("/auth")
	{
		authGroup.POST("/register", r.userHandler.Register)
		authGroup.POST("/login", r.userHandler.Login)
	}

	// 短連結重定向 (通常是公開的)
	r.engine.GET("/:shortURL", r.urlHandler.RedirectToOriginalURL)

	// 創建短連結 (應用可選認證中間件，傳入密鑰)
	r.engine.POST("/url_mapping", middleware.OptionalAuthMiddleware(r.jwtSecret), r.urlHandler.CreateShortURL) // <--- 傳遞密鑰
}

// setupAuthenticatedRoutes 設定需要 JWT 認證才能訪問的路由
func (r *Router) setupAuthenticatedRoutes() {
	// 創建一個新的路由組，並應用 AuthMiddleware
	authenticated := r.engine.Group("/")
	authenticated.Use(middleware.AuthMiddleware(r.jwtSecret)) // <--- 傳遞密鑰
	{
		// URL 映射管理 (現在需要登入才能創建和查看列表)
		authenticated.GET("/url_mapping", r.urlHandler.GetAllURLMappings)
	}
}
