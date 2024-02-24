package entity

import "time"

type Room struct {
	ChatRoomID uint32     `json:"chat_room_id"`
	UserName   string     `json:"user_name"`
	Message    *string    `json:"message"`
	CreatedAt  *time.Time `json:"created_at"`
}

type ChatRooms struct {
	List []Room
}
