package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattrmcg/equitalytics-backend/config"
	"github.com/mattrmcg/equitalytics-backend/internal/db"
	"github.com/mattrmcg/equitalytics-backend/internal/server"
)

func main() {
	// create connection pool
	dbPool, err := db.CreateDBPool(config.Envs.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	// close connection pool when program terminates
	defer db.CloseDBPool(dbPool)

	// Check if connection can be established
	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}
	log.Println("DB: Successfully Connected!")

	// Create and run server
	server := server.NewAPIServer(fmt.Sprintf(":%v", config.Envs.Port), dbPool)
	if err = server.Run(); err != nil {
		log.Fatalf("unable to start server: %v", err)
	}

}
