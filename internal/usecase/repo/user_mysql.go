package repo

import (
	"business/internal/entity"
	"business/pkg/mysql"
	"context"
	"errors"
	"log"
)

type Mysql struct {
	conn *mysql.MySQL
}

func New(db *mysql.MySQL) *Mysql {
	return &Mysql{db}
}

func (r *Mysql) Store(ctx context.Context, user entity.User) error {
	// ユーザー情報をデータベースに保存するロジックを実装
	return nil
}

func (r *Mysql) GetUserList(ctx context.Context) ([]entity.User, error) {
	conn := r.conn // MySQL構造体内の*sql.DBにアクセス
	query := "SELECT id, name FROM `user`"
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Internal Server Error. adapter/gateway/GetHistory")
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			log.Println(err)
			return nil, errors.New("Internal Server Error. adapter/gateway/GetHistory")
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, errors.New("Internal Server Error. adapter/gateway/GetHistory")
	}

	return users, nil
}
