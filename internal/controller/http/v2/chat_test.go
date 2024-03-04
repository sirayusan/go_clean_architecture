package v2

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"business/internal/entity"
	"business/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ChatUseCaseMock はChatUsecaseインターフェースのモックです
type ChatUseCaseMock struct {
	mock.Mock
}

func (m *ChatUseCaseMock) GetChats(userID uint32) (entity.ChatRooms, error) {
	args := m.Called(userID)
	return args.Get(0).(entity.ChatRooms), args.Error(1)
}

// TestChatRoutes_GetChats は正常系のテスト
func TestChatRoutes_GetChats(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()

	// モックのusecaseとloggerを作成
	chatUsecaseMock := new(ChatUseCaseMock)
	loggerMock := new(MockLogger)

	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/chats", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	jwtToken, err := auth.GenerateToken(uint32(1))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
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

// TestChatRoutes_GetChats_Failed は異常系(204) :チャットが存在しない場合のテスト
func TestChatRoutes_GetChats_Failed(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()

	// モックのusecaseとloggerを作成
	chatUsecaseMock := new(ChatUseCaseMock)
	loggerMock := new(MockLogger)

	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/chats", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	jwtToken, err := auth.GenerateToken(uint32(1))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	// レスポンスを記録するためのResponseRecorderを作成
	res := httptest.NewRecorder()
	// 新しいEchoコンテキストを生成
	c := e.NewContext(req, res)

	// AuthRoutesのインスタンスを作成
	routes := ChatRoutes{
		t: chatUsecaseMock,
		l: loggerMock,
	}
	chatUsecaseMock.On("GetChats", uint32(0)).Return(entity.ChatRooms{}, gorm.ErrRecordNotFound)
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.GetChats(c)) {
		assert.Equal(t, http.StatusNoContent, res.Code)
		// レスポンスボディを entity.Chats 型にアンマーシャル
		var response map[string]string
		err := json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)

		// アンマーシャルしたレスポンスボディを期待値と比較
		assert.Equal(t, "record not found", response["error"])
	}
	// モックが期待通りに呼び出されたことを確認
	chatUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

// TestChatRoutes_GetChats_Failed は異常系(500) :DB疎通が失敗したパターン
func TestChatRoutes_GetChats_Failed2(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()

	// モックのusecaseとloggerを作成
	chatUsecaseMock := new(ChatUseCaseMock)
	loggerMock := new(MockLogger)

	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/chats", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	jwtToken, err := auth.GenerateToken(uint32(1))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	// レスポンスを記録するためのResponseRecorderを作成
	res := httptest.NewRecorder()
	// 新しいEchoコンテキストを生成
	c := e.NewContext(req, res)

	// AuthRoutesのインスタンスを作成
	routes := ChatRoutes{
		t: chatUsecaseMock,
		l: loggerMock,
	}
	chatUsecaseMock.On("GetChats", uint32(0)).Return(entity.ChatRooms{}, errors.New("DB疎通によるエラーメッセージ"))
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.GetChats(c)) {
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		// レスポンスボディを entity.Chats 型にアンマーシャル
		var response map[string]string
		err := json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)

		// アンマーシャルしたレスポンスボディを期待値と比較
		assert.Equal(t, "DB疎通によるエラーメッセージ", response["error"])
	}
	// モックが期待通りに呼び出されたことを確認
	chatUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

// TestNewChatRouter NewChatRouter呼び出し
// 実施する意味ないが見栄えの為。
func TestNewChatRouter(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()

	// モックのusecaseとloggerを作成
	chatUsecaseMock := new(ChatUseCaseMock)
	loggerMock := new(MockLogger)

	NewChatRouter(e, chatUsecaseMock, loggerMock)

	// モックが期待通りに呼び出されたことを確認
	chatUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}
