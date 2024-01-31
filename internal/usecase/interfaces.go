// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"business/internal/entity"
	"context"
)

type (
	// User -.
	User interface {
		Translate(context.Context, entity.User) (entity.User, error)
		UserList(context.Context) ([]entity.User, error)
	}

	// UserRepo -.
	UserRepo interface {
		Store(context.Context, entity.User) error
		GetUserList(context.Context) ([]entity.User, error)
	}

	// UserWebAPI -.
	UserWebAPI interface {
		Translate(entity.User) (entity.User, error)
	}
)
