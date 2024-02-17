package usecase

import (
	"business/internal/entity"
)

type (
	// User -.
	User interface {
		UserList() ([]entity.User, error)
	}

	// UserRepo -.
	UserRepo interface {
		GetUserList() ([]entity.User, error)
	}
)
