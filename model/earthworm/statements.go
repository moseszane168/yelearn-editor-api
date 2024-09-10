package earthworm

import (
	"time"
)

// Statements represents the statements table
type Statements struct {
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Id        string    `json:"id" db:"id"`
	CourseId  string    `json:"course_id" db:"course_id"`
	Soundmark string    `json:"soundmark" db:"soundmark"`
	Chinese   string    `json:"chinese" db:"chinese"`
	English   string    `json:"english" db:"english"`
	Pid       string    `json:"pid" db:"pid"`
}
