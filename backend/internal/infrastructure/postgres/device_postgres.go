package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seymourrisey/staredesk/internal/entity"
	"github.com/seymourrisey/staredesk/pkg/idgen"
)

// --- DeviceConfig ---

type DeviceConfigPostgres struct {
	db *pgxpool.Pool
}

func NewDeviceConfigPostgres(db *pgxpool.Pool) *DeviceConfigPostgres {
	return &DeviceConfigPostgres{db: db}
}

func (r *DeviceConfigPostgres) GetByUserID(ctx context.Context, userID string) (*entity.DeviceConfig, error) {
	query := `
		SELECT id, user_id, distance_min_cm, distance_max_cm, ldr_threshold, away_timeout_minutes, config_ack, updated_at
		FROM device_config
		WHERE user_id = $1
		LIMIT 1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	config, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[entity.DeviceConfig])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return config, nil
}

func (r *DeviceConfigPostgres) Upsert(ctx context.Context, config *entity.DeviceConfig) error {
	if config.ID == "" {
		config.ID = idgen.NewDeviceConfigID()
	}
	query := `
		INSERT INTO device_config (id, user_id, distance_min_cm, distance_max_cm, ldr_threshold, away_timeout_minutes, config_ack, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET
			distance_min_cm      = EXCLUDED.distance_min_cm,
			distance_max_cm      = EXCLUDED.distance_max_cm,
			ldr_threshold        = EXCLUDED.ldr_threshold,
			away_timeout_minutes = EXCLUDED.away_timeout_minutes,
			config_ack           = EXCLUDED.config_ack,
			updated_at           = NOW()`

	_, err := r.db.Exec(ctx, query,
		config.ID,
		config.UserID,
		config.DistanceMinCM,
		config.DistanceMaxCM,
		config.LDRThreshold,
		config.AwayTimeoutMinutes,
		config.ConfigAck,
	)
	return err
}

func (r *DeviceConfigPostgres) SetConfigAck(ctx context.Context, userID string, ack bool) error {
	query := `UPDATE device_config SET config_ack = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, ack, userID)
	return err
}

// --- DeviceStatus ---

type DeviceStatusPostgres struct {
	db *pgxpool.Pool
}

func NewDeviceStatusPostgres(db *pgxpool.Pool) *DeviceStatusPostgres {
	return &DeviceStatusPostgres{db: db}
}

func (r *DeviceStatusPostgres) GetByUserID(ctx context.Context, userID string) (*entity.DeviceStatus, error) {
	query := `
		SELECT id, user_id, is_online, last_seen
		FROM device_status
		WHERE user_id = $1
		LIMIT 1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	status, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[entity.DeviceStatus])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return status, nil
}

func (r *DeviceStatusPostgres) Upsert(ctx context.Context, status *entity.DeviceStatus) error {
	if status.ID == "" {
		status.ID = idgen.NewDeviceStatusID()
	}
	query := `
		INSERT INTO device_status (id, user_id, is_online, last_seen)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			is_online = EXCLUDED.is_online,
			last_seen = EXCLUDED.last_seen`

	_, err := r.db.Exec(ctx, query,
		status.ID,
		status.UserID,
		status.IsOnline,
		status.LastSeen,
	)
	return err
}
