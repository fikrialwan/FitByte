package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RunExtension(db *gorm.DB) {
	// extension to generate uuid
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}

func validateRequiredEnvVars() {
	required := []string{"DB_USER", "DB_PASS", "DB_HOST", "DB_NAME", "DB_PORT", "JWT_SECRET"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			panic(fmt.Sprintf("Required environment variable %s is not set", env))
		}
	}
}

func InitDb() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	validateRequiredEnvVars()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
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
