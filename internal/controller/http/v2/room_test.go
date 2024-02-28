package v2

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"business/internal/entity"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestChatRoutes_GetChats は正常系のテスト
func TestNewMessageRouter(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	chatUsecaseMock := new(ChatUseCaseMock)
	loggerMock := new(MockLogger)

	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/chats/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//jwtToken, err := auth.GenerateToken(uint32(1))
	//assert.NoError(t, err)
	//req.Header.Set("Authorization", "Bearer "+jwtToken)
	// レスポンスを記録するためのResponseRecorderを作成
	res := httptest.NewRecorder()
	// 新しいEchoコンテキストを生成
	c := e.NewContext(req, res)

	// AuthRoutesのインスタンスを作成
	routes := ChatRoutes{
		t: chatUsecaseMock,
		l: loggerMock,
	}
	assertChatList := entity.ChatRooms{
		List: []entity.Room{
			{
				ChatRoomID: 1,
				UserName:   "今井次郎",
				Message:    func() *string { s := "テスト1"; return &s }(),
			},
			{
				ChatRoomID: 2,
				UserName:   "斎藤三郎",
				Message:    func() *string { s := "テスト2"; return &s }(),
			},
		},
	}

	chatUsecaseMock.On("GetChats", uint32(0)).Return(assertChatList, nil)
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.GetChats(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
		// レスポンスボディを entity.Chats 型にアンマーシャル
		var actualChatList entity.ChatRooms
		err := json.Unmarshal(res.Body.Bytes(), &actualChatList)
		assert.NoError(t, err)

		// アンマーシャルしたレスポンスボディを期待値と比較
		assert.Equal(t, assertChatList, actualChatList)
	}
	// モックが期待通りに呼び出されたことを確認
	chatUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

}
