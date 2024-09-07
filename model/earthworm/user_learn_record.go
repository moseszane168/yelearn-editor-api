package earthworm

import (
	"time"
)

// UserLearnRecord represents the user_learn_record table
type UserLearnRecord struct {
	Count     int       `json:"count" db:"count"`
	Day       time.Time `json:"day" db:"day"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Id        string    `json:"id" db:"id"`
	UserId    string    `json:"user_id" db:"user_id"`
}
