package usecase

import (
	"errors"
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

func (m *MockChatRepo) GetChatList(userID uint32) (entity.ChatRooms, error) {
	args := m.Called(userID)
	return args.Get(0).(entity.ChatRooms), args.Error(1)
}

// Authenticationメソッドのテスト
func TestAuthentication(t *testing.T) {
	mockRepo := new(MockChatRepo)
	uc := New(mockRepo)

	chatList := entity.ChatRooms{
		List: []entity.Room{
			{
				ChatRoomID: 1,
				UserName:   "斎藤太郎",
				Message:    func() *string { s := "テスト1"; return &s }(),
			},
		},
	}

	// 正常系
	mockRepo.On("GetChatList", uint32(1)).Return(chatList, nil)
	chats, err := uc.GetChats(uint32(1))
	assert.NoError(t, err)
	assert.NotEmpty(t, chats)

	// 異常系
	mockRepo.On("GetChatList", uint32(0)).Return(entity.ChatRooms{}, gorm.ErrRecordNotFound) // int32(0)からuint32(0)へ修正
	chats, err = uc.GetChats(uint32(0))
	assert.Error(t, err)
	assert.Empty(t, chats)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	// 500異常系
	mockRepo.On("GetChatList", uint32(2)).Return(entity.ChatRooms{}, errors.New("予期せぬエラー"))
	chats, err = uc.GetChats(uint32(2))
	assert.Error(t, err)
	assert.Empty(t, chats)
	assert.EqualError(t, err, "GetChats - s.repo.GetChatList: 予期せぬエラー")
}
