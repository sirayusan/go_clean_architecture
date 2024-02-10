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
func (r *AuthRepository) GetUserByMail(email string) (entity.LoginUser, error) {
	var user entity.LoginUser
	err := r.DB.Table("users").
		Select(
			"user_id",
			"password",
		).
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.LoginUser{}, err
	}
	if err != nil {
		return entity.LoginUser{}, fmt.Errorf("DB serveer Error : %w", err)
	}
	return user, nil
}
