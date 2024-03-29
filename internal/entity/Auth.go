package entity

import (
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Mail     string `json:"mail" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginUser -.
type LoginUser struct {
	UserID            uint32 `gorm:"column:user_id"`
	EncryptedPassword string `gorm:"column:password"`
}

func (encryptedPassWord *LoginUser) IsValidPassword(password string) bool {
	// bcryptを使用してハッシュ化されたパスワードと平文のパスワードを比較
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassWord.EncryptedPassword), []byte(password))
	// エラーがなければ、パスワードは一致しているとみなす
	return err == nil
}
