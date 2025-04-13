package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware 創建一個需要密鑰的強制認證中間件
func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc { // 接收 jwtSecret 參數
	if len(jwtSecret) == 0 {
		log.Fatal("CRITICAL: JWT secret is empty in AuthMiddleware configuration.") // 密鑰為空是嚴重錯誤
	}
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			log.Println("AuthMiddleware: Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// 檢查是否為 Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("AuthMiddleware: Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]

		// 解析和驗證 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 驗證簽名算法是否為預期的 HMAC (例如 HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("AuthMiddleware: Unexpected signing method: %v", token.Header["alg"])
				return nil, errors.New("unexpected signing method")
			}
			fmt.Println("jwtSecret", jwtSecret)
			// 返回用於驗證簽名的密鑰
			return jwtSecret, nil
		})

		if err != nil {
			log.Printf("AuthMiddleware: Error parsing token: %v", err)
			errorMsg := "Invalid or expired token"
			if errors.Is(err, jwt.ErrTokenExpired) {
				errorMsg = "Token has expired"
			} else if errors.Is(err, jwt.ErrTokenMalformed) {
				errorMsg = "Malformed token"
			} else if errors.Is(err, jwt.ErrSignatureInvalid) {
				errorMsg = "Invalid token signature"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
			return
		}

		// 檢查 Token 是否有效且 Claims 是否存在
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 從 Claims 中提取使用者 ID (或其他需要的信息)
			// 我們在生成 token 時將 User ID 存在 "sub" claim 中
			if sub, ok := claims["sub"]; ok {
				// JWT 標準中，ID 通常是 float64，需要轉換
				if userIDFloat, ok := sub.(float64); ok {
					userID := uint(userIDFloat) // 轉換為 uint
					// 將 userID 存儲到 Gin 的 Context 中，供後續 Handler 使用
					c.Set("userID", userID)
					log.Printf("AuthMiddleware: User %d authenticated", userID)
					c.Next() // Token 有效，繼續處理請求鏈
					return
				} else {
					log.Printf("AuthMiddleware: Invalid user ID type in token claims: %T", sub)
				}
			} else {
				log.Println("AuthMiddleware: User ID ('sub') claim missing in token")
			}
		} else {
			log.Println("AuthMiddleware: Invalid token or claims")
		}

		// 如果 token 無效或無法提取 UserID
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	}
}

// OptionalAuthMiddleware 創建一個需要密鑰的可選認證中間件
func OptionalAuthMiddleware(jwtSecret []byte) gin.HandlerFunc { // 接收 jwtSecret 參數
	if len(jwtSecret) == 0 {
		log.Println("Warning: JWT secret is empty in OptionalAuthMiddleware configuration. Optional auth may not work.")
		// 不直接 Fatal，但記錄警告
	}
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next() // 無標頭，繼續
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next() // 格式不對，繼續
			return
		}
		tokenString := parts[1]

		// 只有在密鑰有效時才嘗試解析
		if len(jwtSecret) > 0 {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return jwtSecret, nil // 使用傳入的 jwtSecret
			})

			// 解析錯誤或 Token 無效，不中止，只記錄並繼續
			if err != nil {
				log.Printf("OptionalAuthMiddleware: Invalid token provided: %v. Proceeding as anonymous.", err)
			} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Token 有效，嘗試提取 UserID
				if sub, ok := claims["sub"]; ok {
					if userIDFloat, ok := sub.(float64); ok {
						userID := uint(userIDFloat)
						c.Set("userID", userID) // 設置 UserID
						log.Printf("OptionalAuthMiddleware: User %d identified.", userID)
					}
				}
			}
		} else {
			log.Println("OptionalAuthMiddleware: Skipping token validation due to missing secret.")
		}

		c.Next() // 繼續處理請求
	}
}
