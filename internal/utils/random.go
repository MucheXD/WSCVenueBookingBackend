package utils

import (
	"crypto/rand"
	"log/slog"
	"math/big"
)

func GenerateRandomInt(max int) int {
	rnd, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		slog.Error("Failed to generate random int", "error", err)
		return 0
	}
	return int(rnd.Int64())
}

func GenerateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[GenerateRandomInt(len(charset))]
	}
	return string(b)
}

// 生成指定长度的随机字符串，包含大小写字母和数字
func GenerateRandomStringULN(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return GenerateRandomString(length, charset)
}

// 生成指定长度的随机字符串，包含大小写字母
func GenerateRandomStringUL(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return GenerateRandomString(length, charset)
}

// 生成指定长度的随机字符串，包含大写字母和数字
func GenerateRandomStringUN(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return GenerateRandomString(length, charset)
}

// 生成指定长度的随机数字字符串
func GenerateRandomStringN(length int) string {
	const charset = "0123456789"
	return GenerateRandomString(length, charset)
}

// 生成指定长度的随机大写字母字符串
func GenerateRandomStringU(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return GenerateRandomString(length, charset)
}
