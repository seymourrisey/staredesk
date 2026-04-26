package repository

import (
	"context"
	"github.com/seymourrisey/staredesk/internal/entity"
)

type DeviceConfigRepository interface {
	GetByUserID(ctx context.Context, userID string) (*entity.DeviceConfig, error)
	Upsert(ctx context.Context, config *entity.DeviceConfig) error
	SetConfigAck(ctx context.Context, userID string, ack bool) error
}

type DeviceStatusRepository interface {
	GetByUserID(ctx context.Context, userID string) (*entity.DeviceStatus, error)
	Upsert(ctx context.Context, status *entity.DeviceStatus) error
}
