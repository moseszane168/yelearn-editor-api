package earthworm

import (
	"time"
)

// UserLearningActivities represents the user_learning_activities table
type UserLearningActivities struct {
	Duration     int         `json:"duration" db:"duration"`
	Metadata     interface{} `json:"metadata" db:"metadata"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	Date         time.Time   `json:"date" db:"date"`
	CourseId     string      `json:"course_id" db:"course_id"`
	UserId       string      `json:"user_id" db:"user_id"`
	ActivityType string      `json:"activity_type" db:"activity_type"`
	Id           string      `json:"id" db:"id"`
}
