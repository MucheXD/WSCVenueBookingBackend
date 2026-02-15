package userSvc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
)

// 生成的 LoginTokenSalt 有效期
const loginTokenSaltValidDuration = 15 * time.Second

var secretLoadOnce sync.Once
var loginTokenSaltSecret string

// Check Login Token
func CheckLoginToken(user *models.User, loginToken string) error {
	// TODO: 实现登录 Token 的验证逻辑
	return nil
}

// 返回登录名对应的 Salt，若无对应用户返回随机值
func GetUserSalt(loginName string) (string, error) {
	result, err := repository.FoundUserByUsername(loginName)
	if err == nil && len(result) == 1 {
		return result[0].PasswordSalt, nil
	}
	randSalt := utils.GenerateRandomStringULN(16)
	if err != nil {
		return randSalt, err
	}
	if len(result) <= 0 { // 未找到用户，返回随机盐值
		return randSalt, nil
	}
	if len(result) > 1 { // 多个用户同名，记录警告日志并返回随机盐值
		slog.Warn("Multiple users found with the same username", "username", loginName)
		return randSalt, nil
	}
	return randSalt, nil
}

// Generate Login Session Salt
func GenerateLoginSessionSalt(ctx context.Context, loginName string) (string, error) {
	sessionID := utils.GenerateRandomStringUN(32)
	// 使用 Set 在 Redis 去重存储 sessionID
	// 如果发生碰撞，后正确登录者会失败
	if err := database.RDB.SAdd(ctx, database.KeyLoginSessionID, sessionID).Err(); err != nil {
		slog.Error("Failed to add login session ID to Redis set", "Error", err)
		return "", err
	}
	validBefore := time.Now().Add(loginTokenSaltValidDuration).Unix()
	if validBefore <= 0 || validBefore >= 1e11 {
		slog.Error("Generated login token salt valid before is out of range", "validBefore", validBefore)
		return "", fmt.Errorf("generated login token salt valid before is out of range: %d", validBefore)
	}
	saltData := fmt.Sprintf("%010d", validBefore) + sessionID
	loadSecret()
	if loginTokenSaltSecret == "" {
		slog.Error("Login token salt secret is not configured")
		return "", fmt.Errorf("login token salt secret is not configured")
	}
	return saltData + generateSignature(saltData, loginTokenSaltSecret), nil
}

func generateSignature(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sha)
}

func loadSecret() {
	secretLoadOnce.Do(func() {
		loginTokenSaltSecret = config.Config.GetString("service.secret.loginTokenSalt")
	})
}
