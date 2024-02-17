package v2

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"

	"business/internal/entity"
	user "business/internal/usecase/user"
	"business/pkg/logger"
)

type UserRoutes struct {
	t user.User
	l logger.Interface
}

// NewUserRoutes はUserRoutesのコンストラクタです
func NewUserRoutes(g *echo.Group, t user.User, l logger.Interface) {
	routes := &UserRoutes{t, l}
	g.GET("/index", routes.getUserList)
	g.GET("/home", routes.getChatList)
}

// getUserList はユーザーリストを取得するエンドポイントのハンドラです
func (r *UserRoutes) getUserList(c echo.Context) error {
	userList, err := r.t.UserList()
	if err != nil {
		r.l.Error(err, "http - v2 - getUserList")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err)})
	}

	return c.JSON(http.StatusOK, entity.UserListResponse{UserList: userList})
}

// getUserList はユーザーリストを取得するエンドポイントのハンドラです
func (r *UserRoutes) getChatList(c echo.Context) error {
	userList, err := r.t.UserList()
	if err != nil {
		r.l.Error(err, "http - v2 - getUserList")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err)})
	}

	return c.JSON(http.StatusOK, entity.UserListResponse{UserList: userList})
}
