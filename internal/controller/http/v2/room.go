package v2

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"business/internal/entity"
	"business/internal/usecase/room"
	"business/pkg/logger"
)

type MessageRoutes struct {
	t   usecase.Message
	l   logger.Interface
	rdb entity.RedisWrapper
}

// NewMessageRouter はチャット関連のURLからコントローラーを実行します。
func NewMessageRouter(e *echo.Echo, t usecase.Message, l logger.Interface, r entity.RedisWrapper) {
	routes := &MessageRoutes{t, l, r}
	e.GET("/chats/:id", routes.handleConnections, websocketJwtMiddleware())
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
	chatIDStr, _ := strconv.Atoi(c.Param("id"))
	chatRoomID := uint32(chatIDStr)

	// WebSocket　へアップグレード
	upGrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upGrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	defer func() {
		if err := ws.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}()
	subscribe := r.rdb.Subscribe(c.Request().Context(), "room-"+fmt.Sprint(chatRoomID))
	defer func() {
		if err := subscribe.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}()

	// リクエストをルームに参加させる。
	err = r.t.JoinRoom(chatRoomID, ws, roomManager)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	r.t.PubSub(c, ws, r.rdb, subscribe, roomManager, chatRoomID)

	return nil
}
