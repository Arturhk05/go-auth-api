package database

import (
	"database/sql"
	"fmt"

	"github.com/arturhk05/go-auth-api/config"
)

type postgresDB struct {
	Db *sql.DB
}

func NewPostgresDB() (*postgresDB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.Env.Database.User, config.Env.Database.Password, config.Env.Database.Host, config.Env.Database.Port, config.Env.Database.Name, config.Env.Database.SSLMode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &postgresDB{Db: db}, nil
}
