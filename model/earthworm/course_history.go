package earthworm

import (
	"time"
)

// CourseHistory represents the course_history table
type CourseHistory struct {
	CompletionCount int       `json:"completion_count" db:"completion_count"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	Id              string    `json:"id" db:"id"`
	CoursePackId    string    `json:"course_pack_id" db:"course_pack_id"`
	UserId          string    `json:"user_id" db:"user_id"`
	CourseId        string    `json:"course_id" db:"course_id"`
}
