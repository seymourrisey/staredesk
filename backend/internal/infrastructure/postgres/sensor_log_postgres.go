package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seymourrisey/staredesk/internal/entity"
	"github.com/seymourrisey/staredesk/pkg/idgen"
)

type SensorLogPostgres struct {
	db *pgxpool.Pool
}

func NewSensorLogPostgres(db *pgxpool.Pool) *SensorLogPostgres {
	return &SensorLogPostgres{db: db}
}

func (r *SensorLogPostgres) Create(ctx context.Context, log *entity.SensorLog) error {
	log.ID = idgen.NewSensorLogID()

	query := `
		INSERT INTO sensor_logs (id, user_id, distance_cm, ldr_value, pir_detected, condition, log_type, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		log.ID,
		log.UserID,
		log.DistanceCM,
		log.LDRValue,
		log.PIRDetected,
		log.Condition,
		log.LogType,
		log.RecordedAt,
	)
	return err
}

func (r *SensorLogPostgres) GetBySessionID(ctx context.Context, sessionID string) ([]*entity.SensorLog, error) {
	// Fetch logs that fall within the session's time range.
	// We join against sessions to get started_at / ended_at.
	query := `
		SELECT sl.id, sl.user_id, sl.distance_cm, sl.ldr_value, sl.pir_detected,
		       sl.condition, sl.log_type, sl.recorded_at
		FROM sensor_logs sl
		JOIN sessions s ON sl.user_id = s.user_id
		WHERE s.id = $1
		  AND sl.recorded_at >= s.started_at
		  AND (s.ended_at IS NULL OR sl.recorded_at <= s.ended_at)
		ORDER BY sl.recorded_at ASC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.SensorLog
	for rows.Next() {
		l := &entity.SensorLog{}
		err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.DistanceCM,
			&l.LDRValue,
			&l.PIRDetected,
			&l.Condition,
			&l.LogType,
			&l.RecordedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
