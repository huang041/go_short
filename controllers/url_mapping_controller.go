package controllers

import (
	"go_short/models"
	"go_short/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UrlMappingController struct{}

func (umc *UrlMappingController) GetUrlMapping(c *gin.Context) {
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

func (umc *UrlMappingController) SaveUrlMapping(c *gin.Context) {
	// Create a struct to bind the request body
	var request struct {
		URL string `json:"url" binding:"required"`
	}

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format or missing URL"})
		return
	}

	urlMapping := models.UrlMapping{Rename_url: nil, Origin_url: request.URL}

	if err := models.DB.Create(&urlMapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id := int(urlMapping.ID)
	urlMapping.Rename_url = services.DecimalToBase62(id)

	if err := models.DB.Save(&urlMapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(&urlMapping)
	c.JSON(200, gin.H{
		"short_url": urlMapping.Rename_url,
	})
}
