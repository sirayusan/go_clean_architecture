package usecase

import (
	"business/internal/entity"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type (
	// Message -.
	Message interface {
		GetMessages(uint32) ([]entity.Message, error)
		JoinRoom(uint32, *websocket.Conn, map[uint32]*entity.ChatRoom) error
		PubSub(echo.Context, *websocket.Conn, map[uint32]*entity.ChatRoom, uint32)
	}

	// MessageRepo -.
	MessageRepo interface {
		GetMessageList(uint32) ([]entity.Message, error)
	}
)
