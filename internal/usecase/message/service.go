package usecase

import (
	"business/internal/entity"
	"errors"
	"fmt"
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
func (uc *MessageUseCase) GetMessages(chatID uint32) (entity.Messages, error) {
	var chats entity.Messages
	_, err := uc.repo.GetMessageList(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return chats, nil
		}
		return entity.Messages{}, fmt.Errorf("GetMessages - s.repo.GetMessageList: %w", err)
	}

	return chats, nil
}
