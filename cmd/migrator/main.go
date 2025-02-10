package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"url-shortener/internal/config"
)

const BasePrefPath string = "file://"

func main() {
	var migrationsPath, migrationsTable string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()
	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	cfg := config.MustLoad()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&x-migrations-table=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
		migrationsTable,
	)

	m, err := migrate.New(
		BasePrefPath+migrationsPath,
		dsn,
	)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return
		}

		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("migrations applied")
}
