package entity

import (
	"github.com/fikrialwan/FitByte/pkg/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name       string    `gorm:"type:varchar(65)" json:"name"`
	Email      string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Password   string    `gorm:"type:text" json:"password"`
	Preference string    `gorm:"type:varchar(10)" json:"preference"`
	WeightUnit string    `gorm:"type:varchar(5)" json:"weight_unit"`
	HeightUnit string    `gorm:"type:varchar(5)" json:"height_unit"`
	Weight     int       `json:"weight"`
	Height     int       `json:"height"`
	ImageUri   string    `gorm:"type:text" json:"image_uri"`

	Timestamp
}

// BeforeCreate hook to hash password and set defaults
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	// Hash password
	if u.Password != "" {
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}

	// Ensure UUID is set
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	return nil
}

// BeforeUpdate hook to handle password updates
func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	// Only hash password if it has been changed
	if u.Password != "" {
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}
	return nil
}
