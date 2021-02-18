package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

func checkMigrations() {
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:5432/%s?sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_DB"),
		))
	if err != nil {
		log.Fatal(err)
	}

	// XXX: https://github.com/golang-migrate/migrate/issues/179
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	log.Print("migrations checked")
}
