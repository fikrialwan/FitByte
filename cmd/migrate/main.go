package main

import (
	"log"

	"github.com/fikrialwan/FitByte/config"
	"github.com/fikrialwan/FitByte/database"
)

func main() {
	db := config.InitDb()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Error migration: %v", err)
	}
	log.Println("complete migration successfully!")
}
