package usecase

import (
	"business/internal/entity"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockChatMessageRepository struct {
	mock.Mock
}

func (m *MockChatMessageRepository) GetMessageList(chatID uint32) ([]entity.Message, error) {
	args := m.Called(chatID)
	return args.Get(0).([]entity.Message), args.Error(1)
}

func (m *MockChatMessageRepository) CreateMessage(chatRoomID uint32, msg string, SenderUserID uint32) (entity.ChatMessage, error) {
	args := m.Called(chatRoomID)
	return args.Get(0).(entity.ChatMessage), args.Error(1)
}

type MockWebSocket struct {
	mock.Mock
}

// WriteMessage は MockConn インターフェースのメソッドです。
func (m *MockWebSocket) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}

// MockWebSocketConnWrapper は WebSocket コネクションをラップする構造体です。
type MockWebSocketConnWrapper struct {
	Conn entity.WebSocketWrapper
}

// MockWriteMessage は WebSocket コネクションにメッセージを書き込むためのメソッドです。
func (w *MockWebSocketConnWrapper) WriteMessage(messageType int, data []byte) error {
	return w.Conn.WriteMessage(messageType, data)
}

// テストケース
func TestMessageUseCase_JoinRoom(t *testing.T) {
	mockRepo := new(MockChatMessageRepository)
	uc := New(mockRepo)
	_time := time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local)

	// メッセージリストのモック設定
	expectedMessages := []entity.Message{{
		SenderUserID: 1,
		UserName:     "aaa",
		Messages:     "テスト1",
		CreatedAt:    _time,
	}}
	mockRepo.On("GetMessageList", uint32(1)).Return(expectedMessages, nil)

	ws := &MockWebSocket{}
	wws := &MockWebSocketConnWrapper{Conn: ws}

	// WebSocketWrapperのWriteMessageをモック
	ws.On("WriteMessage", mock.Anything, mock.Anything).Return(nil)

	roomManager := make(map[uint32]*entity.ChatRoom)
	err := uc.JoinRoom(uint32(1), wws, roomManager)
	assert.NoError(t, err)

	// クライアントがルームに追加されたことを確認
	room, exists := roomManager[uint32(1)]
	assert.True(t, exists)
	assert.Equal(t, 1, len(room.Clients))

	// メッセージが送信されたことを確認
	ws.AssertCalled(t, "WriteMessage", websocket.TextMessage, mock.Anything)
}
