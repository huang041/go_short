package handler

import (
	"errors"
	"log"
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

func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var request struct {
		URL       string `json:"url" binding:"required,url"`
		ExpiresIn *int   `json:"expires_in,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// --- 嘗試獲取 UserID (如果請求來自已認證用戶) ---
	// 注意：因為端點是公開的，AuthMiddleware 不會執行，除非用戶自己提供了有效的 Token
	// 這裡的邏輯保持不變，如果 AuthMiddleware (或其他方式) 設置了 userID，就使用它
	userIDValue, exists := c.Get("userID")
	var userID *uint
	if exists {
		id, ok := userIDValue.(uint)
		if ok {
			userID = &id
			log.Printf("Creating short URL for authenticated user ID: %d", *userID)
		} else {
			// 這種情況理論上不應發生，如果 userID 存在但類型不對，可能是中間件或程式碼有問題
			log.Printf("Warning: userID found in context but not uint type: %T. Treating as anonymous.", userIDValue)
			// 安全起見，視為匿名
		}
	} else {
		log.Println("Creating short URL anonymously.")
		// userID 保持為 nil
	}
	// --- UserID 獲取結束 ---

	algorithm := c.GetString("algorithm") // 算法設置仍然來自全域中間件
	if algorithm == "" {
		algorithm = "base62"
	}

	var expiresIn *time.Duration
	if request.ExpiresIn != nil {
		if *request.ExpiresIn <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "expires_in must be positive"})
			return
		}
		duration := time.Duration(*request.ExpiresIn) * time.Hour
		expiresIn = &duration
	}

	// 將 userID (可能是 nil) 傳遞給 Service 層
	urlMapping, err := h.urlService.CreateShortURL(c.Request.Context(), request.URL, algorithm, userID, expiresIn)
	if err != nil {
		log.Printf("Error creating short URL: %v", err)
		// 可以根據 Service 返回的錯誤決定是 400 還是 500
		if errors.Is(err, service.ErrFailedToGenerateShortURL) {
			c.JSON(http.StatusConflict, gin.H{"error": "Could not generate a unique short URL, please try again."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		}
		return
	}

	// 返回結果
	c.JSON(http.StatusOK, gin.H{
		"short_url":  urlMapping.ShortURL,
		"algorithm":  urlMapping.Algorithm,
		"expires_at": urlMapping.ExpiresAt,
	})
}

// GetAllURLMappings (仍然需要認證)
func (h *URLHandler) GetAllURLMappings(c *gin.Context) {
	// 從 context 獲取 userID，因為這個路由受 AuthMiddleware 保護
	userIDValue, exists := c.Get("userID")
	if !exists {
		// 理論上不應該發生，因為 AuthMiddleware 應該已中止請求
		log.Println("Error: userID missing in context for authenticated route /url_mapping")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		log.Printf("Error: Invalid userID type in context for authenticated route /url_mapping: %T", userIDValue)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Printf("Fetching URL mappings for user ID: %d", userID)

	// TODO: 修改 Service 層和 Repository 層以支持按 UserID 過濾
	// mappings, err := h.urlService.GetURLMappingsByUserID(c.Request.Context(), userID)
	mappings, err := h.urlService.GetAllURLMappings(c.Request.Context()) // 暫時仍然獲取所有
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
		"data": mappings, // TODO: 應該只返回該用戶的數據
	})
}

// RedirectToOriginalURL 處理重定向到原始 URL 的請求
func (h *URLHandler) RedirectToOriginalURL(c *gin.Context) {
	shortURL := c.Param("shortURL")

	originalURL, err := h.urlService.GetOriginalURL(c.Request.Context(), shortURL)
	if err != nil {
		switch { // 使用 switch true 簡化多個 error 判斷
		case errors.Is(err, service.ErrURLNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"code": http.StatusNotFound,
				"msg":  "Short URL not found",
			})
		case errors.Is(err, service.ErrURLExpired):
			c.JSON(http.StatusGone, gin.H{
				"code": http.StatusGone,
				"msg":  "URL has expired",
			})
		default:
			log.Printf("Error retrieving original URL for %s: %v", shortURL, err)
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
