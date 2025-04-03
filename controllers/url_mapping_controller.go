package controllers

import (
	"go_short/conf"
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

// RedirectToOriginalUrl redirects to the original URL when a user visits /{short_url}
func (umc *UrlMappingController) RedirectToOriginalUrl(c *gin.Context) {
	shortURL := c.Param("shortURL")
	
	var urlMapping models.UrlMapping
	result := models.DB.Where("rename_url = ?", shortURL).First(&urlMapping)
	
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "Short URL not found",
		})
		return
	}
	
	c.Redirect(http.StatusFound, urlMapping.Origin_url)
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

	// Get the algorithm from configuration
	config := conf.Conf()
	algorithm := config.ShortenerAlgorithm

	// Create URL mapping with algorithm from config
	urlMapping := models.UrlMapping{
		Rename_url: nil, 
		Origin_url: request.URL,
		Algorithm:  algorithm,
	}

	// Save to database to get ID
	if err := models.DB.Create(&urlMapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get ID for shortening
	id := int(urlMapping.ID)
	
	// Create shortener with appropriate strategy based on algorithm from config
	var shortener *services.URLShortener
	
	switch algorithm {
	case "base64":
		shortener = services.NewURLShortener(&services.Base64Strategy{})
	case "md5":
		shortener = services.NewURLShortener(&services.MD5Strategy{})
	case "random":
		shortener = services.NewURLShortener(&services.RandomStrategy{})
	default: // base62 is default
		shortener = services.NewURLShortener(&services.Base62Strategy{})
	}
	
	// Generate short URL
	urlMapping.Rename_url = shortener.ShortenURL(request.URL, id)

	// Save updated mapping with short URL
	if err := models.DB.Save(&urlMapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	log.Println(&urlMapping)
	c.JSON(200, gin.H{
		"short_url": urlMapping.Rename_url,
		"algorithm": urlMapping.Algorithm,
	})
}
