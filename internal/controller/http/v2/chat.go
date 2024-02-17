package v2

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"

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
	userIDStr, _ := c.Get("userID").(string)
	userID, _ := strconv.Atoi(userIDStr)

	chats, err := r.t.GetChats(uint32(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNoContent, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, chats)
}
