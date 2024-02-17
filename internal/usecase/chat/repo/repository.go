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

// GetChatList はチャットリストを取得します。
func (r *ChatRepository) GetChatList(userID uint32) (entity.Chats, error) {
	var chatList []entity.Chat
	err := r.DB.Table("chats c").
		Select(`
			DISTINCT
			c.chat_id,
            CASE
               WHEN c.user_id1 = ? THEN CONCAT(u2.last_name, u2.first_name)
               ELSE CONCAT(u.last_name, u.first_name)
            END AS user_name,
			cm.message,
			cm.created_at
		`, userID).
		Joins(`INNER JOIN (
        	    SELECT
        	    	cm.chat_id,
        	    	MAX(cm.chat_message_id) AS max_chat_message_id
        	    FROM
        	    	chat_messages cm
        	    GROUP BY
        	    	cm.chat_id
           ) AS latest_cm ON c.chat_id = latest_cm.chat_id
        `).
		Joins(`INNER JOIN users u ON 
            u.user_id = c.user_id1
        `).
		Joins(`INNER JOIN users u2 ON 
            u2.user_id = c.user_id2
        `).
		Joins(`INNER JOIN chat_messages cm ON
            c.chat_id = cm.chat_id AND
            cm.chat_message_id = latest_cm.max_chat_message_id
        `).
		Where(`c.user_id1 = ? OR c.user_id2 = ?`, userID, userID).
		Find(&chatList).
		Error

	if err != nil {
		return entity.Chats{}, fmt.Errorf("DB serveer Error : %w", err)
	}
	if len(chatList) == 0 {
		return entity.Chats{}, gorm.ErrRecordNotFound
	}

	return entity.Chats{List: chatList}, nil
}
