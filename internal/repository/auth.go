package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) SaveRefreshToken(guid string, hashedRefreshToken []byte) error {
	query := `SELECT COUNT(*) FROM auth WHERE guid = $1`
	var exist bool
	err := r.db.QueryRow(query, guid).Scan(&exist)
	if err != nil {
		return err
	}
	if !exist {
		query := `INSERT INTO auth (guid, hashed_token) VALUES ($1, $2)`
		_, err := r.db.Exec(query, guid, hashedRefreshToken)
		if err != nil {
			return fmt.Errorf("failed to save refresh token: %w", err)
		}
		return nil
	} else {
		query := `UPDATE auth SET hashed_token = $1 WHERE guid = $2`
		_, err := r.db.Exec(query, hashedRefreshToken, guid)
		if err != nil {
			return fmt.Errorf("failed to save refresh token: %w", err)
		}
		return nil
	}

}

func (r *AuthPostgres) CheckRefreshToken(guid string) ([]byte, error) {
	query := `SELECT hashed_token FROM auth WHERE guid = $1`
	var storedHash []byte
	err := r.db.QueryRow(query, guid).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return storedHash, nil
}
