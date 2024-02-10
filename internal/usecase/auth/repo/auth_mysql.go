package repo

import (
	"business/internal/entity"
	"business/pkg/mysql"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

// New は新しいRepositoryインスタンスを生成します。
func New(db *mysql.MySQL) *AuthRepository {
	return &AuthRepository{DB: db.DB} // MySQL構造体のDBフィールドを使ってRepositoryを初期化
}

// GetUserByMail はユーザーリストを取得します。
func (r *AuthRepository) GetUserByMail(email string) (entity.LoginUserPassWord, error) {
	var user entity.LoginUserPassWord
	err := r.DB.Table("users").
		Select(
			"password",
		).
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.LoginUserPassWord{}, err
	}
	if err != nil {
		return entity.LoginUserPassWord{}, fmt.Errorf("DB server Error: %w", err)
	}
	return user, nil
}
