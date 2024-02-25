package usecase

import (
	"business/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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

// MockWebSocketWrapper はテスト用のWebSocket操作を模倣します。
type MockWebSocketWrapper struct {
	// 必要に応じてテスト用のフィールドやメソッドを追加
}

// WriteMessage はメッセージ送信処理のモック
func (m *MockWebSocketWrapper) WriteMessage(messageType int, data []byte) error {
	return nil
}

func TestMessageUseCase_JoinRoom(t *testing.T) {
	mockRepo := new(MockChatMessageRepository)
	uc := New(mockRepo)

	ws := &MockWebSocketWrapper{} // モックのインスタンスを作成
	roomManager := make(map[uint32]*entity.ChatRoom)
	mockRepo.On("GetMessageList", uint32(1)).Return([]entity.Message{}, nil)

	err := uc.JoinRoom(uint32(1), ws, roomManager) // clientを渡す
	assert.NoError(t, err)
}
