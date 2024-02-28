package entity

import "github.com/gorilla/websocket"

type Client struct {
	Ws WebSocketWrapper
}

func (client *Client) Send(msg []byte) error {
	return client.Ws.WriteMessage(websocket.TextMessage, msg)
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
