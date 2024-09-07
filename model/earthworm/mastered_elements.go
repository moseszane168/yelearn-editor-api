package earthworm

import (
	"time"
)

// MasteredElements represents the mastered_elements table
type MasteredElements struct {
	Content    interface{} `json:"content" db:"content"`
	MasteredAt time.Time   `json:"mastered_at" db:"mastered_at"`
	Id         string      `json:"id" db:"id"`
	UserId     string      `json:"user_id" db:"user_id"`
}
