package bootstrap

import (
	"context"
	"log"
	"time"

	"go_short/conf"
	// Identity Domain Imports
	// (如果需要在 bootstrap 中引用)

	identityservice "go_short/domain/identity/service"
	urlshortenerservice "go_short/domain/urlshortener/service"

	// Infrastructure Imports
	"go_short/infra/database"
	gormpersistence "go_short/infra/persistence/gorm"
	redispersistence "go_short/infra/persistence/redis"

	// API Imports
	"go_short/internal/api"
	"go_short/internal/api/handler"

	// Application Imports
	identityapp "go_short/internal/application/identity"
	urlshortenerapp "go_short/internal/application/urlshortener"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Dependencies 包含應用程式啟動所需的所有依賴項
type Dependencies struct {
	Config      *conf.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	GinEngine   *gin.Engine
	URLApp      *urlshortenerapp.App // URL Shortener Application instance
	IdentityApp *identityapp.App     // Identity Application instance
	UserHandler *handler.UserHandler // User Handler instance
	URLHandler  *handler.URLHandler  // URL Handler instance (保持現有)
}

// InitDependencies 初始化應用程式的所有依賴項
func InitDependencies() (*Dependencies, error) {
	log.Println("Initializing dependencies...")

	// 1. 初始化配置
	config := conf.Conf()
	log.Println("Configuration loaded.")

	// 2. 初始化資料庫連接
	db, err := database.InitDB(config)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return nil, err
	}
	log.Println("Database connection initialized.")

	// 3. 初始化 Redis 客戶端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := redisClient.Ping(pingCtx).Result(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		log.Println("Redis connection verified.")
	}

	// --- 依賴注入 ---
	log.Println("Setting up dependency injection...")

	// --- URL Shortener Domain Dependencies ---
	urlRepo := gormpersistence.NewGormURLRepository(db)
	cacheRepo := redispersistence.NewRedisCacheRepository(redisClient)
	urlDomainService := urlshortenerservice.NewURLService(urlRepo, cacheRepo, 24*time.Hour) // Assuming this constructor exists
	urlApp := urlshortenerapp.NewApp(urlDomainService)                                      // Assume NewApp takes the service now
	urlHandler := handler.NewURLHandler(urlDomainService)                                   // Handler depends on Domain Service
	log.Println("URL Shortener dependencies initialized.")

	// --- Identity Domain Dependencies ---
	userRepo := gormpersistence.NewGormUserRepository(db)
	identityDomainService := identityservice.NewIdentityService(userRepo)      // Create Identity Domain Service
	identityApplication := identityapp.NewApp(userRepo, identityDomainService) // Create Identity Application Service
	// (如果 identityApp 需要 identityDomainService, 則注入: identityapp.NewApp(userRepo, identityDomainService))
	userHandler := handler.NewUserHandler(identityApplication) // Create User Handler, inject App Service
	log.Println("Identity dependencies initialized.")

	// --- API Router Setup ---
	ginEngine := gin.Default()
	// 傳遞所有需要的 Handlers 給 Router
	apiRouter := api.NewRouter(ginEngine, urlHandler, userHandler, config) // 修改 NewRouter 以接收 UserHandler
	apiRouter.SetupRoutes()                                                // SetupRoutes 內部應分別設定 URL 和 User 路由
	log.Println("API Router initialized and routes set up.")
	// --- 依賴注入結束 ---

	deps := &Dependencies{
		Config:      config,
		DB:          db,
		RedisClient: redisClient,
		GinEngine:   ginEngine,
		URLApp:      urlApp,
		IdentityApp: identityApplication, // 保存 Identity App 實例
		UserHandler: userHandler,         // 保存 User Handler 實例
		URLHandler:  urlHandler,
	}

	log.Println("Dependencies initialized successfully.")
	return deps, nil
}

// Close gracefully closes the dependencies
func (d *Dependencies) Close() {
	log.Println("Closing resources...")
	if d.RedisClient != nil {
		if err := d.RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		} else {
			log.Println("Redis connection closed.")
		}
	}
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			} else {
				log.Println("Database connection closed.")
			}
		}
	}
	log.Println("Resources closed.")
}
