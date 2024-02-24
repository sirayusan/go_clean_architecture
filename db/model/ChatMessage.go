package model

import "time"

type ChatMessage struct {
	ChatMessageID uint32     `gorm:"column:chat_message_id;primaryKey;autoIncrement:true;comment:チャットメッセージID;"`
	ChatRoomID    uint32     `gorm:"column:chat_room_id;not null;comment:チャットルームID;"`
	Message       string     `gorm:"column:message;not null;comment:メッセージ内容;"`
	SenderUserID  uint32     `gorm:"column:sender_user_id;not null;comment:送信者ID;"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;type:timestamp;comment:登録日時;"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;type:datetime;comment:削除日時;"`
}
