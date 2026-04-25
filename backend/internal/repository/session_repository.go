package repository

import (
	"context"

	"github.com/seymourrisey/staredesk/internal/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	Update(ctx context.Context, session *entity.Session) error
	GetActiveByUserID(ctx context.Context, userID string) (*entity.Session, error)
}
