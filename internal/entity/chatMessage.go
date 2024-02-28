package entity

import (
	"time"
)

// Message -.
type Message struct {
	SenderUserID uint32    `gorm:"column:sender_user_id"  json:"sender_user_id"`
	UserName     string    `gorm:"column:user_name"  json:"user_name"`
	Messages     string    `gorm:"column:message" json:"message"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

// ChatMessage -.
type ChatMessage struct {
	ChatMessageID uint32    `gorm:"column:chat_message_id;primaryKey;"`
	ChatRoomID    uint32    `gorm:"column:chat_room_id"`
	Message       string    `gorm:"column:message"`
	SenderUserID  uint32    `gorm:"column:sender_user_id"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

type RedisMessage struct {
	ServerId  string
	Timestamp time.Time
	MessageId uint32
	Payload   string
}
