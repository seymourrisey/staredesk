package postgres

import (
	"context"
	"errors"

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
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.StartedAt,
		&s.EndedAt,
		&s.DurationSec,
		&s.DominantCondition,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // tidak ada sesi aktif
		}
		return nil, err
	}
	return s, nil
}
