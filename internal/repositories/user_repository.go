package repositories

import (
	"database/sql"
	"fmt"

	"github.com/arturhk05/go-auth-api/internal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	query := `
		INSERT INTO users (id, email, username, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`
	err := r.db.QueryRow(query, user.ID, user.Email, user.Username, user.PasswordHash).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, is_active, email_verified, last_login_at, failed_login_attempts, locked_until
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(id uuid.UUID) (*models.User, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("id cannot be nil")
	}

	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, is_active, email_verified, last_login_at, failed_login_attempts, locked_until
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateLoginAttempts(id uuid.UUID, failedAttempts int) error {
	query := `
		UPDATE users
		SET failed_login_attempts = $2
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, failedAttempts)
	if err != nil {
		return fmt.Errorf("failed to update login attempts: %w", err)
	}

	return nil
}

func (r *UserRepository) LockAccountAndResetLoginAttempts(id uuid.UUID, lockTime sql.NullTime) error {
	query := `
		UPDATE users
		SET locked_until = $2, failed_login_attempts = 0
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, lockTime)
	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}

	return nil
}

func (r *UserRepository) ResetLoginAttempts(id uuid.UUID) error {
	query := `
		UPDATE users
		SET failed_login_attempts = 0
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to reset login attempts: %w", err)
	}

	return nil
}
