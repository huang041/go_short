package bootstrap

import (
	"context"
	"log"
	"os"
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
	JWTSecret   []byte               // 新增：保存 JWT 密鑰
}

// InitDependencies 初始化應用程式的所有依賴項
func InitDependencies() (*Dependencies, error) {
	log.Println("Initializing dependencies...")

	// 1. 初始化配置
	config := conf.Conf()
	log.Println("Configuration loaded.")

	// --- 讀取 JWT 配置 (只讀一次) ---
	jwtSecretString := os.Getenv("JWT_SECRET")
	if jwtSecretString == "" {
		log.Println("CRITICAL WARNING: JWT_SECRET environment variable not set. Using default insecure key.")
		jwtSecretString = "a_very_insecure_default_secret_key_change_me" // 極不安全
	}
	jwtSecretBytes := []byte(jwtSecretString)

	jwtExpStr := os.Getenv("JWT_EXPIRATION_HOURS")
	jwtExpirationDuration, err := time.ParseDuration(jwtExpStr + "h")
	if err != nil || jwtExpirationDuration <= 0 {
		log.Printf("Warning: Invalid or missing JWT_EXPIRATION_HOURS. Using default 24 hours. Error: %v", err)
		jwtExpirationDuration = 24 * time.Hour
	}
	log.Printf("JWT Expiration set to: %v", jwtExpirationDuration)
	// --- JWT 配置讀取結束 ---

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
	urlDomainService := urlshortenerservice.NewURLService(urlRepo, cacheRepo, 24*time.Hour)
	urlApp := urlshortenerapp.NewApp(urlDomainService)
	urlHandler := handler.NewURLHandler(urlDomainService)
	log.Println("URL Shortener dependencies initialized.")

	// --- Identity Domain Dependencies ---
	userRepo := gormpersistence.NewGormUserRepository(db)
	identityDomainService := identityservice.NewIdentityService(userRepo)
	identityApplication := identityapp.NewApp(userRepo, identityDomainService, jwtSecretBytes, jwtExpirationDuration)
	userHandler := handler.NewUserHandler(identityApplication)
	log.Println("Identity dependencies initialized.")

	// --- API Router Setup ---
	ginEngine := gin.Default()
	// 傳遞所有需要的 Handlers 給 Router
	apiRouter := api.NewRouter(ginEngine, urlHandler, userHandler, config, jwtSecretBytes)
	apiRouter.SetupRoutes()
	log.Println("API Router initialized and routes set up.")
	// --- 依賴注入結束 ---

	deps := &Dependencies{
		Config:      config,
		DB:          db,
		RedisClient: redisClient,
		GinEngine:   ginEngine,
		URLApp:      urlApp,
		IdentityApp: identityApplication,
		UserHandler: userHandler,
		URLHandler:  urlHandler,
		JWTSecret:   jwtSecretBytes,
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
