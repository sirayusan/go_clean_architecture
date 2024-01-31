package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type MySQL struct {
	*sql.DB
}

// New Config引数を削除
func New() (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.ExpandEnv("${MYSQL_USER}"), os.ExpandEnv("${MYSQL_PASSWORD}"), os.ExpandEnv("${DB_HOST}"), os.ExpandEnv("${DB_PORT}"), os.ExpandEnv("${MYSQL_DATABASE}"))
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &MySQL{conn}, nil
}
