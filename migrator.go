package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Application entry point

	migraton, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			"admin",
			"1234",
			"localhost",
			5432,
			"sobaki",
		))

	if err != nil {
		fmt.Println(err)
	}

	if err := migraton.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}
