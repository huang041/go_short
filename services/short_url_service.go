package services

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

// ShortenerStrategy defines the interface for URL shortening algorithms
type ShortenerStrategy interface {
	Generate(input string, id int) *string
}

// Base62Strategy implements ShortenerStrategy using Base62 encoding
type Base62Strategy struct{}

// Base64Strategy implements ShortenerStrategy using Base64 encoding
type Base64Strategy struct{}

// MD5Strategy implements ShortenerStrategy using MD5 hash
type MD5Strategy struct{}

// RandomStrategy implements ShortenerStrategy using random characters
type RandomStrategy struct{}

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// URLShortener is the main service for shortening URLs
type URLShortener struct {
	strategy ShortenerStrategy
}

// NewURLShortener creates a new URL shortener with the specified strategy
func NewURLShortener(strategy ShortenerStrategy) *URLShortener {
	return &URLShortener{
		strategy: strategy,
	}
}

// SetStrategy changes the shortening strategy
func (s *URLShortener) SetStrategy(strategy ShortenerStrategy) {
	s.strategy = strategy
}

// ShortenURL generates a short URL using the current strategy
func (s *URLShortener) ShortenURL(originalURL string, id int) *string {
	return s.strategy.Generate(originalURL, id)
}

// Generate implements ShortenerStrategy for Base62Strategy
func (s *Base62Strategy) Generate(input string, id int) *string {
	return DecimalToBase62(id)
}

// Generate implements ShortenerStrategy for Base64Strategy
func (s *Base64Strategy) Generate(input string, id int) *string {
	// Convert ID to string and encode with base64
	idStr := strings.TrimSpace(strings.Replace(base64.StdEncoding.EncodeToString([]byte(input)), "=", "", -1))
	// Take first 8 characters or less if string is shorter
	length := 8
	if len(idStr) < length {
		length = len(idStr)
	}
	result := idStr[:length]
	return &result
}

// Generate implements ShortenerStrategy for MD5Strategy
func (s *MD5Strategy) Generate(input string, id int) *string {
	// Create MD5 hash of the original URL + ID
	hasher := md5.New()
	hasher.Write([]byte(input + string(rune(id))))
	hashStr := hex.EncodeToString(hasher.Sum(nil))
	// Take first 8 characters
	result := hashStr[:8]
	return &result
}

// Generate implements ShortenerStrategy for RandomStrategy
func (s *RandomStrategy) Generate(input string, id int) *string {
	// Initialize random number generator with seed
	rand.Seed(time.Now().UnixNano())
	// Generate 8 random characters
	var result strings.Builder
	for i := 0; i < 8; i++ {
		randomIndex := rand.Intn(len(charset))
		result.WriteByte(charset[randomIndex])
	}
	resultStr := result.String()
	return &resultStr
}

// DecimalToBase62 converts a decimal number to a base62 string
func DecimalToBase62(decimalNum int) *string {
	var result strings.Builder
	base := 62

	for decimalNum > 0 {
		remainder := decimalNum % base
		result.WriteByte(charset[remainder])
		decimalNum /= base
	}

	// Reverse the string because we built it from right to left
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
