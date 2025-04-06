package handler

import (
	"net/http"
	"time"

	"go_short/domain/urlshortener/service"

	"github.com/gin-gonic/gin"
)

// URLHandler 處理 URL 相關的 HTTP 請求
type URLHandler struct {
	urlService service.URLShortenerService
}

// NewURLHandler 創建一個新的 URL 處理器
func NewURLHandler(urlService service.URLShortenerService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// CreateShortURL 處理創建短 URL 的請求
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var request struct {
		URL       string `json:"url" binding:"required"`
		ExpiresIn *int   `json:"expires_in,omitempty"` // 過期時間（以小時為單位）
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format or missing URL",
		})
		return
	}

	// 從配置中獲取算法
	algorithm := c.GetString("algorithm")
	if algorithm == "" {
		algorithm = "base62" // 默認算法
	}

	// 設置過期時間（如果有）
	var expiresIn *time.Duration
	if request.ExpiresIn != nil {
		duration := time.Duration(*request.ExpiresIn) * time.Hour
		expiresIn = &duration
	}

	// 創建短 URL
	urlMapping, err := h.urlService.CreateShortURL(c.Request.Context(), request.URL, algorithm, expiresIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url":  urlMapping.ShortURL,
		"algorithm":  urlMapping.Algorithm,
		"expires_at": urlMapping.ExpiresAt,
	})
}

// GetAllURLMappings 處理獲取所有 URL 映射的請求
func (h *URLHandler) GetAllURLMappings(c *gin.Context) {
	mappings, err := h.urlService.GetAllURLMappings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Error while fetching data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": mappings,
	})
}

// RedirectToOriginalURL 處理重定向到原始 URL 的請求
func (h *URLHandler) RedirectToOriginalURL(c *gin.Context) {
	shortURL := c.Param("shortURL")

	originalURL, err := h.urlService.GetOriginalURL(c.Request.Context(), shortURL)
	if err != nil {
		switch err {
		case service.ErrURLNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"code": http.StatusNotFound,
				"msg":  "Short URL not found",
			})
		case service.ErrURLExpired:
			c.JSON(http.StatusGone, gin.H{
				"code": http.StatusGone,
				"msg":  "URL has expired",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Server error",
			})
		}
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

// HealthCheck 處理健康檢查請求
func (h *URLHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
