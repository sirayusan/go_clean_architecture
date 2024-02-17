package usecase

import (
	"gorm.io/gorm"
	"testing"

	"business/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// AuthRepoのモックを作成
type MockChatRepo struct {
	mock.Mock
}

func (m *MockChatRepo) GetChatList(userID uint32) (entity.Chats, error) {
	args := m.Called(userID)
	return args.Get(0).(entity.Chats), args.Error(1)
}

// Authenticationメソッドのテスト
func TestAuthentication(t *testing.T) {
	mockRepo := new(MockChatRepo)
	uc := New(mockRepo)

	chatList := entity.Chats{
		[]entity.Chat{
			{
				RecipientUserName: "斎藤太郎",
				Message:           func() *string { s := "テスト1"; return &s }(),
			},
		},
	}

	// 正常系
	mockRepo.On("GetChatList", "test@example.com").Return(chatList, nil)
	chats, err := uc.GetChats(uint32(1))
	assert.NoError(t, err)
	assert.NotEmpty(t, chats)

	// 異常系
	mockRepo.On("GetChatList", "test@example.com").Return(entity.Chats{}, gorm.ErrRecordNotFound)
	chats, err = uc.GetChats(uint32(0))
	assert.Error(t, err)
	assert.Empty(t, chats)
	assert.Equal(t, err, gorm.ErrRecordNotFound)
}
