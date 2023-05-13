package router

import(
	"github.com/gin-gonic/gin"
	. "go_short/app/apis"
)

func InitRouter() *gin.Engine{
	router := gin.Default()
	router.GET("/ping", Ping)
	return router
}