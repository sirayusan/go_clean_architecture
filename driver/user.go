package driver

/*
driver パッケージは，技術的な実装を持ちます．．
*/

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"

	// blank import for MySQL driver
	"business/adapter/controller"
	"business/adapter/gateway"
	"business/adapter/presenter"
	"business/usecase/interactor"
	_ "github.com/go-sql-driver/mysql"
)

// Serve はserverを起動させます．
// e はechoのインスタンス
// port は":8080"の形式で渡される。
func Serve(e *echo.Echo, port string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DATABASE"))
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return
	}

	// ルートルートに対するハンドラを登録
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	user := controller.User{
		OutputFactory: presenter.NewUserOutputPort,
		InputFactory:  interactor.NewUserInputPort,
		RepoFactory:   gateway.NewUserRepository,
		Conn:          conn,
	}
	http.HandleFunc("/user/", user.GetUserByID)
}
