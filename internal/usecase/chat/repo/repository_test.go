package repo

import (
	"testing"
	"time"

	"business/db/model"
	"business/internal/entity"
	pkgmysql "business/pkg/mysql"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepository_GetUserByMail(t *testing.T) {

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
	repo := New(conn)

	// 期待値
	assertChatList := entity.ChatRooms{
		List: []entity.Room{
			{
				ChatRoomID: 1,
				UserName:   "今井次郎",
				Message:    func() *string { s := "テスト2"; return &s }(),
				CreatedAt:  func() *time.Time { time_ := time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local); return &time_ }(),
			},
			{
				ChatRoomID: 2,
				UserName:   "斎藤三郎",
				Message:    func() *string { s := "テスト3"; return &s }(),
				CreatedAt:  func() *time.Time { time_ := time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local); return &time_ }(),
			},
		},
	}
	// 正常系
	chatList, err := repo.GetChatList(uint32(1))
	assert.NoError(t, err)
	assert.NotEmpty(t, chatList)
	assert.Equal(t, chatList, assertChatList)

	// 異常系:データが存在しないパターン
	chatList, err = repo.GetChatList(uint32(10))
	assert.Error(t, err)
	assert.Empty(t, chatList)
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
			CreatedAt:    time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local),
		},
		{
			ChatRoomID:   1,
			Message:      "テスト2",
			SenderUserID: 1,
			CreatedAt:    time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local),
		},
		{
			ChatRoomID:   2,
			Message:      "テスト3",
			SenderUserID: 3,
			CreatedAt:    time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local),
		},
	}
	for _, chatMessage := range ChatMessageList {
		err = conn.DB.Create(&chatMessage).Error
		assert.NoError(t, err)
	}
}
