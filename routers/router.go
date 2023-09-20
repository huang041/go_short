package router

import(
	"github.com/gin-gonic/gin"
	. "go_short/controllers"
)

func InitRouter() *gin.Engine{
	router := gin.Default()
	router.GET("/ping", Ping)

	urlMappingController := UrlMappingController{}
	router.GET("/url_mapping", urlMappingController.GetUrlMapping)
	router.POST("/url_mapping", urlMappingController.SaveUrlMapping)
	return router
}