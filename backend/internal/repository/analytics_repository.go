package repository

import (
	"context"
	"time"
)

// PeakHourEntry adalah data fokus per jam dalam sehari.
type PeakHourEntry struct {
	Hour         int // 0–23
	TotalSec     int
	SessionCount int
}

// ConditionBreakdownEntry adalah proporsi per kondisi.
type ConditionBreakdownEntry struct {
	Condition string
	Count     int
}

// TimelineEntry adalah data kondisi dominan per jam untuk satu hari.
type TimelineEntry struct {
	Hour              int
	DominantCondition string
	TotalSec          int
}

type AnalyticsRepository interface {
	// GetPeakHours mengembalikan total durasi fokus per jam dalam range waktu.
	GetPeakHours(ctx context.Context, userID string, from time.Time, to time.Time) ([]*PeakHourEntry, error)

	// GetConditionBreakdown mengembalikan jumlah log per kondisi dalam range waktu.
	GetConditionBreakdown(ctx context.Context, userID string, from time.Time, to time.Time) ([]*ConditionBreakdownEntry, error)

	// GetTimeline mengembalikan kondisi dominan per jam untuk satu hari.
	GetTimeline(ctx context.Context, userID string, date time.Time) ([]*TimelineEntry, error)
}
