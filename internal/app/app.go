package app

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"os"

	"business/config"
	"business/internal/controller/http/v2"
	"business/pkg/logger"
	"business/pkg/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// DefaultValidator はecho.Validatorインターフェースを実装します
type DefaultValidator struct {
	validator *validator.Validate
}

// Validate　はバリデーションメソッドを定義します。
func (dv *DefaultValidator) Validate(i interface{}) error {
	return dv.validator.Struct(i)
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// フロントのURLをCORS承認する
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.ExpandEnv("${FRONT_DOMAIN}")},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	conn, err := mysql.New()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - mysql.New: %w", err))
	}

	v2.NewRouter(e, conn, l)
}
