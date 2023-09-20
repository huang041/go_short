package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"go_short/models"
)

type UrlMappingController struct{}

func ( umc * UrlMappingController ) GetUrlMapping(c *gin.Context) {
	var urlMappings []models.UrlMapping
    if err := models.DB.Find(&urlMappings).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "msg":  "Error while fetching data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": 200,
        "msg":  "success",
        "data": urlMappings,
    })
}

func ( umc * UrlMappingController ) SaveUrlMapping(c *gin.Context) {
	originURL := c.Query("url")
	models.DB.Create(&models.UrlMapping{Rename_url: "12.st.com", Origin_url: originURL})
}