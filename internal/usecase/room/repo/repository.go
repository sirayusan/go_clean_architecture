package repo

import (
	"business/internal/entity"
	"business/pkg/mysql"
	"fmt"
	"gorm.io/gorm"
)

type ChatRepository struct {
	DB *gorm.DB
}

// New は新しいRepositoryインスタンスを生成します。
func New(db *mysql.MySQL) *ChatRepository {
	return &ChatRepository{DB: db.DB} // MySQL構造体のDBフィールドを使ってRepositoryを初期化
}

// GetMessageList はチャットに紐づくメッセージ一覧を取得します。
func (r *ChatRepository) GetMessageList(chatRoomID uint32) ([]entity.Message, error) {
	var messageList []entity.Message
	err := r.DB.Table("chat_messages cm").
		Select(`
			cm.chat_message_id,
            CONCAT(u.last_name, u.first_name) as user_name,
			cm.message,
			cm.created_at
		`).
		Joins(`INNER JOIN users u ON u.user_id = cm.sender_user_id`).
		Where(`cm.chat_room_id = ?`, chatRoomID).
		Find(&messageList).
		Error

	if err != nil {
		return []entity.Message{}, fmt.Errorf("DB serveer Error : %w", err)
	}

	return messageList, nil
}
