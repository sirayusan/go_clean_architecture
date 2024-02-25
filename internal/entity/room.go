package entity

import (
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	Ws WebSocketConn // WebSocketWrapper から WebSocketConn へ変更
}

type ChatRoom struct {
	Clients []*Client
}

func (rooms *ChatRoom) AddClient(client *Client) {
	rooms.Clients = append(rooms.Clients, client)
}

func (rooms *ChatRoom) GetClients() []Client {
	var c []Client
	for _, client := range rooms.Clients {
		c = append(c, *client)
	}

	return c
}

func (rooms *ChatRoom) Publish(msg []byte) {
	for _, client := range rooms.Clients {
		client.Send(msg)
	}
}

func (client *Client) Send(msg []byte) error {
	return client.Ws.WriteMessage(websocket.TextMessage, msg)
}

// MessageResponse -.
type MessageResponse struct {
	UserID uint32    `json:":user_id"`
	List   []Message `json:"list"`
}

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

type WebSocketConn interface {
	WriteMessage(messageType int, data []byte) error
}

type WebSocketWrapper struct {
	Conn *websocket.Conn
}

func (w *WebSocketWrapper) WriteMessage(messageType int, data []byte) error {
	return w.Conn.WriteMessage(messageType, data)
}
