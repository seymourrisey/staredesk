package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seymourrisey/staredesk/internal/repository"
)

type AnalyticsPostgres struct {
	db *pgxpool.Pool
}

func NewAnalyticsPostgres(db *pgxpool.Pool) *AnalyticsPostgres {
	return &AnalyticsPostgres{db: db}
}

// GetPeakHours — total duration_sec per jam dari tabel sessions dalam range waktu.
// Menggunakan sessions (bukan sensor_logs) karena duration_sec sudah dihitung.
func (r *AnalyticsPostgres) GetPeakHours(ctx context.Context, userID string, from time.Time, to time.Time) ([]*repository.PeakHourEntry, error) {
	query := `
		SELECT
			EXTRACT(HOUR FROM started_at)::int AS hour,
			COALESCE(SUM(duration_sec), 0)::int AS total_sec,
			COUNT(*)::int AS session_count
		FROM sessions
		WHERE user_id = $1
		  AND ended_at IS NOT NULL
		  AND started_at >= $2
		  AND started_at < $3
		GROUP BY hour
		ORDER BY hour
	`
	rows, err := r.db.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*repository.PeakHourEntry
	for rows.Next() {
		e := &repository.PeakHourEntry{}
		if err := rows.Scan(&e.Hour, &e.TotalSec, &e.SessionCount); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

// GetConditionBreakdown — jumlah log per kondisi dari sensor_logs dalam range waktu.
func (r *AnalyticsPostgres) GetConditionBreakdown(ctx context.Context, userID string, from time.Time, to time.Time) ([]*repository.ConditionBreakdownEntry, error) {
	query := `
		SELECT
			condition,
			COUNT(*)::int AS count
		FROM sensor_logs
		WHERE user_id = $1
		  AND recorded_at >= $2
		  AND recorded_at < $3
		GROUP BY condition
		ORDER BY count DESC
	`
	rows, err := r.db.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*repository.ConditionBreakdownEntry
	for rows.Next() {
		e := &repository.ConditionBreakdownEntry{}
		if err := rows.Scan(&e.Condition, &e.Count); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

// GetTimeline — kondisi dominan per jam untuk satu hari (00:00 – 23:59).
func (r *AnalyticsPostgres) GetTimeline(ctx context.Context, userID string, date time.Time) ([]*repository.TimelineEntry, error) {
	// Ambil seluruh hari: date 00:00:00 sampai date+1 00:00:00
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	query := `
		SELECT
			EXTRACT(HOUR FROM recorded_at)::int AS hour,
			condition,
			COUNT(*)::int AS total_count
		FROM sensor_logs
		WHERE user_id = $1
		  AND recorded_at >= $2
		  AND recorded_at < $3
		GROUP BY hour, condition
		ORDER BY hour, total_count DESC
	`
	rows, err := r.db.Query(ctx, query, userID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Per jam, ambil kondisi dengan count tertinggi (dominant)
	type hourEntry struct {
		condition string
		count     int
		totalSec  int
	}
	dominantPerHour := make(map[int]*hourEntry)

	for rows.Next() {
		var hour, count int
		var condition string
		if err := rows.Scan(&hour, &condition, &count); err != nil {
			return nil, err
		}
		if existing, ok := dominantPerHour[hour]; !ok || count > existing.count {
			dominantPerHour[hour] = &hourEntry{condition: condition, count: count}
		}
	}

	var result []*repository.TimelineEntry
	for hour := 0; hour < 24; hour++ {
		e := &repository.TimelineEntry{Hour: hour}
		if entry, ok := dominantPerHour[hour]; ok {
			e.DominantCondition = entry.condition
			e.TotalSec = entry.count * 30 // estimasi: tiap log ~30 detik (heartbeat interval)
		}
		result = append(result, e)
	}
	return result, nil
}
