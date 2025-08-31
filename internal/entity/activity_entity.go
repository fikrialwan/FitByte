package entity

import (
	"time"

	"github.com/google/uuid"
)

type Acticity struct {
	ID                uuid.UUID
	ActivityType      string    `gorm:"type:varchar(15)" json:"activity_type"`
	DoneAt            time.Time `json:"done_at"`
	DurationInMinutes int       `json:"duration_in_minutes"`
	CaloriesBurned    int       `json:"calories_burned"`

	UserID uuid.UUID
	User   User

	Timestamp
}
