package entity

import "time"

type SensorLog struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	DistanceCM  *float64  `db:"distance_cm"`
	LDRValue    *int      `db:"ldr_value"`
	PIRDetected bool      `db:"pir_detected"`
	Condition   string    `db:"condition"`
	LogType     string    `db:"log_type"` // "heartbeat" | "condition_change"
	RecordedAt  time.Time `db:"recorded_at"`
}
