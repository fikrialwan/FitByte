package database

import (
	"github.com/fikrialwan/FitByte/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Acticity{},
	); err != nil {
		return err
	}

	return nil
}
