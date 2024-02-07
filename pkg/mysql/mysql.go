package mysql

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	DB *gorm.DB
}

// New はGORMを使用してMySQLデータベースに接続するための新しいMySQLインスタンスを生成します。
func New() (*MySQL, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.ExpandEnv("${MYSQL_USER}"), os.ExpandEnv("${MYSQL_PASSWORD}"),
		os.ExpandEnv("${DB_HOST}"), os.ExpandEnv("${DB_PORT}"),
		os.ExpandEnv("${MYSQL_DATABASE}"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &MySQL{DB: db}, nil
}

// NewTest はGORMを使用してMySQLデータベースに接続するための新しいMySQLインスタンスを生成します。
func NewTest() (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.ExpandEnv("${MYSQL_USER}"), os.ExpandEnv("${MYSQL_PASSWORD}"),
		os.ExpandEnv("${DB_HOST}"), os.ExpandEnv("${DB_PORT}"),
		os.ExpandEnv("${MYSQL_TEST_DATABASE}"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &MySQL{DB: db}, nil
}

// Transactional は新しいトランザクションを開始しインスタンスを返す関数です。
// 戻り値の関数はcleanUPとして受け取り、"defer cleanUP()"を直下の行に記述してください。
func Transactional(db *gorm.DB) (*gorm.DB, func()) {
	tx := db.Begin()
	if tx.Error != nil {
		panic("トランザクションの開始に失敗しました。")
	}

	// エラーハンドリングとロールバックを行うクロージャを返す。
	return tx, func() {
		// panicによるエラーの場合
		if r := recover(); r != nil {
			fmt.Printf("予期せぬエラーが発生したため、トランザクションを巻き戻しました。\n")
			tx.Rollback()
		}

		// tx.Errorが設定されている場合（明示的なエラー設定）
		if tx.Error != nil {
			fmt.Printf("エラーが発生したため、トランザクションを巻き戻しました:\n %v\n", tx.Error)
			tx.Rollback()
		}

		tx.Commit()
	}
}
