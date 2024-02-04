package repo

import (
	"business/internal/entity"
	"business/pkg/mysql"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// New は新しいRepositoryインスタンスを生成します。
func New(db *mysql.MySQL) *UserRepository {
	return &UserRepository{DB: db.DB} // MySQL構造体のDBフィールドを使ってRepositoryを初期化
}

// GetUserList はユーザーリストを取得します。
func (r *UserRepository) GetUserList(ctx context.Context) ([]entity.User, error) {
	var userList []entity.User
	err := r.DB.Table("user").
		Find(&userList).
		Error // GORMのTableメソッドを使用

	if err != nil {
		return nil, fmt.Errorf("failed to get user list: %w", err)
	}
	return userList, nil
}
