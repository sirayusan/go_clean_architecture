package v2

import (
	"business/db/model"
	pkgmysql "business/pkg/mysql"
	"net/http"
	"net/http/httptest"
	"testing"

	"business/internal/entity"
	"business/pkg/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ChatUseCaseMock はChatUsecaseインターフェースのモックです
type ChatUseCaseMock struct {
	mock.Mock
}

func (m *ChatUseCaseMock) GetChats(userID uint32) (entity.Chats, error) {
	args := m.Called(userID)
	return args.Get(0).(entity.Chats), args.Error(1)
}

func TestChatRoutes_GetChats(t *testing.T) {
	// テスト用のデータベース接続を作成
	conn, err := pkgmysql.NewTest()
	assert.NoError(t, err)
	err = conn.DB.Migrator().DropTable(model.User{}, model.Chat{}, model.ChatMessage{})
	assert.NoError(t, err)
	err = conn.DB.AutoMigrate(model.User{}, model.Chat{}, model.ChatMessage{})
	assert.NoError(t, err)
	CreateData(conn, t)

	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

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
	chatUsecaseMock.On("GetChats", nil).Return(entity.Chats{}, nil)
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.GetChats(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
		// レスポンスボディを検証する
		assert.Equal(t, res.Body, nil)
	}
	// モックが期待通りに呼び出されたことを確認
	chatUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

	// 異常系(204) :チャットが存在しない場合
	jwtToken, err = auth.GenerateToken(uint32(3))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)
	// テスト対象のメソッドを実行
	chatUsecaseMock.On("GetChats", nil).Return(entity.Chats{}, nil) // ここをuint32(3)に修正
	if assert.NoError(t, routes.GetChats(c)) {
		assert.Equal(t, http.StatusNoContent, res.Code)
		// レスポンスボディを検証する
		assert.Equal(t, "", res.Body.String())
	}
	// モックが期待通りに呼び出されたことを確認
	chatUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
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

	chatList := []model.Chat{
		{
			ChatID:  1,
			UserID1: 1,
			UserID2: 2,
		},
		{
			ChatID:  2,
			UserID1: 1,
			UserID2: 3,
		},
	}
	for _, chat := range chatList {
		err = conn.DB.Create(&chat).Error
		assert.NoError(t, err)
	}

	ChatMessageList := []model.ChatMessage{
		{
			ChatID:       1,
			Message:      "テスト1",
			SenderUserID: 1,
		},
		{
			ChatID:       1,
			Message:      "テスト2",
			SenderUserID: 3,
		},
	}
	for _, chatMessage := range ChatMessageList {
		err = conn.DB.Create(&chatMessage).Error
		assert.NoError(t, err)
	}
}
