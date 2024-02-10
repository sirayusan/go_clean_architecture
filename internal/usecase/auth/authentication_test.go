package usecase

import (
	"errors"
	"testing"

	"business/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// AuthRepoのモックを作成
type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) GetUserByMail(mail string) (entity.LoginUser, error) {
	args := m.Called(mail)
	return args.Get(0).(entity.LoginUser), args.Error(1)
}

// Authenticationメソッドのテスト
func TestAuthentication(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	uc := New(mockRepo)

	// モックが返すべきユーザー情報を設定
	// EncryptedPasswordはパスワードを暗号化した場合。
	user := entity.LoginUser{UserID: 1, EncryptedPassword: "$2a$10$AjABQamoZJoLASTRKUzQr.o5eMaYRMnUDkjDl2vbUHoEIXZ0ZyaJi"}
	mockRepo.On("GetUserByMail", "test@example.com").Return(user, nil)

	// Authenticationメソッドを呼び出し
	token, err := uc.GenerateJwtToken(entity.LoginRequest{Mail: "test@example.com", Password: "パスワード"})

	// パスワード検証を行う
	isValid := user.IsValidPassword("パスワード")

	// 正常系
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, isValid)

	// 存在しないユーザーのケース
	mockRepo.On("GetUserByMail", "unknown@example.com").Return(entity.LoginUser{}, gorm.ErrRecordNotFound)
	token, err = uc.GenerateJwtToken(entity.LoginRequest{Mail: "unknown@example.com", Password: "password"})
	assert.Error(t, err)
	assert.Empty(t, token)

	//　パスワードが一致しないケース
	user = entity.LoginUser{UserID: 1, EncryptedPassword: "暗号化後の値でないパスワード"}
	isValid = user.IsValidPassword("パスワード")
	assert.False(t, isValid)

	// DBエラーのケース
	mockRepo.On("GetUserByMail", "dberror@example.com").Return(entity.LoginUser{}, errors.New("db error"))
	token, err = uc.GenerateJwtToken(entity.LoginRequest{Mail: "dberror@example.com", Password: "password"})
	assert.Error(t, err)
	assert.Empty(t, token)
}
