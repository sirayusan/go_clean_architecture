package usecase

import (
	"business/internal/entity"
	"github.com/labstack/echo/v4"
)

type (
	// Message -.
	Message interface {
		JoinRoom(uint32, entity.WebSocketWrapper, map[uint32]*entity.ChatRoom) error
		PubSub(echo.Context, entity.WebSocketWrapper, entity.RedisWrapper, *entity.PubSub, map[uint32]*entity.ChatRoom, uint32)
	}

	WebSocketConnInterface interface {
		WriteMessage(messageType int, data []byte) error
	}

	// MessageRepo -.
	MessageRepo interface {
		GetMessageList(uint32) ([]entity.Message, error)
		CreateMessage(uint32, string, uint32) (entity.ChatMessage, error)
	}
)
