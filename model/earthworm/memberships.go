package earthworm

import (
	"time"
)

// Memberships represents the memberships table
type Memberships struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	IsActive  bool      `json:"isActive" db:"isActive"`
	Type      string    `json:"type" db:"type"`
	UserId    string    `json:"user_id" db:"user_id"`
	Id        string    `json:"id" db:"id"`
}
