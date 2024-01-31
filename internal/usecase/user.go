package usecase

import (
	"context"
	"fmt"

	"business/internal/entity"
)

// UserUseCase -.
type UserUseCase struct {
	repo UserRepo
}

// New -.
func New(r UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
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
