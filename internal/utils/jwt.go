package utils

import (
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenData struct {
	UserID             uint `json:"uid"`
	VenueAccessGroupID uint `json:"vag,omitempty"`
}
type tokenClaims struct {
	TokenData
	jwt.RegisteredClaims
}

// 生成token
func GenerateToken(tokenData TokenData) (string, error) {
	claims := tokenClaims{
		TokenData: tokenData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(config.Config.GetInt("jwt.expiration")))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.Config.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func GetTokenData(tokenStr string) (TokenData, error) {
	var claims tokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Config.GetString("jwt.secret")), nil
	})
	if err != nil {
		return TokenData{}, err
	}
	if !token.Valid {
		// jwt.Parse 通过而 Valid 为 false 不是一个预期行为
		// 这里做一个特殊的错误返回
		return TokenData{}, jwt.ErrTokenNotValidYet
	}
	return claims.TokenData, nil
}
