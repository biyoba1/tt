package repository

import (
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	SaveRefreshToken(guid string, hashedRefreshToken []byte) error
	CheckRefreshToken(guid string) ([]byte, error)
}

type Repository struct {
	AuthRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthRepository: NewAuthPostgres(db),
	}
}
