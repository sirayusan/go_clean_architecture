package seeders

import (
	"business/db/model"
	"gorm.io/gorm"
)

// CreateChatMessage　はサンプルデータを投入する。
func CreateChatMessage(tx *gorm.DB) error {
	var err error
	chatMessages := []model.ChatMessage{
		{
			ChatRoomID:   uint32(1),
			Message:      "テスト1",
			SenderUserID: uint32(1),
		},
		{
			ChatRoomID:   uint32(1),
			Message:      "テスト2",
			SenderUserID: uint32(2),
		},
		{
			ChatRoomID:   uint32(2),
			Message:      "テスト3",
			SenderUserID: uint32(1),
		},
		{
			ChatRoomID:   uint32(2),
			Message:      "テスト4",
			SenderUserID: uint32(3),
		},
	}

	for _, chatMessage := range chatMessages {
		err := tx.Create(&chatMessage).Error
		if err != nil {
			return err
		}
	}

	return err
}
