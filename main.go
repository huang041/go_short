package main

import (
	"go_short/models"
	"go_short/routers"
	"go_short/services"
	"log"
)

func main() {
	models.InitDatabase()
	services.InitRedisClient()
	log.Println(models.DB)
	router := routers.InitRouter()
	router.Run("0.0.0.0:8080")
}
