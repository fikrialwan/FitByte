package config

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RunExtension(db *gorm.DB) {
	// extension to generate uuid
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}

func InitDb(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.GetDSN(),
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(50)
		sqlDB.SetMaxOpenConns(200)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	}

	RunExtension(db)

	return db
}
