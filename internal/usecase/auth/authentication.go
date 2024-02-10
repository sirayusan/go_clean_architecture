package usecase

import (
	"errors"
	"fmt"
	"gorm.io/gorm"

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
func (uc *AuthUseCase) Authentication(param entity.LoginRequest) (string, error) {
	loginUserPassWord, err := uc.repo.GetUserByMail(param.Mail)
	// DBと疎通できずエラーなのか、存在せずエラー(401)を分ける必要がある。
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("authentication - s.repo.GetUserByMail: %w", err)
	}

	// isValidPassword　DBから取得した暗号化パスワードとリクエストのパスワードを暗号化して比較した結果、一致したらtrueが入る。
	isValidPassword := loginUserPassWord.IsValidPassword(param.Password)
	if isValidPassword {
		// TODO jwtトークンを生成して返す。
		return "", nil
	} else {
		return "", errors.New("invalid password")
	}
}
