package usecase

import (
	"context"
	"fmt"

	"business/internal/entity"
)

// UserUseCase -.
type UserUseCase struct {
	repo   UserRepo
	webAPI UserWebAPI
}

// New -.
func New(r UserRepo, w UserWebAPI) *UserUseCase {
	return &UserUseCase{
		repo:   r,
		webAPI: w,
	}
}

// UserList - getting translate history from store.
func (uc *UserUseCase) UserList(ctx context.Context) ([]entity.User, error) {
	translations, err := uc.repo.GetUserList(ctx)
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return translations, nil
}

// Translate -.
func (uc *UserUseCase) Translate(ctx context.Context, t entity.User) (entity.User, error) {
	//translation, err := uc.webAPI.Translate(t)
	//if err != nil {
	//	return entity.User{}, fmt.Errorf("UserUseCase - Translate - s.webAPI.Translate: %w", err)
	//}
	//
	//err = uc.repo.Store(context.Background(), translation)
	//if err != nil {
	//	return entity.User{}, fmt.Errorf("UserUseCase - Translate - s.repo.Store: %w", err)
	//}

	return entity.User{}, nil
}
