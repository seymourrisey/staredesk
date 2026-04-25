package repository

import (
	"context"

	"github.com/seymourrisey/staredesk/internal/entity"
)

type SensorLogRepository interface {
	Create(ctx context.Context, log *entity.SensorLog) error
}
