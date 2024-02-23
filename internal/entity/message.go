package entity

import "time"

type MessageResponse struct {
	UserID uint32     `json:":user_id"`
	List   []Messages `json:"list"`
}

// Message -.
type Message struct {
	UserName  string    `gorm:"column:user_name"  json:":user_name"`
	Messages  string    `gorm:"column:message" json:":message"`
	CreatedAt time.Time `gorm:"column:created_at" json:":created_at"`
}

// Messages -.
type Messages struct {
	List []Message `json:"list"`
}
