package entity

import (
	"github.com/google/uuid"
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
}
