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

func (m *MessageMock) PubSub(c echo.Context, wsw entity.WebSocketWrapper, roomManager map[uint32]*entity.ChatRoom, chatRoomID uint32, rdb entity.RedisWrapper) {
	return
}

func (m *MessageMock) RedisPubSub(c echo.Context, subscribe *entity.PubSub, roomManager map[uint32]*entity.ChatRoom, chatRoomID uint32) {
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
	e.GET("/chats/:id", routes.handleConnections)

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
	e.GET("/chats/:id", routes.handleConnections)

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
	e.GET("/chats/:id", routes.handleConnections)

	// テスト用サーバーの立ち上げ
	server := httptest.NewServer(e)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/chats/1?jwt=" + jwtToken
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NotEmpty(t, err)
	assert.Equal(t, "websocket: bad handshake", err.Error())
}

// TestWebSocketHandlerException3 異常系:JWTトークンの有効期限が切れています
//func TestWebSocketHandlerException3(t *testing.T) {
//	rdb := redis.NewRedis()
//	wrappedRdb := &entity.RedisConn{Conn: rdb}
//	messageUseCaseMock := new(MessageMock)
//	loggerMock := new(MockLogger)
//	routes := MessageRoutes{
//		t:   messageUseCaseMock,
//		l:   loggerMock,
//		rdb: wrappedRdb,
//	}
//	messageUseCaseMock.On(
//		"JoinRoom",
//		mock.AnythingOfType("uint32"),
//		mock.AnythingOfType("*websocket.Conn"),
//		mock.AnythingOfType("map[uint32]*entity.ChatRoom"),
//	).Return(nil)
//
//	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTIxNTQwMzIsImp0aSI6IjEiLCJpYXQiOjE3MDk1NjIwMzIsImlzcyI6Ijpnb19jbGVhbl9hcmNoaXRlY3R1cmUifQ.lrR_1ufVsWlnVqqMhgF1w_RImMdiHVuVJ_MO7FcSxAc"
//
//	e := echo.New()
//	e.GET("/chats/:id", routes.handleConnections)
//
//	// テスト用サーバーの立ち上げ
//	server := httptest.NewServer(e)
//	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/chats/1?jwt=" + jwtToken
//	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
//	assert.NotEmpty(t, err)
//	assert.Equal(t, "websocket: bad handshake", err.Error())
//}
