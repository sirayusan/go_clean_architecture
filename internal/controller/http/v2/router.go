package v2

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"

	"business/internal/entity"
	authusecase "business/internal/usecase/auth"
	authrepo "business/internal/usecase/auth/repo"
	chatusecase "business/internal/usecase/chat"
	chatrepo "business/internal/usecase/chat/repo"
	massageusecase "business/internal/usecase/room"
	massagerepo "business/internal/usecase/room/repo"
	user "business/internal/usecase/user"
	"business/internal/usecase/user/repo"
	"business/pkg/logger"
	"business/pkg/mysql"
	ct "business/pkg/time"
)

func NewRouter(e *echo.Echo, conn *mysql.MySQL, rdb *redis.Client, l logger.Interface) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	u := e.Group("/user")
	u.Use(jwtMiddleware())

	NewUserRoutes(u, user.New(repo.New(conn)), l)
	NewAuthRouter(e, authusecase.New(authrepo.New(conn)), l)
	NewChatRouter(e, chatusecase.New(chatrepo.New(conn)), l)

	_ct := ct.CustomTime{}
	wrappedRdb := &entity.RedisConn{Conn: rdb}
	NewMessageRouter(e, massageusecase.New(massagerepo.New(conn, _ct), _ct), l, wrappedRdb)

	// jwt認証URL
	e.GET("/auth", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}, jwtMiddleware())

	if os.Getenv("USE_SSL") == "FALSE" {
		e.Logger.Fatal(e.Start(":" + os.Getenv("GO_PORT")))
	} else {
		certKey := os.Getenv("CERT_PATH") + "/" + os.Getenv("CERT_FILE_NAME")
		secretKey := os.Getenv("SECRET_KEY_PATH") + "/" + os.Getenv("SECRET_KEY_FILE_NAME")

		e.Logger.Fatal(e.StartTLS(":"+os.Getenv("GO_TLS_PORT"), certKey, secretKey))
	}
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
			c.Set("userID", claims.Id)

			// トークンが有効な場合は次のハンドラーに処理を渡す
			return next(c)
		}
	}
}

// websocketJwtMiddleware は、WebSocket用のJWTトークンを検証する認証ミドルウェアです。
func websocketJwtMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			chatIDStr := c.Param("id")
			_, err := strconv.Atoi(chatIDStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "不正なリクエストパラメータ")
			}

			tokenString := c.QueryParam("jwt")
			token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.ExpandEnv("${JWT_SECRET_KEY}")), nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "トークンの解析に失敗しました")
			}

			claims, ok := token.Claims.(*jwt.StandardClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "JWT claimsの取得に失敗しました")
			}

			c.Set("userID", claims.Id)

			return next(c)
		}
	}
}
