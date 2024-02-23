package v2

import (
	"business/internal/usecase/message"
	"business/pkg/logger"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

var upGrade = websocket.Upgrader{}
var rooms = Rooms{}

// NewMessageRouter はチャット関連のURLからコントローラーを実行します。
func NewMessageRouter(e *echo.Echo, t usecase.Message, l logger.Interface) {
	routes := &MessageRoutes{t, l}
	e.GET("/chats/:id", routes.handleConnections)
}

func (r *MessageRoutes) handleConnections(c echo.Context) error {
	chatIDStr := c.Param("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		fmt.Printf("不正なリクエストパラメータ \n")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "不正なリクエストパラメータです。"})
	}

	upGrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upGrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
	}

	defer ws.Close()
	client := &Client{Ws: ws}
	rooms.AddClient(client)

	MessagesList, err := r.t.GetMessages(uint32(chatID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// TODO 受信メッセージを画面にすべてJson形式で返す。
	for _, data := range MessagesList.List {
		client.Send([]byte(data.Messages))
	}

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
