package repository

import (
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	SaveRefreshToken(guid string, hashedRefreshToken []byte) error
	CheckRefreshToken(hashedRefreshToken []byte) (string, string, error)

	GetRefreshTokenByToken(refreshToken string) (string, string, error) // hashedToken, tokenID, error
	UpdateRefreshToken(userID int64, newHashedRefreshToken string, tokenID string) error
	GetUserByRefreshToken(hashedRefreshToken string) (int64, error)
}

type Repository struct {
	AuthRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthRepository: NewAuthPostgres(db),
	}
}
