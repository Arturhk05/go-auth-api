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
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", config.Env.DbUser, config.Env.DbPassword, config.Env.DbHost, config.Env.DbPort, config.Env.DbName, config.Env.DbSSLMode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &postgresDB{Db: db}, nil
}
