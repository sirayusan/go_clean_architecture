package app

import (
	"business/config"
	"business/internal/controller/http/v2"
	"business/internal/usecase"
	"business/internal/usecase/repo"
	"business/pkg/logger"
	"business/pkg/mysql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	conn, err := mysql.New() // MySQL設定を使用
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - mysql.New: %w", err))
	}
	defer conn.Close()

	UserUseCase := usecase.New(
		repo.New(conn), // MySQL接続インスタンス
	)

	v2.NewRouter(e, UserUseCase, l)
}
