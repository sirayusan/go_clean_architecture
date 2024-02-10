package v2

import (
	"business/internal/entity"
	"business/internal/usecase/auth"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"

	"business/pkg/logger"
)

type AuthRoutes struct {
	t usecase.Auth
	l logger.Interface
}

// NewAuthRouter はAuthRoutesのインスタンスを作成します
func NewAuthRouter(e *echo.Echo, t usecase.Auth, l logger.Interface) {
	routes := &AuthRoutes{t, l}
	e.POST("/login", func(c echo.Context) error {
		return routes.Authentication(c)
	})
}

// Authentication はユーザーリストを取得するエンドポイントのハンドラです
func (r *AuthRoutes) Authentication(c echo.Context) error {
	// 引数の値のフォーマットチェック
	// リクエストパラメータ取得
	var param entity.LoginRequest

	if err := c.Bind(&param); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(param); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	jwtToken, err := r.t.GenerateJwtToken(param)
	if err != nil {
		// ユーザーが存在しない場合もパスワードが一致しない場合も同じエラーを返す。
		errMsg := "ユーザー認証に失敗しました。"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.l.Error(err, errMsg)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": errMsg})
		}
		if err.Error() == "invalid password" {
			r.l.Error(err, errMsg)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": errMsg})
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"jwt": jwtToken})
}
