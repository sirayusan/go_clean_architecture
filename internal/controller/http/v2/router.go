package v2

import (
	"business/internal/usecase"
	"business/pkg/logger"
	"github.com/labstack/echo/v4"
	"os"
)

func NewRouter(e *echo.Echo, t usecase.User, l logger.Interface) {
	port := os.ExpandEnv(":${GO_PORT}")

	// UserRoutesのインスタンスを作成
	userRouteHandlers := NewUserRoutes(t, l)

	// UserRoutesのメソッドをルートとして登録
	u := e.Group("/user")
	u.GET("/index", userRouteHandlers.GetUserList)

	e.Logger.Fatal(e.Start(port))
}
