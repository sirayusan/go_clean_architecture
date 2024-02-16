package v2

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
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

	// ログイン後URL
	e.POST("/home", func(c echo.Context) error {
		ID := c.Get("user_id")
		fmt.Printf("%v \n", ID)
		fmt.Printf("%v \n", ID)
		fmt.Printf("%v \n", ID)
		fmt.Printf("%v \n", ID)
		fmt.Printf("%v \n", ID)
		fmt.Printf("%v \n", ID)
		return c.NoContent(http.StatusOK)
	}, jwtMiddleware())

	e.Logger.Fatal(e.Start(os.ExpandEnv(":${GO_PORT}")))
}

// jwtMiddleware は、JWTトークンを検証する認証ミドルウェアです。
func jwtMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// リクエストヘッダーからトークンを取得
			tokenString := c.Request().Header.Get("Authorization")
			const bearerPrefix = "Bearer "
			if len(tokenString) > len(bearerPrefix) && tokenString[:len(bearerPrefix)] == bearerPrefix {
				tokenString = tokenString[len(bearerPrefix):]
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed JWT")
			}

			// トークンをパースし、検証する
			token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWT token")
			}

			claims, ok := token.Claims.(*jwt.StandardClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT claims22")
			}

			// EchoのContextにユーザーIDを設定
			c.Set("user_id", claims.Id)

			// トークンが有効な場合は次のハンドラーに処理を渡す
			return next(c)
		}
	}
}
