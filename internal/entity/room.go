package entity

import (
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	Ws *websocket.Conn
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
	UserID uint32     `json:":user_id"`
	List   []Messages `json:"list"`
}

// Message -.
type Message struct {
	UserName  string    `gorm:"column:user_name"  json:"user_name"`
	Messages  string    `gorm:"column:message" json:"message"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// Messages -.
type Messages struct {
	List []Message `json:"list"`
}
