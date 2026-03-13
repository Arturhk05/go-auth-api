package repositories

import (
	"database/sql"
	"fmt"

	"github.com/arturhk05/go-auth-api/internal/models"
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
