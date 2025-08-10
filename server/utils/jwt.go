package utils

import (
	"server/global"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 定义JWT声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成访问令牌
// 生成访问令牌
func GenerateToken(userID uint, username string) (string, error) {
    // 从配置获取过期时间
    expireTime := time.Now().Add(time.Duration(global.Config.Jwt.AccessTokenExpiryTime) * time.Minute)
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}

	// 使用访问令牌密钥签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.Config.Jwt.AccessTokenSecret))
}

// 生成刷新令牌
func GenerateRefreshToken(userID uint, username string) (string, error) {
	// 从配置获取刷新令牌过期时间
	expireTime := time.Now().Add(time.Duration(global.Config.Jwt.RefreshTokenExpiryTime) * time.Minute)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}

	// 使用刷新令牌密钥签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.Config.Jwt.RefreshTokenSecret))
}

// 解析令牌
func ParseToken(tokenString string, isRefresh bool) (*Claims, error) {
	// 根据令牌类型选择对应的密钥
	secret := global.Config.Jwt.AccessTokenSecret
	if isRefresh {
		secret = global.Config.Jwt.RefreshTokenSecret
	}

	// 解析令牌
	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
