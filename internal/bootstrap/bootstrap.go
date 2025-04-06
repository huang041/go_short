package bootstrap

import (
	"context"
	"log"
	"time"

	"go_short/conf"
	"go_short/infra/database"
	"go_short/internal/api"
	"go_short/internal/api/handler"
	urlshortenerapp "go_short/internal/application/urlshortener"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	Config      *conf.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	GinEngine   *gin.Engine
	URLApp      *urlshortenerapp.App
}

func InitDependencies() (*Dependencies, error) {
	log.Println("Initializing dependencies...")

	config := conf.Conf()
	log.Println("Configuration loaded.")

	db, err := database.InitDB(config)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return nil, err
	}
	log.Println("Database connection initialized.")

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

	log.Println("Setting up dependency injection...")

	urlApp := urlshortenerapp.NewApp(db, redisClient)
	urlService := urlApp.GetURLService()
	log.Println("Application layer initialized.")

	urlHandler := handler.NewURLHandler(urlService)
	log.Println("API Handlers initialized.")

	ginEngine := gin.Default()
	apiRouter := api.NewRouter(ginEngine, urlHandler, config)
	apiRouter.SetupRoutes()
	log.Println("API Router initialized and routes set up.")

	deps := &Dependencies{
		Config:      config,
		DB:          db,
		RedisClient: redisClient,
		GinEngine:   ginEngine,
		URLApp:      urlApp,
	}

	log.Println("Dependencies initialized successfully.")
	return deps, nil
}

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
