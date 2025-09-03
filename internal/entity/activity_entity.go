package entity

import (
	"time"

	"github.com/google/uuid"
)

type Acticity struct {
	ID                uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"activityId"`
	ActivityType      ActivityType `gorm:"type:varchar(15)" json:"activityType"`
	DoneAt            time.Time    `json:"doneAt"`
	DurationInMinutes int          `json:"durationInMinutes"`
	CaloriesBurned    int          `json:"caloriesBurned"`

	UserID uuid.UUID
	User   User

	Timestamp
}

type ActivityType string

const (
	Walking    ActivityType = "Walking"
	Yoga       ActivityType = "Yoga"
	Stretching ActivityType = "Stretching"
	Cycling    ActivityType = "Cycling"
	Swimming   ActivityType = "Swimming"
	Dancing    ActivityType = "Dancing"
	Hiking     ActivityType = "Hiking"
	Running    ActivityType = "Running"
	HIIT       ActivityType = "HIIT"
	JumpRope   ActivityType = "JumpRope"
)

// CaloriesPerMinute returns the calories burned per minute for each activity type
func (a ActivityType) CaloriesPerMinute() int {
	switch a {
	case Walking, Yoga, Stretching:
		return 4
	case Cycling, Swimming, Dancing:
		return 8
	case Hiking, Running, HIIT, JumpRope:
		return 10
	default:
		return 0
	}
}
