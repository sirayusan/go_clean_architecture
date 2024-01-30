package gateway

/*
gateway パッケージは，DB操作に対するアダプターです．
*/

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"business/entity"
	"business/usecase/port"
)

type UserRepository struct {
	conn *sql.DB
}

// NewUserRepository はUserRepositoryを返します．
func NewUserRepository(conn *sql.DB) port.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

// GetUserByID はDBからデータを取得します．
func (u *UserRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	conn := u.GetDBConn()
	query := "SELECT id, name FROM `user` WHERE id = ?"
	row := conn.QueryRowContext(ctx, query, userID)
	user := entity.User{}
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User Not Found. UserID = %s", userID)
		}
		log.Println(err)
		return nil, errors.New("Internal Server Error. adapter/gateway/GetUserByID")
	}
	return &user, nil
}

// GetDBConn はconnectionを取得します．
func (u *UserRepository) GetDBConn() *sql.DB {
	return u.conn
}
