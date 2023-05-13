package main

import (
	"go_short/app/router"
)



func main() {
	router := router.InitRouter()
	router.Run(":8080")
}