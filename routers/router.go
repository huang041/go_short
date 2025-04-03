package routers

import (
	. "go_short/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", Ping)

	urlMappingController := UrlMappingController{}
	router.GET("/url_mapping", urlMappingController.GetUrlMapping)
	router.POST("/url_mapping", urlMappingController.SaveUrlMapping)

	// Route for redirecting to original URL when visiting /{short_url}
	router.GET("/:shortURL", urlMappingController.RedirectToOriginalUrl)

	return router
}
