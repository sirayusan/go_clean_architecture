package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"business/internal/entity"
	ct "business/pkg/time"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// MessageUseCase -.
type MessageUseCase struct {
	repo MessageRepo
	wt   ct.WrapperTime
}

// New -.
func New(r MessageRepo, wt ct.WrapperTime) *MessageUseCase {
	return &MessageUseCase{
		repo: r,
		wt:   wt,
	}
}

// JoinRoom は参加ルームを判定して追加し、メッセージ一覧を参加者へ送信する。
func (uc *MessageUseCase) JoinRoom(
	chatRoomID uint32,
	wsw entity.WebSocketWrapper,
	roomManager map[uint32]*entity.ChatRoom,
) error {
	client := &entity.Client{Ws: wsw}

	if _room, exists := roomManager[chatRoomID]; exists {
		_room.AddClient(client)
	} else {
		newRoom := entity.ChatRoom{}
		newRoom.AddClient(client)
		roomManager[chatRoomID] = &newRoom
	}

	messageList, err := uc.repo.GetMessageList(chatRoomID)

	// レコードが存在しないエラーは許容する。
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("GetMessages - s.repo.GetMessageList: %w", err)
	}

	for _, data := range messageList {
		_json, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = client.Send(_json)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *MessageUseCase) PubSub(
	c echo.Context,
	wsw entity.WebSocketWrapper,
	roomManager map[uint32]*entity.ChatRoom,
	chatRoomID uint32,
	rdb entity.RedisWrapper,
) {
	// chatRoomID を文字列に変換
	channelName := "room-" + fmt.Sprint(chatRoomID)

	for {
		_, msg, err := wsw.ReadMessage()
		if err != nil {
			if err == io.EOF {
				// クライアントが接続を閉じた
				c.Logger().Error("WebSocket connection closed by the client")
			} else {
				c.Logger().Error("Error reading from WebSocket: ", err)
			}
			break
		}

		_msg := string(msg)
		if _msg != "" { // msgが空でなければ処理を続ける
			// メッセージをDBに保存し、Redisを介して他の参加者に送信
			if err := uc.processMessage(c, roomManager, rdb, channelName, chatRoomID, _msg); err != nil {
				c.Logger().Error("Error processing message: ", err)
				break
			}
		}
	}

}

// processMessageはメッセージを処理するヘルパー関数です。
func (uc *MessageUseCase) processMessage(
	c echo.Context,
	roomManager map[uint32]*entity.ChatRoom,
	rdb entity.RedisWrapper,
	channelName string,
	chatRoomID uint32,
	_msg string,
) error {
	userIDStr, _ := c.Get("userID").(string) // 自前でセットしているのでエラーハンドリングはしない。
	userID, _ := strconv.Atoi(userIDStr)
	chatMessage, err := uc.repo.CreateMessage(chatRoomID, _msg, uint32(userID))
	if err != nil {
		return err
	}

	_jsonRedis, err := json.Marshal(entity.RedisMessage{
		ServerId:  os.ExpandEnv("${CHANNEL}"), // TODO: kubernetes構築時にPod名を入れる
		Timestamp: time.Time{},
		MessageId: chatMessage.ChatMessageID,
		Payload:   _msg,
	})
	if err != nil {
		return err
	}

	if err := rdb.Publish(c.Request().Context(), channelName, _jsonRedis).Err(); err != nil {
		return err
	}

	// メッセージを指定したチャンネルに送信
	_json, err := json.Marshal(entity.Message{
		UserName:  "",
		Messages:  _msg,
		CreatedAt: uc.wt.Now(),
	})
	if err != nil {
		return err
	}

	roomManager[chatRoomID].Publish(_json)
	return nil
}
