package main

import (
	"business/db/migrations"
	"fmt"
	"os"

	"business/pkg/mysql"
)

// main は引数からテーブル作成を行います
// 引数:
// - arg1: 接続環境の指定。期待する語群:dev or test
// - arg1: テーブルの作成か、削除の指定 期待する語群:create or drop
func main() {
	// コマンドラインのバリデーション
	err := migrations.CheckArgs()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	var conn *mysql.MySQL
	if os.Args[1] == "dev" {
		conn, err = mysql.New()
	} else if os.Args[1] == "test" {
		conn, err = mysql.NewTest()
	}
	if err != nil {
		panic(err)
	}

	// connがnilでないことを確認
	if conn == nil || conn.DB == nil {
		panic("データベース接続が初期化されていません。")
	}

	if os.Args[2] == "create" {
		err = conn.DB.AutoMigrate(migrations.CreateArrayMigrationSlice()...)
	} else if os.Args[2] == "drop" {
		err = conn.DB.Migrator().DropTable(migrations.CreateArrayMigrationSlice()...)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("正常に終了しました。\n")
}
