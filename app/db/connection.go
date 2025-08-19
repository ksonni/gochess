package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const kMigrationsPath = "file://db/migrations"

func Connection() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	return sql.Open("postgres", connStr)
}

func Migrate() error {
	conn, err := Connection()
	if err != nil {
		return err
	}
	defer conn.Close()

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(kMigrationsPath, "postgres", driver)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}
	return nil
}
