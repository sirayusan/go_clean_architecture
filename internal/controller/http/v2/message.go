package v2

import (
	"business/internal/usecase/message"
	"business/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MessageRoutes struct {
	t usecase.Message
	l logger.Interface
}

type Client struct {
	Ws *websocket.Conn
}
type Rooms struct {
	Clients []*Client
}

func (rooms *Rooms) AddClient(client *Client) {
	rooms.Clients = append(rooms.Clients, client)
}

func (rooms *Rooms) GetClients() []Client {
	var cs []Client

	for _, client := range rooms.Clients {
		cs = append(cs, *client)
	}

	return cs
}

func (rooms *Rooms) Publish(msg []byte) {
	for _, client := range rooms.Clients {
		client.Send(msg) // この後に作る
	}
}

func (client *Client) Send(msg []byte) error {
	return client.Ws.WriteMessage(websocket.TextMessage, msg)
}

var upgrader = websocket.Upgrader{}
var rooms = Rooms{}

// NewMessageRouter はチャット関連のURLからコントローラーを実行します。
func NewMessageRouter(e *echo.Echo, t usecase.Message, l logger.Interface) {
	//routes := &MessageRoutes{t, l}
	e.GET("/message", handleConnections)
}

func handleConnections(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		c.Logger().Error(err)
	}

	defer ws.Close()

	client := &Client{Ws: ws}

	rooms.AddClient(client)

	for {
		_, msg, err := ws.ReadMessage()

		if err != nil {
			c.Logger().Error(err)
			break
		}

		rooms.Publish(msg)
	}

	return nil
}
