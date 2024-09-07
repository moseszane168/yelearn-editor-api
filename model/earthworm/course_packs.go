package earthworm

import (
	"time"
)

// CoursePacks represents the course_packs table
type CoursePacks struct {
	IsFree      bool      `json:"is_free" db:"is_free"`
	Order       int       `json:"order" db:"order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatorId   string    `json:"creator_id" db:"creator_id"`
	ShareLevel  string    `json:"share_level" db:"share_level"`
	Cover       string    `json:"cover" db:"cover"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Id          string    `json:"id" db:"id"`
}
