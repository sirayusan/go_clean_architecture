package main

import (
	"business/driver"
	"github.com/labstack/echo/v4"
)

func main() {
	// port　は環境変数GO_PORTを取得した値。
	//port := os.ExpandEnv(":${GO_PORT}")
	port := ":8080"

	// Echoインスタンスを作成
	e := echo.New()

	driver.Serve(e, port)

	if err := e.Start(port); err != nil {
		e.Logger.Fatal(err)
	}
}
