package repo

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"business/db/model"
	"business/internal/entity"
	pkgmysql "business/pkg/mysql"
	"github.com/stretchr/testify/assert"
)

type MockTimeProvider struct {
	mock.Mock
}

func (m *MockTimeProvider) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func TestChatRepository_GetMessageList(t *testing.T) {
	// テスト用のデータベース接続を作成
	conn, err := pkgmysql.NewTest()
	assert.NoError(t, err)
	err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	err = conn.DB.AutoMigrate(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	CreateData(conn, t)

	defer func() {
		err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
		assert.NoError(t, err)
		err = conn.DB.AutoMigrate(model.User{}, model.ChatRoom{}, model.ChatMessage{})
		assert.NoError(t, err)
	}()

	// リポジトリのインスタンスを作成
	mockTimeProvider := &MockTimeProvider{}
	mockTimeProvider.
		On("Now").
		Return(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local))

	repo := New(conn, mockTimeProvider)

	// 期待値
	assertChatList := []entity.Message{
		{
			SenderUserID: 1,
			UserName:     "高橋太郎",
			Messages:     "テスト1",
			CreatedAt:    time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local),
		},
		{
			SenderUserID: 2,
			UserName:     "今井次郎",
			Messages:     "テスト2",
			CreatedAt:    time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local),
		},
	}
	// 正常系
	chatList, err := repo.GetMessageList(uint32(1))
	assert.NoError(t, err)
	assert.NotEmpty(t, chatList)
	assert.Equal(t, chatList, assertChatList)

	// 異常系:データが存在しないパターン
	chatList, err = repo.GetMessageList(uint32(10))
	assert.Empty(t, err)
	assert.Empty(t, chatList)
}

func TestChatRepository_CreateMessage(t *testing.T) {
	// テスト用のデータベース接続を作成
	conn, err := pkgmysql.NewTest()
	assert.NoError(t, err)
	err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	err = conn.DB.AutoMigrate(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	CreateData(conn, t)

	defer func() {
		err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
		assert.NoError(t, err)
		err = conn.DB.AutoMigrate(model.User{}, model.ChatRoom{}, model.ChatMessage{})
		assert.NoError(t, err)
	}()

	// リポジトリのインスタンスを作成
	mockTimeProvider := &MockTimeProvider{}
	mockTimeProvider.
		On("Now").
		Return(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local))

	repo := New(conn, mockTimeProvider)

	// 正常系
	chatMessage, err := repo.CreateMessage(1, "テスト4", 1)
	require.NoError(t, err)
	require.NotEmpty(t, chatMessage)
	require.Greater(t, chatMessage.ChatMessageID, uint32(0)) // IDが有効かチェック
	require.Equal(t, uint32(1), chatMessage.ChatRoomID)
	require.Equal(t, "テスト4", chatMessage.Message)
	require.Equal(t, uint32(1), chatMessage.SenderUserID)
	require.True(t, chatMessage.CreatedAt.Equal(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local)))

	// 異常系
	err = conn.DB.Migrator().DropTable(model.User{}, model.ChatRoom{}, model.ChatMessage{})
	assert.NoError(t, err)
	chatMessage, err = repo.CreateMessage(1, "テスト4", 1)
	assert.NotEmpty(t, err)
	assert.Empty(t, chatMessage)
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
	err = conn.DB.Create(&userList).Error
	assert.NoError(t, err)

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
	err = conn.DB.Create(&chatList).Error
	assert.NoError(t, err)

	ChatMessageList := []model.ChatMessage{
		{
			ChatRoomID:   1,
			Message:      "テスト1",
			SenderUserID: 1,
			CreatedAt:    time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local),
		},
		{
			ChatRoomID:   1,
			Message:      "テスト2",
			SenderUserID: 2,
			CreatedAt:    time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local),
		},
		{
			ChatRoomID:   2,
			Message:      "テスト3",
			SenderUserID: 3,
			CreatedAt:    time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local),
		},
	}
	err = conn.DB.Create(&ChatMessageList).Error
	assert.NoError(t, err)
}
