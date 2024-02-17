package v2

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"time"

	authusecase "business/internal/usecase/auth"
	authrepo "business/internal/usecase/auth/repo"
	user "business/internal/usecase/user"
	"business/internal/usecase/user/repo"
	"business/pkg/logger"
	"business/pkg/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, conn *mysql.MySQL, l logger.Interface) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// UserRoutesのインスタンスを作成
	UserUseCase := user.New(
		repo.New(conn),
	)
	u := e.Group("/user")
	u.Use(jwtMiddleware())
	NewUserRoutes(u, UserUseCase, l)

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

			// トークンの有効期限が現在時刻を過ぎていないか確認
			if claims.ExpiresAt < time.Now().Unix() {
				return echo.NewHTTPError(http.StatusUnauthorized, "JWT token has expired")
			}

			// EchoのContextにユーザーIDを設定
			c.Set("user_id", claims.Id)

			// トークンが有効な場合は次のハンドラーに処理を渡す
			return next(c)
		}
	}
}
