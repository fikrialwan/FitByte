package entity

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
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

type ActivityConfig struct {
	CaloriesPerMinute int
}

var activityConfigs = map[ActivityType]ActivityConfig{
	Walking:    {CaloriesPerMinute: 4},
	Yoga:       {CaloriesPerMinute: 4},
	Stretching: {CaloriesPerMinute: 4},

	Cycling:  {CaloriesPerMinute: 8},
	Swimming: {CaloriesPerMinute: 8},
	Dancing:  {CaloriesPerMinute: 8},

	Hiking:   {CaloriesPerMinute: 10},
	Running:  {CaloriesPerMinute: 10},
	HIIT:     {CaloriesPerMinute: 10},
	JumpRope: {CaloriesPerMinute: 10},
}

var (
	validActivityTypesList = []ActivityType{
		Walking, Yoga, Stretching, Cycling, Swimming, Dancing, Hiking, Running, HIIT, JumpRope,
	}

	validActivityTypeStrings = []string{
		"Walking", "Yoga", "Stretching", "Cycling", "Swimming", "Dancing", "Hiking", "Running", "HIIT", "JumpRope",
	}
)

func (a ActivityType) IsValid() bool {
	_, exists := activityConfigs[a]
	return exists
}

func (a ActivityType) CaloriesPerMinute() int {
	if config, exists := activityConfigs[a]; exists {
		return config.CaloriesPerMinute
	}
	return 0
}

func (a ActivityType) CalculateBurnedCalories(durationInMinutes int) int {
	return a.CaloriesPerMinute() * durationInMinutes
}

func GetValidActivityTypes() []ActivityType {
	return validActivityTypesList
}

func GetValidActivityTypeStrings() []string {
	return validActivityTypeStrings
}
