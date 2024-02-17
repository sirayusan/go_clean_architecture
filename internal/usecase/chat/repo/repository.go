package repo

import (
	"business/internal/entity"
	"business/pkg/mysql"
	"gorm.io/gorm"
)

type ChatRepository struct {
	DB *gorm.DB
}

// New は新しいRepositoryインスタンスを生成します。
func New(db *mysql.MySQL) *ChatRepository {
	return &ChatRepository{DB: db.DB} // MySQL構造体のDBフィールドを使ってRepositoryを初期化
}

// GetChatList はチャットリストを取得します。
func (r *ChatRepository) GetChatList(userID uint32) (entity.Chats, error) {
	var chats entity.Chats
	//err := r.DB.Table("users").
	//	Select(
	//		"user_id",
	//		"password",
	//	).
	//	Where("email = ?", userID).
	//	First(&user).
	//	Error

	//if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	//	return entity.Chat{}, err
	//}
	//if err != nil {
	//	return entity.Chat{}, fmt.Errorf("DB serveer Error : %w", err)
	//}
	return chats, nil
}
