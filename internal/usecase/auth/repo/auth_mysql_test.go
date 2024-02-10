package repo

import (
	"business/db/model"
	"testing"

	pkgmysql "business/pkg/mysql"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepository_GetUserByMail(t *testing.T) {
	// テスト用のデータベース接続を作成
	conn, err := pkgmysql.NewTest()
	assert.NoError(t, err)
	err = conn.DB.Migrator().DropTable(model.User{})
	assert.NoError(t, err)
	err = conn.DB.AutoMigrate(model.User{})
	assert.NoError(t, err)

	userList := []model.User{
		{
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
			LastName:          "今井",
			FirstName:         "次郎",
			HiraganaLastName:  "いまい",
			HiraganaFirstName: "たろう",
			Email:             "abcd@co.jp",
			Password:          "パスワード",
			CreatedUserID:     0,
			UpdateUserID:      0,
		},
	}

	for _, user := range userList {
		err = conn.DB.Create(&user).Error
		assert.NoError(t, err)
	}

	// リポジトリのインスタンスを作成
	repo := New(conn) // New関数には*gorm.DBインスタンスを渡す

	// 正常系
	user, err := repo.GetUserByMail("abc@co.jp")
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	// 異常系
	user, err = repo.GetUserByMail("存在しないメールアドレス")
	assert.Error(t, err)
	assert.Empty(t, user)

	err = conn.DB.Migrator().DropTable(model.User{})
	assert.NoError(t, err)
	err = conn.DB.AutoMigrate(model.User{})
	assert.NoError(t, err)
}
