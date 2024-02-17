package v2

import (
	"net/http"
	"strconv"

	"business/internal/entity"
	"business/internal/usecase/chat"
	"business/pkg/logger"
	"github.com/labstack/echo/v4"
)

type ChatRoutes struct {
	t usecase.Chat
	l logger.Interface
}

// NewChatRouter はチャット関連のURLからコントローラーを実行します。
func NewChatRouter(e *echo.Echo, t usecase.Chat, l logger.Interface) {
	routes := &ChatRoutes{t, l}
	e.GET("/chats", func(c echo.Context) error {
		return routes.GetChats(c)
	}, jwtMiddleware())
}

// GetChats はユーザーリストを取得するエンドポイントのハンドラです
func (r *ChatRoutes) GetChats(c echo.Context) error {
	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		// userIDが文字列ではない場合
		return echo.NewHTTPError(http.StatusBadRequest, "userID must be a string")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		// 文字列から整数への変換に失敗した場合のエラー処理
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid userID")
	}

	chats, err := r.t.GetChats(uint32(userID))

	return c.JSON(http.StatusOK, entity.Chats{List: chats})
}
