package conf

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	// URL Shortener
	ShortenerAlgorithm string
	// Server
	ServerPort string
}

var config *Config
var loadConfigOnce sync.Once

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using environment variables")
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0 // Default Redis DB
	}

	config = &Config{
		// Database
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		// Redis
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		// URL Shortener
		ShortenerAlgorithm: os.Getenv("SHORTENER_ALGORITHM"),
	}

	// Set default algorithm if not specified
	if config.ShortenerAlgorithm == "" {
		config.ShortenerAlgorithm = "base62"
	}
}

func Conf() *Config {
	loadConfigOnce.Do(loadConfig)
	return config
}
