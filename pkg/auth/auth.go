package auth

import (
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"time"
)

// GenerateToken はjwtトークンを生成して返却します。
func GenerateToken(userID uint32) (string, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(24 * 30 * time.Hour)

	// JWTのクレームを設定
	claims := &jwt.StandardClaims{
		Issuer:    os.Getenv("APP_NAME"),
		IssuedAt:  time.Now().Unix(),
		Id:        strconv.FormatUint(uint64(userID), 10), // ユーザーID
		ExpiresAt: expirationTime.Unix(),                  // 有効期限
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵で署名
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
