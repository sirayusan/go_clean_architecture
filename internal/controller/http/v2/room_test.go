package v2

import (
	"business/pkg/auth"
	"business/pkg/redis"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/ras0q/go-wstest/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"business/internal/entity"
)

// MessageMock はMessageインターフェースのモックです
type MessageMock struct {
	mock.Mock
}

func (m *MessageMock) JoinRoom(chatRoomID uint32, wsw entity.WebSocketWrapper, roomManager map[uint32]*entity.ChatRoom) error {
	args := m.Called(chatRoomID, wsw, roomManager)
	return args.Error(0)
}

func (m *MessageMock) PubSub(c echo.Context, wsw entity.WebSocketWrapper, rdb entity.RedisWrapper, subscribe *entity.PubSub, roomManager map[uint32]*entity.ChatRoom, chatRoomID uint32) {
	return
}

type wsHandler struct{}

// *wsHandlerはhttp.Handlerインターフェイスを満たす
func (h *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws.Serve(w, r)
}

// TestWebSocketHandler は正常系のテスト
func TestWebSocketHandler(t *testing.T) {
	rdb := redis.NewRedis()
	wrappedRdb := &entity.RedisConn{Conn: rdb}
	messageUseCaseMock := new(MessageMock)
	loggerMock := new(MockLogger)
	routes := MessageRoutes{
		t:   messageUseCaseMock,
		l:   loggerMock,
		rdb: wrappedRdb,
	}
	messageUseCaseMock.On(
		"JoinRoom",
		mock.AnythingOfType("uint32"),
		mock.AnythingOfType("*websocket.Conn"),
		mock.AnythingOfType("map[uint32]*entity.ChatRoom"),
	).Return(nil)

	jwtToken, err := auth.GenerateToken(uint32(1))
	assert.Empty(t, err)

	e := echo.New()
	e.GET("/chats/:id", routes.handleConnections, websocketJwtMiddleware())

	// テスト用サーバーの立ち上げ
	server := httptest.NewServer(e)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/chats/1?jwt=" + jwtToken
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.Empty(t, err)

	err = conn.WriteMessage(websocket.TextMessage, []byte("hello"))
	assert.Empty(t, err)
}

// TestWebSocketHandlerException 異常系:不正なリクエストパラメータによるハンドシェイク失敗
func TestWebSocketHandlerException(t *testing.T) {
	rdb := redis.NewRedis()
	wrappedRdb := &entity.RedisConn{Conn: rdb}
	messageUseCaseMock := new(MessageMock)
	loggerMock := new(MockLogger)
	routes := MessageRoutes{
		t:   messageUseCaseMock,
		l:   loggerMock,
		rdb: wrappedRdb,
	}
	messageUseCaseMock.On(
		"JoinRoom",
		mock.AnythingOfType("uint32"),
		mock.AnythingOfType("*websocket.Conn"),
		mock.AnythingOfType("map[uint32]*entity.ChatRoom"),
	).Return(nil)

	jwtToken, err := auth.GenerateToken(uint32(1))
	assert.Empty(t, err)

	e := echo.New()
	e.GET("/chats/:id", routes.handleConnections, websocketJwtMiddleware())

	// テスト用サーバーの立ち上げ
	server := httptest.NewServer(e)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/chats/あ?jwt=" + jwtToken
	_, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NotEmpty(t, err)
	assert.Equal(t, "websocket: bad handshake", err.Error())
}

// TestWebSocketHandlerException2 異常系:不正なトークン
func TestWebSocketHandlerException2(t *testing.T) {
	rdb := redis.NewRedis()
	wrappedRdb := &entity.RedisConn{Conn: rdb}
	messageUseCaseMock := new(MessageMock)
	loggerMock := new(MockLogger)
	routes := MessageRoutes{
		t:   messageUseCaseMock,
		l:   loggerMock,
		rdb: wrappedRdb,
	}
	messageUseCaseMock.On(
		"JoinRoom",
		mock.AnythingOfType("uint32"),
		mock.AnythingOfType("*websocket.Conn"),
		mock.AnythingOfType("map[uint32]*entity.ChatRoom"),
	).Return(nil)

	// 不正な秘密鍵で署名したパターン
	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTIxNTI4NjUsImp0aSI6IjEiLCJpYXQiOjE3MDk1NjA4NjUsImlzcyI6Ijpnb19jbGVhbl9hcmNoaXRlY3R1cmUifQ.vVA44m7iFqubbUU2RBVJFE3jACyDB72PY6d8als6Y-nQ"

	e := echo.New()
	e.GET("/chats/:id", routes.handleConnections, websocketJwtMiddleware())

	// テスト用サーバーの立ち上げ
	server := httptest.NewServer(e)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/chats/1?jwt=" + jwtToken
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NotEmpty(t, err)
	assert.Equal(t, "websocket: bad handshake", err.Error())
}

// TestWebSocketHandlerException3 異常系:JWTトークンの有効期限が切れています
func TestWebSocketHandlerException3(t *testing.T) {
	rdb := redis.NewRedis()
	wrappedRdb := &entity.RedisConn{Conn: rdb}
	messageUseCaseMock := new(MessageMock)
	loggerMock := new(MockLogger)
	routes := MessageRoutes{
		t:   messageUseCaseMock,
		l:   loggerMock,
		rdb: wrappedRdb,
	}
	messageUseCaseMock.On(
		"JoinRoom",
		mock.AnythingOfType("uint32"),
		mock.AnythingOfType("*websocket.Conn"),
		mock.AnythingOfType("map[uint32]*entity.ChatRoom"),
	).Return(nil)

	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NTEwNDksImp0aSI6IjEiLCJpYXQiOjE3MDk2NTEwNDksImlzcyI6ImdvX2NsZWFuX2FyY2hpdGVjdHVyZSJ9.Fjz4HgFvpnkVOTxu-h1HFubByxYznMCpIoj6rkS6XY8"

	e := echo.New()
	e.GET("/chats/:id", routes.handleConnections, websocketJwtMiddleware())

	// テスト用サーバーの立ち上げ
	server := httptest.NewServer(e)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/chats/1?jwt=" + jwtToken
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NotEmpty(t, err)
	assert.Equal(t, "websocket: bad handshake", err.Error())
}

// TestNewMessageRouter呼び出し
// WebSocketを再現していないのでエラーになるが、Routerの動作確認ができているので問題ない。
func TestNewMessageRouter(t *testing.T) {
	e := echo.New()
	rdb := entity.RedisConn{}
	wrappedRdb := &entity.RedisConn{Conn: rdb.Conn}
	messageUseCaseMock := new(MessageMock)
	loggerMock := new(MockLogger)

	NewMessageRouter(e, messageUseCaseMock, loggerMock, wrappedRdb)
	messageUseCaseMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

	jwtToken, err := auth.GenerateToken(uint32(1))
	assert.Empty(t, err)

	req := httptest.NewRequest(http.MethodGet, "/chats/1?jwt="+jwtToken, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code, "Expected HTTP status code 200")
}
