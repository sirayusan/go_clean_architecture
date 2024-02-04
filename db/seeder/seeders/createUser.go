package seeders

import (
	"business/db/model"
	"gorm.io/gorm"
)

// CreateUser　はサンプルデータを投入する。
func CreateUser(db *gorm.DB) {
	users := []model.User{
		{Email: "abc@co.jp"},
	}

	for _, user := range users {
		db.Create(&user) // データの投入
	}
}
