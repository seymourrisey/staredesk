package main

import (
	"log"

	"github.com/seymourrisey/staredesk/config"
	"github.com/seymourrisey/staredesk/internal/infrastructure/postgres"
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewPool(&cfg.DB)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	log.Printf("Connected to database")
	log.Printf("StareDesk backend starting on port %s", cfg.App.Port)
}
