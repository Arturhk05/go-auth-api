package database

import (
	"database/sql"
	"fmt"

	"github.com/arturhk05/go-auth-api/config"
)

type postgresDB struct {
	Db *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*postgresDB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.SSLMode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &postgresDB{Db: db}, nil
}
