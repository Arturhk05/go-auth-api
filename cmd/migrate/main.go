package main

import (
	"log"
	"os"

	"github.com/arturhk05/go-auth-api/config"
	"github.com/arturhk05/go-auth-api/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Db.Close()

	driver, err := postgres.WithInstance(postgresDB.Db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	switch cmd := os.Args[len(os.Args)-1]; cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatalf("Unknown command: %s. Use 'up' or 'down'.", cmd)
	}
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
}
