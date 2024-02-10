package usecase

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"

	"business/internal/entity"
)

// AuthUseCase -.
type AuthUseCase struct {
	repo AuthRepo
}

// New -.
func New(r AuthRepo) *AuthUseCase {
	return &AuthUseCase{
		repo: r,
	}
}

// Authentication は入力値の顧客が存在したら、jwtトークンを生成して返す。
func (uc *AuthUseCase) GenerateJwtToken(param entity.LoginRequest) (string, error) {
	user, err := uc.repo.GetUserByMail(param.Mail)
	// DBと疎通できずエラーなのか、存在せずエラー(401)を分ける必要がある。
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("authentication - s.repo.GetUserByMail: %w", err)
	}

	// isValidPassword　DBから取得した暗号化パスワードとリクエストのパスワードを暗号化して比較した結果、一致したらtrueが入る。
	isValidPassword := user.IsValidPassword(param.Password)
	if isValidPassword {
		token, err := generateToken(user.UserID)
		if err != nil {
			return "", err
		}
		return token, nil
	} else {
		return "", errors.New("invalid password")
	}
}

func generateToken(userID uint32) (string, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(24 * time.Hour)

	// JWTのクレームを設定
	claims := &jwt.StandardClaims{
		Issuer:    os.ExpandEnv(":${APP_NAME}"),
		IssuedAt:  time.Now().Unix(),
		Id:        strconv.FormatUint(uint64(userID), 10), // ユーザーID
		ExpiresAt: expirationTime.Unix(),                  // 有効期限
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵で署名
	tokenString, err := token.SignedString([]byte(os.ExpandEnv(":${JWT_SECRET_KEY}")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
