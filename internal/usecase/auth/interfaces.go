package usecase

import (
	"business/internal/entity"
)

type (
	// Auth -.
	Auth interface {
		Authentication(entity.LoginRequest) (string, error)
	}

	// AuthRepo -.
	AuthRepo interface {
		GetUserByMail(string) (entity.LoginUser, error)
	}
)
