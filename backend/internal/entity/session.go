package entity

import "time"

type Session struct {
	ID                string     `db:"id"`
	UserID            string     `db:"user_id"`
	StartedAt         time.Time  `db:"started_at"`
	EndedAt           *time.Time `db:"ended_at"`
	DurationSec       *int       `db:"duration_sec"`
	DominantCondition *string    `db:"dominant_condition"`
}
