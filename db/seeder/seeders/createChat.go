package seeders

import (
	"business/db/model"
	"gorm.io/gorm"
)

// CreateChat　はサンプルデータを投入する。
func CreateChat(tx *gorm.DB) error {
	var err error
	chats := []model.Chat{
		{
			ChatID:  uint32(1),
			UserID1: 1,
			UserID2: 2,
		},
		{
			ChatID:  uint32(2),
			UserID1: 1,
			UserID2: 3,
		},
	}

	for _, chat := range chats {
		err := tx.Create(&chat).Error
		if err != nil {
			return err
		}
	}

	return err
}
