package usecase

import (
	"business/internal/entity"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// ChatUseCase -.
type ChatUseCase struct {
	repo ChatRepo
}

// New -.
func New(r ChatRepo) *ChatUseCase {
	return &ChatUseCase{
		repo: r,
	}
}

// GetChats はチャット一覧を取得して返す。
func (uc *ChatUseCase) GetChats(userID uint32) (entity.Chats, error) {
	chats, err := uc.repo.GetChatList(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Chats{}, err
		}
		return entity.Chats{}, fmt.Errorf("GetChats - s.repo.GetChatList: %w", err)
	}

	return chats, nil
}
