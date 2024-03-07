package usecase

import (
	"business/internal/entity"
	"business/pkg/auth"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

type StandardClaims struct {
	Issuer    string
	IssuedAt  uint32
	Id        string // ユーザーID
	ExpiresAt uint32 // 有効期限
	abc       string
}

// GenerateJwtToken は入力値の顧客が存在したら、jwtトークンを生成して返す。
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
		token, err := auth.GenerateToken(user.UserID)
		if err != nil {
			return "", err
		}
		return token, nil
	} else {
		return "", errors.New("invalid password")
	}
}
