package usecase

import (
	"business/internal/entity"
)

type (
	// Chat -.
	Chat interface {
		GetChats(uint32) (entity.ChatRooms, error)
	}

	// ChatRepo -.
	ChatRepo interface {
		GetChatList(uint32) (entity.ChatRooms, error)
	}
)
