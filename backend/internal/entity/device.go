package entity

import "time"

type DeviceConfig struct {
	ID                 string    `db:"id"`
	UserID             string    `db:"user_id"`
	DistanceMinCM      int       `db:"distance_min_cm"`
	DistanceMaxCM      int       `db:"distance_max_cm"`
	LDRThreshold       int       `db:"ldr_threshold"`
	AwayTimeoutMinutes int       `db:"away_timeout_minutes"`
	ConfigAck          bool      `db:"config_ack"`
	UpdatedAt          time.Time `db:"updated_at"`
}

type DeviceStatus struct {
	ID       string     `db:"id"`
	UserID   string     `db:"user_id"`
	IsOnline bool       `db:"is_online"`
	LastSeen *time.Time `db:"last_seen"`
}
