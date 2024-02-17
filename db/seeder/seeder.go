package seeder

import (
	"errors"
	"gorm.io/gorm"
	"os"

	"business/db/seeder/seeders"
)

// CheckArgs はコマンドライン引数を確認する。
func CheckArgs() error {
	if len(os.Args) != 2 {
		return errors.New("期待している引数は1つです。引数を確認してください。")
	}

	if os.Args[1] != "dev" && os.Args[1] != "test" {
		return errors.New("第一引数が期待している語群は以下の通りです。\n1:dev\n2:test")
	}

	return nil
}

// Seed　はサンプルデータを投入する。
func Seed(tx *gorm.DB) error {
	var err error
	if err = seeders.CreateUser(tx); err != nil {
		return err
	}
	if err = seeders.CreateChat(tx); err != nil {
		return err
	}
	if err = seeders.CreateChatMessage(tx); err != nil {
		return err
	}

	return nil
}
