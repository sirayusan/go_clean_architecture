package v2

import (
	"business/internal/entity"
	"business/internal/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"

	"business/pkg/logger"
)

type UserRoutes struct {
	t usecase.User
	l logger.Interface
}

type doTranslateRequest struct {
	Source      string `json:"source"       binding:"required"  example:"auto"`
	Destination string `json:"destination"  binding:"required"  example:"en"`
	Original    string `json:"original"     binding:"required"  example:"текст для перевода"`
}

// NewUserRoutes はUserRoutesのコンストラクタです
func NewUserRoutes(t usecase.User, l logger.Interface) *UserRoutes {
	return &UserRoutes{t, l}
}

// GetUserList はユーザーリストを取得するエンドポイントのハンドラです
func (r *UserRoutes) getUserList(c echo.Context) error {
	userList, err := r.t.UserList(c.Request().Context()) // Echoでリクエストからcontext.Contextを取得
	if err != nil {
		r.l.Error(err, "http - v2 - getUserList")
		// Echoでエラーレスポンスを送信
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err)})
	}

	// EchoでJSONレスポンスを送信
	return c.JSON(http.StatusOK, entity.UserListResponse{UserList: userList})
}
