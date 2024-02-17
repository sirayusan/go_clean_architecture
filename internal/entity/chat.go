package entity

import "time"

type Chat struct {
	UserName  string     `json:"user_name"`
	Message   *string    `json:"message"`
	CreatedAt *time.Time `json:"created_at"`
}

type Chats struct {
	List []Chat
}
