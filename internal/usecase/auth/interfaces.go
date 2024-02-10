package usecase

import (
	"business/internal/entity"
)

type (
	// Auth -.
	Auth interface {
		GenerateJwtToken(entity.LoginRequest) (string, error)
	}

	// AuthRepo -.
	AuthRepo interface {
		GetUserByMail(string) (entity.LoginUser, error)
	}
)
