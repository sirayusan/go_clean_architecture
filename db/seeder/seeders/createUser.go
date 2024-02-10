package seeders

import (
	"business/db/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser　はサンプルデータを投入する。
func CreateUser(tx *gorm.DB) error {
	var err error
	users := []model.User{
		{
			LastName:          "高橋",
			FirstName:         "太郎",
			HiraganaLastName:  "たかはし",
			HiraganaFirstName: "たろう",
			Email:             "abc@co.jp",
			Password:          encrypt("パスワード"),
			CreatedUserID:     0,
			UpdateUserID:      0,
		},
		{
			LastName:          "今井",
			FirstName:         "次郎",
			HiraganaLastName:  "いまい",
			HiraganaFirstName: "たろう",
			Email:             "abcd@co.jp",
			Password:          encrypt("パスワード"),
			CreatedUserID:     0,
			UpdateUserID:      0,
		},
	}

	for _, user := range users {
		err := tx.Create(&user).Error
		if err != nil {
			return err
		}
	}

	return err
}

func encrypt(plainText string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
