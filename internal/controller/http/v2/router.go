package v2

import (
	"os"

	"business/internal/usecase"
	authusecase "business/internal/usecase/auth"
	authrepo "business/internal/usecase/auth/repo"
	"business/internal/usecase/repo"
	"business/pkg/logger"
	"business/pkg/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, conn *mysql.MySQL, l logger.Interface) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	UserUseCase := usecase.New(
		repo.New(conn),
	)
	// UserRoutesのインスタンスを作成
	userRouteHandlers := NewUserRoutes(UserUseCase, l)
	u := e.Group("/user")
	u.GET("/index", userRouteHandlers.getUserList)

	authUseCase := authusecase.New(
		authrepo.New(conn),
	)
	NewAuthRouter(e, authUseCase, l)

	e.Logger.Fatal(e.Start(os.ExpandEnv(":${GO_PORT}")))
}
