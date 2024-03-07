package v2

import (
	"bytes"
	"errors"
	"testing"

	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"

	"business/internal/entity"
)

// DefaultValidator はecho.Validatorインターフェースを実装します
type DefaultValidator struct {
	validator *validator.Validate
}

// Validate　はバリデーションメソッドを定義します。
func (dv *DefaultValidator) Validate(i interface{}) error {
	return dv.validator.Struct(i)
}

// AuthUsecaseMock はAuthUsecaseインターフェースのモックです
type AuthUsecaseMock struct {
	mock.Mock
}

func (m *AuthUsecaseMock) GenerateJwtToken(email entity.LoginRequest) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

type MockLogger struct {
	logger mock.Mock
	mock.Mock
}

func (m *MockLogger) Debug(message interface{}, args ...interface{}) {
}

func (m *MockLogger) Info(message string, args ...interface{}) {
}

func (m *MockLogger) Warn(message string, args ...interface{}) {
}

// Error -.
func (m *MockLogger) Error(message interface{}, args ...interface{}) {
}

func (m *MockLogger) Fatal(message interface{}, args ...interface{}) {
}

func (m *MockLogger) log(message string, args ...interface{}) {
}

func (m *MockLogger) msg(level string, message interface{}, args ...interface{}) {
}

// TestAuthRoutes_Authentication 正常系
func TestAuthRoutes_Authentication(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	// モックのusecaseとloggerを作成
	authUsecaseMock := new(AuthUsecaseMock)
	loggerMock := new(MockLogger)

	// テスト用のリクエストボディを作成
	loginRequest := entity.LoginRequest{
		Mail:     "test@example.com",
		Password: "password",
	}
	// jsonへ変換
	requestBody, _ := json.Marshal(loginRequest)
	// テスト用のHTTPリクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// レスポンスを記録するためのResponseRecorderを作成
	res := httptest.NewRecorder()
	// 新しいEchoコンテキストを生成
	c := e.NewContext(req, res)

	// AuthRoutesのインスタンスを作成
	routes := AuthRoutes{
		t: authUsecaseMock,
		l: loggerMock,
	}

	authUsecaseMock.On("GenerateJwtToken", loginRequest).Return("mockedJwtToken", nil)

	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.Authentication(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
		// レスポンスボディを検証する
		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response)
		assert.Equal(t, "mockedJwtToken", response["jwt"])
	}

	// モックが期待通りに呼び出されたことを確認
	authUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

// TestAuthRoutes_Authentication2 異常系: メールアドレスが入力されていない場合。
func TestAuthRoutes_Authentication2(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	// モックのusecaseとloggerを作成
	authUsecaseMock := new(AuthUsecaseMock)
	loggerMock := new(MockLogger)

	// AuthRoutesのインスタンスを作成
	routes := AuthRoutes{
		t: authUsecaseMock,
		l: loggerMock,
	}

	// モックが期待通りに呼び出されたことを確認
	authUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

	//異常系: メールアドレスが入力されていない場合。
	//テスト用のリクエストボディを作成
	loginRequest := entity.LoginRequest{
		Mail:     "",
		Password: "password",
	}
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	authUsecaseMock.On("GenerateJwtToken", loginRequest).Return("", nil)
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.Authentication(c)) {
		assert.Equal(t, http.StatusBadRequest, res.Code)
		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response) // レスポンスボディを検証する
		assert.Equal(t, "Key: 'LoginRequest.Mail' Error:Field validation for 'Mail' failed on the 'required' tag", response["error"])
	}

}

// TestAuthRoutes_Authentication2 異常系: パスワードが入力されていない場合。
func TestAuthRoutes_Authentication3(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	// モックのusecaseとloggerを作成
	authUsecaseMock := new(AuthUsecaseMock)
	loggerMock := new(MockLogger)

	// AuthRoutesのインスタンスを作成
	routes := AuthRoutes{
		t: authUsecaseMock,
		l: loggerMock,
	}

	// モックが期待通りに呼び出されたことを確認
	authUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

	//異常系: パスワードが入力されていない場合。
	//テスト用のリクエストボディを作成
	loginRequest := entity.LoginRequest{
		Mail:     "test@example.com",
		Password: "",
	}
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	authUsecaseMock.On("GenerateJwtToken", loginRequest).Return("", nil)
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.Authentication(c)) {
		assert.Equal(t, http.StatusBadRequest, res.Code)
		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response) // レスポンスボディを検証する
		assert.Equal(t, "Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag", response["error"])
	}

}

// TestAuthRoutes_Authentication2 異常系: 不正なJSON項目
func TestAuthRoutes_Authentication4(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	// モックのusecaseとloggerを作成
	authUsecaseMock := new(AuthUsecaseMock)
	loggerMock := new(MockLogger)

	// AuthRoutesのインスタンスを作成
	routes := AuthRoutes{
		t: authUsecaseMock,
		l: loggerMock,
	}

	// モックが期待通りに呼び出されたことを確認
	authUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

	// 異常系: 不正なJSON項目
	loginRequest_ := `
    aaa
    `
	requestBody, _ := json.Marshal(loginRequest_)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	authUsecaseMock.On("GenerateJwtToken", loginRequest_).Return("", gorm.ErrRecordNotFound)
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.Authentication(c)) {
		assert.Equal(t, http.StatusBadRequest, res.Code)
		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response) // レスポンスボディを検証する
		assert.Equal(t, "", response["error"])
	}

}

// TestAuthRoutes_Authentication2 異常系: メールアドレスのユーザーが存在しない場合
func TestAuthRoutes_Authentication5(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	// モックのusecaseとloggerを作成
	authUsecaseMock := new(AuthUsecaseMock)
	loggerMock := new(MockLogger)

	// AuthRoutesのインスタンスを作成
	routes := AuthRoutes{
		t: authUsecaseMock,
		l: loggerMock,
	}

	// モックが期待通りに呼び出されたことを確認
	authUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

	// 異常系: メールアドレスのユーザーが存在しない場合
	loginRequest := entity.LoginRequest{
		Mail:     "test@example.com12345",
		Password: "パスワード",
	}
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	authUsecaseMock.On("GenerateJwtToken", loginRequest).Return("", gorm.ErrRecordNotFound)
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.Authentication(c)) {
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response) // レスポンスボディを検証する
		assert.Equal(t, "ユーザー認証に失敗しました。", response["error"])
	}

}

// TestAuthRoutes_Authentication2 異常系: パスワードが一致しない場合
func TestAuthRoutes_Authentication6(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	// モックのusecaseとloggerを作成
	authUsecaseMock := new(AuthUsecaseMock)
	loggerMock := new(MockLogger)

	// AuthRoutesのインスタンスを作成
	routes := AuthRoutes{
		t: authUsecaseMock,
		l: loggerMock,
	}

	// モックが期待通りに呼び出されたことを確認
	authUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

	// 異常系: パスワードが一致しない場合
	loginRequest := entity.LoginRequest{
		Mail:     "test@example.com",
		Password: "パスワード",
	}
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	authUsecaseMock.On("GenerateJwtToken", loginRequest).Return("", errors.New("invalid password"))
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.Authentication(c)) {
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response) // レスポンスボディを検証する
		assert.Equal(t, "ユーザー認証に失敗しました。", response["error"])
	}

}

// TestAuthRoutes_Authentication2 異常系: DBとの疎通時に何らかの原因でエラーになった場合
func TestAuthRoutes_Authentication7(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	// モックのusecaseとloggerを作成
	authUsecaseMock := new(AuthUsecaseMock)
	loggerMock := new(MockLogger)

	// AuthRoutesのインスタンスを作成
	routes := AuthRoutes{
		t: authUsecaseMock,
		l: loggerMock,
	}

	// モックが期待通りに呼び出されたことを確認
	authUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

	// 異常系: DBとの疎通時に何らかの原因でエラーになった場合
	loginRequest := entity.LoginRequest{
		Mail:     "test@example.com",
		Password: "存在しないパスワード",
	}
	requestBody, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	authUsecaseMock.On("GenerateJwtToken", loginRequest).Return("", errors.New("DB疎通時の想定外のエラー"))
	// テスト対象のメソッドを実行
	if assert.NoError(t, routes.Authentication(c)) {
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response) // レスポンスボディを検証する
		assert.Equal(t, "DB疎通時の想定外のエラー", response["error"])
	}

}

// NewAuthRouter呼び出し
// 実施する意味ないが見栄えの為。
func TestNewAuthRouter(t *testing.T) {
	// Echoのインスタンスを生成
	e := echo.New()
	// カスタムバリデーターをechoインスタンスに登録
	e.Validator = &DefaultValidator{validator: validator.New()}

	// モックのusecaseとloggerを作成
	authUsecaseMock := new(AuthUsecaseMock)
	loggerMock := new(MockLogger)

	NewAuthRouter(e, authUsecaseMock, loggerMock)

	// モックが期待通りに呼び出されたことを確認
	authUsecaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}
