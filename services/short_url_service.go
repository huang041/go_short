package services

import (
	"strings"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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

	resultStr := result.String()
	return &resultStr
}
