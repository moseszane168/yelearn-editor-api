package earthworm

import (
	"time"
)

// UserCourseProgress represents the user_course_progress table
type UserCourseProgress struct {
	StatementIndex int       `json:"statement_index" db:"statement_index"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	Id             string    `json:"id" db:"id"`
	CourseId       string    `json:"course_id" db:"course_id"`
	UserId         string    `json:"user_id" db:"user_id"`
	CoursePackId   string    `json:"course_pack_id" db:"course_pack_id"`
}
