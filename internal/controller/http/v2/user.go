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

type UserResponse struct {
	UserList []entity.User `json:"history"`
}

// NewUserRoutes はUserRoutesのコンストラクタです
func NewUserRoutes(t usecase.User, l logger.Interface) *UserRoutes {
	return &UserRoutes{t, l}
}

// GetUserList はユーザーリストを取得するエンドポイントのハンドラです
func (r *UserRoutes) GetUserList(c echo.Context) error {
	userList, err := r.t.UserList(c.Request().Context()) // Echoでリクエストからcontext.Contextを取得
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		// Echoでエラーレスポンスを送信
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err)})
	}

	// EchoでJSONレスポンスを送信
	return c.JSON(http.StatusOK, UserResponse{UserList: userList})
}
