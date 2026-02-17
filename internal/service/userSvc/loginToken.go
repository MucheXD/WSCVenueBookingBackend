package userSvc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/redis/go-redis/v9"
)

// LoginToken 构成
// 1-10 位：有效期截止时间戳（UNIX秒）
// 11-42 位：SessionID（随机字符串）
// 43-106 位：签名（HMAC-SHA256-HEX）
// 107- 位：R=Hash(LoginSessionSalt+Hash(Password+Salt))

// 生成的 LoginTokenSalt 有效期
const loginTokenSaltValidDuration = 15 * time.Second

var secretLoadOnce sync.Once
var loginTokenSaltSecret string

func TryPasswordLogin(ctx context.Context, loginName, loginToken string) (bool, error) {
	userFound, err := repository.GetUniqueUserByUsername(loginName)
	if err != nil {
		return false, err
	}
	// 此处可以添加对 Email 登录的尝试，此处跳过
	if userFound == nil {
		return false, nil
	}
	return checkLoginToken(ctx, userFound, loginToken)
}

// Check Login Token
func checkLoginToken(ctx context.Context, user *models.User, loginToken string) (bool, error) {

	// 检查步骤需按照：
	// 1. 验证 Token 格式和有效期，如果无效可避免后续步骤
	// 2. 验证签名，确保 Token 未被篡改
	// 3. 验证 SessionID 是否存在于 Redis 中，确保 Token 是由系统生成的
	// 如果先验证 SessionID，攻击者可以通过暴力生成大量 Token 来探测有效的 SessionID

	if len(loginToken) < 106 {
		return false, fmt.Errorf("%w: %d", ErrInvalidLoginTokenLength, len(loginToken))
	}
	validBeforeStr := loginToken[:10]
	validBefore, err := strconv.ParseInt(validBeforeStr, 10, 64)
	if err != nil {
		return false, err
	}
	if validBefore < time.Now().Unix() {
		return false, fmt.Errorf("%w at %d", ErrLoginTokenExpired, validBefore)
	}
	sessionID := loginToken[10:42]
	signature := loginToken[42:106]
	saltData := validBeforeStr + sessionID

	loadSecret()
	if loginTokenSaltSecret == "" {
		return false, ErrLoginTokenSaltSecretNotConfigured
	}

	// 验证签名,确保 Token 未被篡改
	expectedSignature := generateSignature(saltData, loginTokenSaltSecret)
	if signature != expectedSignature {
		return false, ErrLoginTokenSignatureMismatch
	}

	// 检查 SessionID 是否存在于 Redis 中
	val, err := database.RDB.GetDel(ctx, database.KeyLoginSessionID(sessionID)).Result()
	if err == redis.Nil {
		return false, ErrLoginSessionIDInvalid
	} else if err != nil {
		return false, fmt.Errorf("%w: %w", ErrCheckLoginSessionIDInRedis, err)
	}
	if val != user.Username { // 若有其它登录方式，此处需要追加对比逻辑
		// 存在但不匹配，返回无效而不是其它错误避免被探测
		return false, ErrLoginSessionIDInvalid
	}

	// base64 解码传入 R 值
	hashedPwdHash, err := base64.StdEncoding.DecodeString(loginToken[106:])
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrDecodeLoginTokenPasswordHash, err)
	}
	// 通过 loginToken 内的 loginSessionSalt 加盐计算正确 R 值并比较
	expectedHash := sha256.Sum256([]byte(loginToken[:106] + user.PasswordHash))
	if !hmac.Equal(hashedPwdHash, expectedHash[:]) {
		return false, nil // 密码不匹配算预期发生的行为
	}

	return true, nil
}

// 返回登录名对应的 Salt，若无对应用户返回随机值
func GetUserSalt(loginName string) (string, error) {
	result, err := repository.GetUniqueUserByUsername(loginName)
	if err == nil && result != nil {
		return result.PasswordSalt, nil
	}
	randSalt := utils.GenerateRandomStringULN(16)
	if err != nil {
		return randSalt, err
	}
	return randSalt, nil
}

// Generate Login Session Salt
func GenerateLoginSessionSalt(ctx context.Context, loginName string) (string, error) {
	sessionID := utils.GenerateRandomStringUN(32)
	// 在 Redis 去重存储 sessionID，有效期等同 SessionSalt 有效期
	if err := database.RDB.SetNX(ctx,
		database.KeyLoginSessionID(sessionID),
		loginName,
		loginTokenSaltValidDuration).Err(); err != nil {
		return "", fmt.Errorf("%w: %w", ErrStoreLoginSessionIDInRedis, err)
	}
	// 生成 LoginTokenSalt
	validBefore := time.Now().Add(loginTokenSaltValidDuration).Unix()
	if validBefore <= 0 || validBefore >= 1e11 {
		return "", fmt.Errorf("%w: %d", ErrLoginTokenSaltValidBeforeOutOfRange, validBefore)
	}
	saltData := fmt.Sprintf("%10d", validBefore) + sessionID
	loadSecret()
	if loginTokenSaltSecret == "" {
		return "", ErrLoginTokenSaltSecretNotConfigured
	}
	return saltData + generateSignature(saltData, loginTokenSaltSecret), nil
}

func generateSignature(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := h.Sum(nil)
	return hex.EncodeToString(sha)
}

// loginTokenSaltSecret 懒加载函数，每次使用前调用
func loadSecret() {
	secretLoadOnce.Do(func() {
		loginTokenSaltSecret = config.Config.GetString("service.secret.login_token_salt")
	})
}
