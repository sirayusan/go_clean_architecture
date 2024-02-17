package usecase

import (
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
func (uc *UserUseCase) UserList() ([]entity.User, error) {
	translations, err := uc.repo.GetUserList()
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return translations, nil
}
