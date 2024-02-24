package usecase

import (
	"encoding/json"
	"errors"
	"fmt"

	"business/internal/entity"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// MessageUseCase -.
type MessageUseCase struct {
	repo MessageRepo
}

// New -.
func New(r MessageRepo) *MessageUseCase {
	return &MessageUseCase{
		repo: r,
	}
}

// GetMessages はチャット一覧を取得して返す。
func (uc *MessageUseCase) GetMessages(chatID uint32) ([]entity.Message, error) {
	var messageList []entity.Message
	messageList, err := uc.repo.GetMessageList(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Message{}, nil
		}
		return []entity.Message{}, fmt.Errorf("GetMessages - s.repo.GetMessageList: %w", err)
	}

	return messageList, nil
}

// JoinRoom は参加ルームを判定して追加し、メッセージ一覧を参加者へ送信する。
func (uc *MessageUseCase) JoinRoom(
	chatRoomID uint32,
	ws *websocket.Conn,
	roomManager map[uint32]*entity.ChatRoom,
) error {
	client := &entity.Client{Ws: ws}

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

// PubSub はルーム参加者全員に新規メッセージを送信します。
func (uc *MessageUseCase) PubSub(c echo.Context, ws *websocket.Conn, roomManager map[uint32]*entity.ChatRoom, chatRoomID uint32) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			break
		}
		roomManager[chatRoomID].Publish(msg)
	}
}
