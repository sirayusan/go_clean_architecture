package v2

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"business/internal/entity"
	"business/internal/usecase/room"
	"business/pkg/logger"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type MessageRoutes struct {
	t   usecase.Message
	l   logger.Interface
	rdb *redis.Client
}

// NewMessageRouter はチャット関連のURLからコントローラーを実行します。
func NewMessageRouter(e *echo.Echo, t usecase.Message, l logger.Interface, r *redis.Client) {
	routes := &MessageRoutes{t, l, r}
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
		return err
	}

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

	wsWrapper := &entity.WebSocketWrapper{Conn: ws}

	// リクエストをルームに参加させる。
	err = r.t.JoinRoom(chatRoomID, wsWrapper, roomManager)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	subscribe := r.rdb.Subscribe(c.Request().Context(), "room-"+fmt.Sprint(chatRoomID))

	defer func() {
		if err := subscribe.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}()

	// ルーム参加者からのメッセージを検知し送信する。
	r.t.PubSub(c, ws, roomManager, chatRoomID, r.rdb)
	currentServerID := os.ExpandEnv("${CHANNEL}") // 現在のサーバーIDを取得

	for {
		msg, err := subscribe.ReceiveMessage(c.Request().Context())
		if err != nil {
			// エラーハンドリング
			break
		}

		var receivedMessage entity.RedisMessage
		err = json.Unmarshal([]byte(msg.Payload), &receivedMessage)
		if err != nil {
			// JSONのアンマーシャル中にエラーが発生した場合のエラーハンドリング
			break
		}

		// 現在のサーバーIDと受信したメッセージのサーバーIDが異なる場合に処理を実行
		if receivedMessage.ServerId != currentServerID {
			_json, err := json.Marshal(entity.Message{
				SenderUserID: 1,
				UserName:     "",
				Messages:     receivedMessage.Payload,
				CreatedAt:    time.Now(),
			})
			if err != nil {
				break
			}

			roomManager[chatRoomID].Publish(_json)
		}
	}

	return nil
}

// validate はリクエストパラメータのバリデーションを行う。
func validate(c echo.Context) (uint32, error) {
	chatIDStr := c.Param("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "不正なリクエストパラメータ")
	}

	return uint32(chatID), nil
}
