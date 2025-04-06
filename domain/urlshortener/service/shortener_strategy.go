package service

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

// ShortenerStrategy 定義了 URL 縮短算法的介面
type ShortenerStrategy interface {
	Generate(input string, id int) *string
}

// Base62Strategy 使用 Base62 編碼實現 ShortenerStrategy
type Base62Strategy struct{}

// Base64Strategy 使用 Base64 編碼實現 ShortenerStrategy
type Base64Strategy struct{}

// MD5Strategy 使用 MD5 哈希實現 ShortenerStrategy
type MD5Strategy struct{}

// RandomStrategy 使用隨機字符實現 ShortenerStrategy
type RandomStrategy struct{}

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Generate 實現 Base62Strategy 的 Generate 方法
func (s *Base62Strategy) Generate(input string, id int) *string {
	return DecimalToBase62(id)
}

// Generate 實現 Base64Strategy 的 Generate 方法
func (s *Base64Strategy) Generate(input string, id int) *string {
	// 將 URL 轉換為 Base64 編碼並取前 8 個字符
	idStr := strings.TrimSpace(strings.Replace(base64.StdEncoding.EncodeToString([]byte(input)), "=", "", -1))
	// 取前 8 個字符，如果字符串較短則取全部
	length := 8
	if len(idStr) < length {
		length = len(idStr)
	}
	result := idStr[:length]
	return &result
}

// Generate 實現 MD5Strategy 的 Generate 方法
func (s *MD5Strategy) Generate(input string, id int) *string {
	// 創建 URL + ID 的 MD5 哈希
	hasher := md5.New()
	hasher.Write([]byte(input + string(rune(id))))
	hashStr := hex.EncodeToString(hasher.Sum(nil))
	// 取前 8 個字符
	result := hashStr[:8]
	return &result
}

// Generate 實現 RandomStrategy 的 Generate 方法
func (s *RandomStrategy) Generate(input string, id int) *string {
	// 初始化隨機數生成器
	rand.Seed(time.Now().UnixNano())
	// 生成 8 個隨機字符
	var result strings.Builder
	for i := 0; i < 8; i++ {
		randomIndex := rand.Intn(len(charset))
		result.WriteByte(charset[randomIndex])
	}
	resultStr := result.String()
	return &resultStr
}

// DecimalToBase62 將十進制數轉換為 Base62 字符串
func DecimalToBase62(decimalNum int) *string {
	var result strings.Builder
	base := 62

	for decimalNum > 0 {
		remainder := decimalNum % base
		result.WriteByte(charset[remainder])
		decimalNum /= base
	}

	// 反轉字符串，因為我們是從右到左構建的
	reversed := result.String()
	var final strings.Builder
	for i := len(reversed) - 1; i >= 0; i-- {
		final.WriteByte(reversed[i])
	}

	resultStr := final.String()
	if resultStr == "" {
		resultStr = "0"
	}
	return &resultStr
}
