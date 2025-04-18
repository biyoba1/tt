package repository

import (
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) SaveRefreshToken(guid string, hashedRefreshToken []byte) error {
	query := `INSERT INTO auth (guid, hashed_token) VALUES ($1, $2)`
	_, err := r.db.Exec(query, guid, hashedRefreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) CheckRefreshToken(guid []byte) (string, string, error) {
	query := `SELECT hashed_token FROM %s WHERE guid=$1`
	res, err := r.db.Exec(query, guid)
}

func (r *AuthPostgres) GetRefreshTokenByToken(refreshToken string) (string, string, error) {
	return "", "", nil
}

func (r *AuthPostgres) UpdateRefreshToken(userID int64, newHashedRefreshToken string, tokenID string) error {
	return nil
}

func (r *AuthPostgres) GetUserByRefreshToken(hashedRefreshToken string) (int64, error) {
	return 0, nil
}
