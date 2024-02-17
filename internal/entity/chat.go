package entity

import "time"

type Chat struct {
	RecipientUserName string     `json:"recipient_user_name"`
	Message           *string    `json:"message"`
	CreateAt          *time.Time `json:"create_at"`
}

type Chats struct {
	List []Chat
}
