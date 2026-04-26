package repository

import (
	"context"
	"time"

	"github.com/seymourrisey/staredesk/internal/entity"
)

type SensorLogRepository interface {
	Create(ctx context.Context, log *entity.SensorLog) error
	GetByDateRange(ctx context.Context, userID string, from time.Time, to time.Time, limit int) ([]*entity.SensorLog, error)
}
