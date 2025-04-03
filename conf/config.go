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

type Config struct {
	// Database
	DB Database
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

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	config = Config{
		DB: Database{
			Host: os.Getenv("DB_HOST"),
			Port: port,
			User: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DB_name: os.Getenv("DB_NAME"),
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