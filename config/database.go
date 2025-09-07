package config

import (
	"log"

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
		sqlDB.SetMaxIdleConns(cfg.GetDBMaxIdleConns())
		sqlDB.SetMaxOpenConns(cfg.GetDBMaxOpenConns())
		sqlDB.SetConnMaxLifetime(cfg.GetDBConnMaxLifetime())
		sqlDB.SetConnMaxIdleTime(cfg.GetDBConnMaxIdleTime())
	}

	RunExtension(db)

	return db
}
