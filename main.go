package main

import (
	"go_short/routers"
	"go_short/models"
	"log"
)



func main() {
	models.InitDatabase()
	log.Println(models.DB)
	router := router.InitRouter()
	router.Run("127.0.0.1:8080")
}