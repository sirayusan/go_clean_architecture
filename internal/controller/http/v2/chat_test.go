package v2

import (
	"business/db/model"
	pkgmysql "business/pkg/mysql"
	"encoding/json"
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
	// テスト用のデータベース接続を作成
	conn, err := pkgmysql.NewTest()
	assert.NoError(t, err)
	err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	err = conn.DB.AutoMigrate(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	CreateData(conn, t)

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

	defer func() {
		err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
		assert.NoError(t, err)
		err = conn.DB.AutoMigrate(model.User{}, model.ChatRoom{}, model.ChatMessage{})
		assert.NoError(t, err)
	}()
}

// TestChatRoutes_GetChats_Failed は異常系(204) :チャットが存在しない場合のテスト
func TestChatRoutes_GetChats_Failed(t *testing.T) {
	// テスト用のデータベース接続を作成
	conn, err := pkgmysql.NewTest()
	assert.NoError(t, err)
	err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	err = conn.DB.AutoMigrate(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	CreateData(conn, t)

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

	defer func() {
		err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
		assert.NoError(t, err)
		err = conn.DB.AutoMigrate(model.User{}, model.ChatRoom{}, model.ChatMessage{})
		assert.NoError(t, err)
	}()
}

func CreateData(conn *pkgmysql.MySQL, t *testing.T) {
	var err error
	userList := []model.User{
		{
			UserID:            1,
			LastName:          "高橋",
			FirstName:         "太郎",
			HiraganaLastName:  "たかはし",
			HiraganaFirstName: "たろう",
			Email:             "abc@co.jp",
			Password:          "パスワード",
			CreatedUserID:     0,
			UpdateUserID:      0,
		},
		{
			UserID:            2,
			LastName:          "今井",
			FirstName:         "次郎",
			HiraganaLastName:  "いまい",
			HiraganaFirstName: "たろう",
			Email:             "abcd@co.jp",
			Password:          "パスワード",
			CreatedUserID:     0,
			UpdateUserID:      0,
		},
		{
			UserID:            3,
			LastName:          "斎藤",
			FirstName:         "三郎",
			HiraganaLastName:  "いまい",
			HiraganaFirstName: "たろう",
			Email:             "abcde@co.jp",
			Password:          "パスワード",
			CreatedUserID:     0,
			UpdateUserID:      0,
		},
	}
	for _, user := range userList {
		err = conn.DB.Create(&user).Error
		assert.NoError(t, err)
	}

	chatList := []model.ChatRoom{
		{
			ChatRoomID: 1,
			UserID1:    1,
			UserID2:    2,
		},
		{
			ChatRoomID: 2,
			UserID1:    1,
			UserID2:    3,
		},
	}
	for _, chat := range chatList {
		err = conn.DB.Create(&chat).Error
		assert.NoError(t, err)
	}

	ChatMessageList := []model.ChatMessage{
		{
			ChatRoomID:   1,
			Message:      "テスト1",
			SenderUserID: 1,
		},
		{
			ChatRoomID:   1,
			Message:      "テスト2",
			SenderUserID: 3,
		},
	}
	for _, chatMessage := range ChatMessageList {
		err = conn.DB.Create(&chatMessage).Error
		assert.NoError(t, err)
	}
}
