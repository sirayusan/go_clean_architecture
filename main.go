package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func main() {
	// Echoインスタンスを作成
	e := echo.New()
	go

	// ルートルートに対するハンドラを登録
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	// サーバーを8080ポートで起動
	if err := e.Start(os.ExpandEnv(":${GO_PORT}")); err != nil {
		e.Logger.Fatal(err)
	}
}
