// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"business/internal/entity"
	"context"
)

type (
	// User -.
	User interface {
		UserList(context.Context) ([]entity.User, error)
	}

	// UserRepo -.
	UserRepo interface {
		GetUserList(context.Context) ([]entity.User, error)
	}
)
