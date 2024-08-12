package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattrmcg/equitalytics-backend/config"
)

func main() {
	// db, err := db.CreateDBPool(config.Envs.DBURL)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// driver, err := pgx.WithInstance()

	dbUrl := fmt.Sprintf("%v?sslmode=require", config.Envs.DBURL)
	m, err := migrate.New(
		"file://cmd/migrate/migrations",
		dbUrl,
	)
	if err != nil {
		log.Fatalf("Couldn't create migrate: %v\n", err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

}
