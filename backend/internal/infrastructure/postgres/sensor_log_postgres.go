package postgres

import (
	"context"
	"time"

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

func (r *SensorLogPostgres) GetByDateRange(ctx context.Context, userID string, from time.Time, to time.Time, limit int) ([]*entity.SensorLog, error) {
	query := `
		SELECT id, user_id, distance_cm, ldr_value, pir_detected, condition, log_type, recorded_at
		FROM sensor_logs
		WHERE user_id = $1
		  AND recorded_at >= $2
		  AND recorded_at < $3
		ORDER BY recorded_at DESC
		LIMIT $4
	`
	rows, err := r.db.Query(ctx, query, userID, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.SensorLog
	for rows.Next() {
		l := &entity.SensorLog{}
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.DistanceCM, &l.LDRValue,
			&l.PIRDetected, &l.Condition, &l.LogType, &l.RecordedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
