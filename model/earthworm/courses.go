package earthworm

import (
	"time"
)

// Courses represents the courses table
type Courses struct {
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Order        int       `json:"order" db:"order"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	Video        string    `json:"video" db:"video"`
	Id           string    `json:"id" db:"id"`
	CoursePackId string    `json:"course_pack_id" db:"course_pack_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
}
