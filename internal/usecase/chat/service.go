package usecase

import "business/internal/entity"

// ChatUseCase -.
type ChatUseCase struct {
	repo ChatRepo
}

// New -.
func New(r ChatRepo) *ChatUseCase {
	return &ChatUseCase{
		repo: r,
	}
}

type StandardClaims struct {
	Issuer    string
	IssuedAt  uint32
	Id        string // ユーザーID
	ExpiresAt uint32 // 有効期限
	abc       string
}

// GetChats はチャット一覧を取得して返す。
func (uc *ChatUseCase) GetChats(userID uint32) (entity.Chats, error) {
	var chatList entity.Chats
	//user, err := uc.repo.GetUserByMail(param.Mail)
	//// DBと疎通できずエラーなのか、存在せずエラー(401)を分ける必要がある。
	//if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	//	return "", err
	//}
	//if err != nil {
	//	return "", fmt.Errorf("authentication - s.repo.GetUserByMail: %w", err)
	//}
	//
	//// isValidPassword　DBから取得した暗号化パスワードとリクエストのパスワードを暗号化して比較した結果、一致したらtrueが入る。
	//isValidPassword := user.IsValidPassword(param.Password)
	//if isValidPassword {
	//	token, err := generateToken(user.UserID)
	//	if err != nil {
	//		return "", err
	//	}
	//	return token, nil
	//} else {
	//	return "", errors.New("invalid password")
	//}

	return chatList, nil
}
