package v2

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"business/internal/entity"
	"business/internal/usecase/room"
	"business/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type MessageRoutes struct {
	t usecase.Message
	l logger.Interface
}

// NewMessageRouter はチャット関連のURLからコントローラーを実行します。
func NewMessageRouter(e *echo.Echo, t usecase.Message, l logger.Interface) {
	routes := &MessageRoutes{t, l}
	e.GET("/chats/:id", routes.handleConnections)
}

// グローバルに宣言するのは、関数内に記述するとアクセス毎に初期されるため。
var upGrade websocket.Upgrader
var roomManager map[uint32]*entity.ChatRoom

// init関数を使用してroomManagerを初期化
func init() {
	roomManager = make(map[uint32]*entity.ChatRoom)
}

// handleConnections GETリクエストをWebSocketへアップグレードし、Pub/Subを管理します。
func (r *MessageRoutes) handleConnections(c echo.Context) error {
	chatRoomID, err := validate(c)
	if err != nil {
		fmt.Printf("不正なリクエストパラメータ \n")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "不正なリクエストパラメータです。"})
	}

	// WebSocket　へアップグレード
	upGrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upGrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
	}

	defer func() {
		if err := ws.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}()

	// リクエストをルームに参加させる。
	err = r.t.JoinRoom(chatRoomID, ws, roomManager)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// ルーム参加者からのメッセージを検知し送信する。
	r.t.PubSub(c, ws, roomManager, chatRoomID)

	return nil
}

// validate はリクエストパラメータののバリデーションを行う。
func validate(c echo.Context) (uint32, error) {
	chatIDStr := c.Param("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		err := errors.New("不正なリクエストパラメータ")
		return 0, err
	}

	return uint32(chatID), nil
}
