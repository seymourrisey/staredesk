package repository

import (
	"context"
	"time"

	"github.com/seymourrisey/staredesk/internal/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	Update(ctx context.Context, session *entity.Session) error
	GetActiveByUserID(ctx context.Context, userID string) (*entity.Session, error)

	// Untuk REST API
	GetByID(ctx context.Context, id string, userID string) (*entity.Session, error)
	GetAll(ctx context.Context, userID string, limit int, offset int) ([]*entity.Session, int, error)
	GetByDateRange(ctx context.Context, userID string, from time.Time, to time.Time) ([]*entity.Session, error)
}
