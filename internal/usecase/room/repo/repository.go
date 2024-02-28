package repo

import (
	"fmt"

	"business/internal/entity"
	"business/pkg/mysql"
	wt "business/pkg/time"
	"gorm.io/gorm"
)

type ChatRepository struct {
	DB *gorm.DB
	wt wt.WrapperTime
}

// New は新しいRepositoryインスタンスを生成します。
func New(db *mysql.MySQL, tp wt.WrapperTime) *ChatRepository {
	return &ChatRepository{DB: db.DB, wt: tp} // MySQL構造体のDBフィールドを使ってRepositoryを初期化
}

// GetMessageList はチャットに紐づくメッセージ一覧を取得します。
func (r *ChatRepository) GetMessageList(chatRoomID uint32) ([]entity.Message, error) {
	var messageList []entity.Message
	err := r.DB.Table("chat_messages cm").
		Select(`
			cm.chat_message_id,
			cm.sender_user_id,
            CONCAT(u.last_name, u.first_name) as user_name,
			cm.message,
			cm.created_at
		`).
		Joins(`INNER JOIN users u ON u.user_id = cm.sender_user_id`).
		Where(`cm.chat_room_id = ?`, chatRoomID).
		Order(`cm.chat_message_id ASC`).
		Find(&messageList).
		Error

	if err != nil {
		return []entity.Message{}, fmt.Errorf("DB serveer Error : %w", err)
	}

	return messageList, nil
}

// CreateMessage はメッセージを作成します。
func (r *ChatRepository) CreateMessage(chatRoomID uint32, msg string, SenderUserID uint32) (entity.ChatMessage, error) {
	chatMessage := entity.ChatMessage{
		ChatMessageID: 0, // オートインクリメントで自動で設定されるが、明示的に指定しておく。
		ChatRoomID:    chatRoomID,
		Message:       msg,
		SenderUserID:  SenderUserID,
		CreatedAt:     r.wt.Now(),
	}
	err := r.DB.Create(&chatMessage).Error

	if err != nil {
		return entity.ChatMessage{}, fmt.Errorf("DB serveer Error : %w", err)
	}

	return chatMessage, nil
}
