package repo

import (
	"fmt"
	"gorm.io/gorm"

	"business/internal/entity"
	"business/pkg/mysql"
)

type UserRepository struct {
	DB *gorm.DB
}

// New は新しいRepositoryインスタンスを生成します。
func New(db *mysql.MySQL) *UserRepository {
	return &UserRepository{DB: db.DB} // MySQL構造体のDBフィールドを使ってRepositoryを初期化
}

// GetUserList はユーザーリストを取得します。
func (r *UserRepository) GetUserList() ([]entity.User, error) {
	var userList []entity.User
	err := r.DB.Table("users").
		Select(
			"user_id",
			"last_name",
			"first_name",
		).
		Find(&userList).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user list: %w", err)
	}
	return userList, nil
}
