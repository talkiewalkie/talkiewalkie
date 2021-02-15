package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

func checkMigrations() {
	m, err := migrate.New(
		"file://migrations",
		"postgres://theo:@localhost:5432/talkiewalkie?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// XXX: https://github.com/golang-migrate/migrate/issues/179
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	log.Print("migrations checked")
}
