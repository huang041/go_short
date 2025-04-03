package conf

import (
	"os"
	"log"
	"github.com/joho/godotenv"
	"strconv"
	"sync"
)

type Database struct {
	Host string
	Port int
	User string
	Password string
	DB_name string
}

type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type Config struct {
	// Database
	DB Database
	// Redis
	Redis Redis
	// URL Shortener
	ShortenerAlgorithm string
}

var config Config
var loadConfigOnce sync.Once

func loadConfig() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432 // Default PostgreSQL port
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		redisPort = 6379 // Default Redis port
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0 // Default Redis DB
	}

	config = Config{
		DB: Database{
			Host: os.Getenv("DB_HOST"),
			Port: dbPort,
			User: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DB_name: os.Getenv("DB_NAME"),
		},
		Redis: Redis{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     redisPort,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		ShortenerAlgorithm: os.Getenv("SHORTENER_ALGORITHM"),
	}

	// Set default algorithm if not specified
	if config.ShortenerAlgorithm == "" {
		config.ShortenerAlgorithm = "base62"
	}
}

func Conf() Config {
	loadConfigOnce.Do(loadConfig)
	return config
}