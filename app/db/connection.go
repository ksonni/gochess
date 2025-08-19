package db

import (
	"database/sql"
	"fmt"
	"gochess/lib/env"

	_ "github.com/lib/pq"
)

const kMigrationsPath = "file://db/migrations"

func Connection() (*sql.DB, error) {
	host := env.MustEnv("DB_HOST")
	port := env.MustEnv("DB_PORT")
	user := env.MustEnv("DB_USER")
	password := env.MustEnv("DB_PASSWORD")
	dbname := env.MustEnv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	return sql.Open("postgres", connStr)
}
