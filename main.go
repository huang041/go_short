package main

import (
	"go_short/models"
	router "go_short/routers"
	"log"
)

func main() {
	models.InitDatabase()
	log.Println(models.DB)
	router := router.InitRouter()
	router.Run("0.0.0.0:8080")
}
