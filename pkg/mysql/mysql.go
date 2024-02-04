package mysql

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	DB *gorm.DB
}

// New はGORMを使用してMySQLデータベースに接続するための新しいMySQLインスタンスを生成します。
func New() (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.ExpandEnv("${MYSQL_USER}"), os.ExpandEnv("${MYSQL_PASSWORD}"),
		os.ExpandEnv("${DB_HOST}"), os.ExpandEnv("${DB_PORT}"),
		os.ExpandEnv("${MYSQL_DATABASE}"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
