package main

import (
	"log"

	"github.com/fikrialwan/FitByte/config"
	"github.com/fikrialwan/FitByte/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config")
	}
	db := config.InitDb(cfg)

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Error migration: %v", err)
	}
	log.Println("complete migration successfully!")
}
