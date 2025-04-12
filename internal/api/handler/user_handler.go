package handler

import (
	"errors" // 引入 errors
	"net/http"

	identityapp "go_short/internal/application/identity" // 引入 Identity 應用服務

	"github.com/gin-gonic/gin"
)

// UserHandler 處理使用者相關的 HTTP 請求
type UserHandler struct {
	identityApp *identityapp.App // 依賴 Identity 應用服務
}

// NewUserHandler 創建 User Handler 實例
func NewUserHandler(identityApp *identityapp.App) *UserHandler {
	return &UserHandler{
		identityApp: identityApp,
	}
}

// Register 處理使用者註冊請求
func (h *UserHandler) Register(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user, err := h.identityApp.RegisterUser(c.Request.Context(), request.Username, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, identityapp.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			// 記錄內部錯誤
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	// 註冊成功，返回部分使用者資訊（避免洩漏密碼雜湊等）
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login 處理使用者登入請求
func (h *UserHandler) Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user, err := h.identityApp.AuthenticateUser(c.Request.Context(), request.Username, request.Password)
	if err != nil {
		// 對於認證失敗，統一返回未授權錯誤
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 登入成功
	// 在實際應用中，這裡通常會生成一個 JWT 或 Session Token 返回給客戶端
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user_id": user.ID, // 示例：返回用戶 ID
		// "token": "your_generated_jwt_token", // 返回 token
	})
}

// --- 可以添加其他 Handler 方法，如 GetProfile, Logout 等 ---
