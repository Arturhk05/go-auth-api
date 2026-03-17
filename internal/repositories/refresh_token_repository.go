package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, uuid.New(), userID, tokenHash, expiresAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" { // foreign key violation
				return fmt.Errorf("invalid user ID: %w", err)
			}
		}
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) ValidateRefreshToken(tokenHash string) (uuid.UUID, error) {
	query := `
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1 
			AND expires_at > NOW()
		 	AND revoked_at IS NULL
	`
	var userID uuid.UUID
	err := r.db.QueryRow(query, tokenHash).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("invalid refresh token: %w", err)
		}
		return uuid.Nil, err
	}

	return userID, nil
}

func (r *RefreshTokenRepository) RevokeByTokenHash(tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1
	`
	_, err := r.db.Exec(query, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) RevokeByUserId(userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteRefreshTokenByUserId(userID uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}
