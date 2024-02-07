package main

import (
	"business/db/seeder"
	"business/pkg/mysql"
	"fmt"
	"os"
)

func main() {
	// コマンドラインのバリデーション
	err := seeder.CheckArgs()
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

	tx, cleanUP := mysql.Transactional(conn.DB)
	defer cleanUP()

	err = seeder.Seed(tx)
	if err != nil {
		tx.Error = err
		fmt.Printf("データ投入中にエラーが発生しました。\n")
		return
	}

	fmt.Printf("正常に終了しました。\n")
}
