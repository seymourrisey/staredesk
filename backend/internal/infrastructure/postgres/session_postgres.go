package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seymourrisey/staredesk/internal/entity"
	"github.com/seymourrisey/staredesk/pkg/idgen"
)

type SessionPostgres struct {
	db *pgxpool.Pool
}

func NewSessionPostgres(db *pgxpool.Pool) *SessionPostgres {
	return &SessionPostgres{db: db}
}

func (r *SessionPostgres) Create(ctx context.Context, session *entity.Session) error {
	session.ID = idgen.NewSessionID()
	query := `
		INSERT INTO sessions (id, user_id, started_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, session.ID, session.UserID, session.StartedAt)
	return err
}

func (r *SessionPostgres) Update(ctx context.Context, session *entity.Session) error {
	query := `
		UPDATE sessions
		SET ended_at = $1, duration_sec = $2, dominant_condition = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query,
		session.EndedAt,
		session.DurationSec,
		session.DominantCondition,
		session.ID,
	)
	return err
}

func (r *SessionPostgres) GetActiveByUserID(ctx context.Context, userID string) (*entity.Session, error) {
	query := `
		SELECT id, user_id, started_at, ended_at, duration_sec, dominant_condition
		FROM sessions
		WHERE user_id = $1 AND ended_at IS NULL
		ORDER BY started_at DESC
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, query, userID)
	s := &entity.Session{}
	err := row.Scan(&s.ID, &s.UserID, &s.StartedAt, &s.EndedAt, &s.DurationSec, &s.DominantCondition)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *SessionPostgres) GetByID(ctx context.Context, id string, userID string) (*entity.Session, error) {
	query := `
		SELECT id, user_id, started_at, ended_at, duration_sec, dominant_condition
		FROM sessions
		WHERE id = $1 AND user_id = $2
	`
	row := r.db.QueryRow(ctx, query, id, userID)
	s := &entity.Session{}
	err := row.Scan(&s.ID, &s.UserID, &s.StartedAt, &s.EndedAt, &s.DurationSec, &s.DominantCondition)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *SessionPostgres) GetAll(ctx context.Context, userID string, limit int, offset int) ([]*entity.Session, int, error) {
	// Total count
	var total int
	countQuery := `SELECT COUNT(*) FROM sessions WHERE user_id = $1 AND ended_at IS NOT NULL`
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data
	query := `
		SELECT id, user_id, started_at, ended_at, duration_sec, dominant_condition
		FROM sessions
		WHERE user_id = $1 AND ended_at IS NOT NULL
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []*entity.Session
	for rows.Next() {
		s := &entity.Session{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.StartedAt, &s.EndedAt, &s.DurationSec, &s.DominantCondition); err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, s)
	}
	return sessions, total, nil
}

func (r *SessionPostgres) GetByDateRange(ctx context.Context, userID string, from time.Time, to time.Time) ([]*entity.Session, error) {
	query := `
		SELECT id, user_id, started_at, ended_at, duration_sec, dominant_condition
		FROM sessions
		WHERE user_id = $1
		  AND ended_at IS NOT NULL
		  AND started_at >= $2
		  AND started_at < $3
		ORDER BY started_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*entity.Session
	for rows.Next() {
		s := &entity.Session{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.StartedAt, &s.EndedAt, &s.DurationSec, &s.DominantCondition); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}
