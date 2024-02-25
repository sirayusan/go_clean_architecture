package usecase

import (
	"business/internal/entity"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type (
	// Message -.
	Message interface {
		JoinRoom(uint32, entity.WebSocketConn, map[uint32]*entity.ChatRoom) error
		PubSub(echo.Context, *websocket.Conn, map[uint32]*entity.ChatRoom, uint32, *redis.Client)
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
