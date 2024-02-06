package migrations

import (
	"business/db/model"
	"errors"
	"os"
)

// CheckArgs はコマンドライン引数を確認する。
func CheckArgs() error {
	if len(os.Args) != 3 {
		return errors.New("期待している引数は2つです。引数を確認してください。")
	}

	if os.Args[1] != "dev" && os.Args[1] != "test" {
		return errors.New("第一引数が期待している語群は以下の通りです。\n1:dev\n2:test")
	}

	if os.Args[2] != "create" && os.Args[2] != "drop" {
		return errors.New("第二引数が期待している語群は以下の通りです。\n1:create\n2:drop")
	}

	return nil
}

// CreateArrayMigrationSlice はマイグレーション用の構造体が入った配列を返す。
func CreateArrayMigrationSlice() []interface{} {
	return []interface{}{
		model.User{},
	}
}
