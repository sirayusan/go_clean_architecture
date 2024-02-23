package usecase

import (
	"business/internal/entity"
)

type (
	// Message -.
	Message interface {
		GetMessages(uint32) (entity.Messages, error)
	}

	// MessageRepo -.
	MessageRepo interface {
		GetMessageList(uint32) (entity.Messages, error)
	}
)
