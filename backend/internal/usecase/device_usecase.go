package usecase

import (
	"context"
	"time"

	"github.com/seymourrisey/staredesk/internal/entity"
	"github.com/seymourrisey/staredesk/internal/repository"
)

// Default threshold values — used when no config row exists in DB
const (
	DefaultDistanceMinCM      = 40
	DefaultDistanceMaxCM      = 90
	DefaultLDRThreshold       = 500
	DefaultAwayTimeoutMinutes = 5
)

type DeviceUsecase struct {
	configRepo repository.DeviceConfigRepository
	statusRepo repository.DeviceStatusRepository
}

// DeviceConfigPayload digunakan untuk transfer data config ke MQTT layer.
type DeviceConfigPayload struct {
	DistanceMinCM      int
	DistanceMaxCM      int
	LDRThreshold       int
	AwayTimeoutMinutes int
}

func NewDeviceUsecase(
	configRepo repository.DeviceConfigRepository,
	statusRepo repository.DeviceStatusRepository,
) *DeviceUsecase {
	return &DeviceUsecase{
		configRepo: configRepo,
		statusRepo: statusRepo,
	}
}

// GetConfig returns existing config or default values if no row exists.
func (u *DeviceUsecase) GetConfig(ctx context.Context, userID string) (*entity.DeviceConfig, error) {
	config, err := u.configRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return &entity.DeviceConfig{
			UserID:             userID,
			DistanceMinCM:      DefaultDistanceMinCM,
			DistanceMaxCM:      DefaultDistanceMaxCM,
			LDRThreshold:       DefaultLDRThreshold,
			AwayTimeoutMinutes: DefaultAwayTimeoutMinutes,
			ConfigAck:          true,
		}, nil
	}
	return config, nil
}

// UpdateConfig persists new thresholds and resets config_ack to false.
// Caller (handler) is responsible for publishing to MQTT after this returns.
func (u *DeviceUsecase) UpdateConfig(ctx context.Context, userID string, distanceMin, distanceMax, ldrThreshold, awayTimeout int) (*entity.DeviceConfig, error) {
	existing, err := u.configRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var config entity.DeviceConfig
	if existing != nil {
		config = *existing
	} else {
		config.UserID = userID
	}

	config.DistanceMinCM = distanceMin
	config.DistanceMaxCM = distanceMax
	config.LDRThreshold = ldrThreshold
	config.AwayTimeoutMinutes = awayTimeout
	config.ConfigAck = false // reset — waiting for ESP32 ack

	if err := u.configRepo.Upsert(ctx, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SetConfigAck updates config_ack flag — called by MQTT handler on config/ack topic.
func (u *DeviceUsecase) SetConfigAck(ctx context.Context, userID string, ack bool) error {
	return u.configRepo.SetConfigAck(ctx, userID, ack)
}

// GetAwayTimeoutMinutes returns away_timeout_minutes from DB, falls back to default.
func (u *DeviceUsecase) GetAwayTimeoutMinutes(ctx context.Context, userID string) (int, error) {
	config, err := u.configRepo.GetByUserID(ctx, userID)
	if err != nil {
		return DefaultAwayTimeoutMinutes, err
	}
	if config == nil {
		return DefaultAwayTimeoutMinutes, nil
	}
	return config.AwayTimeoutMinutes, nil
}

// GetStatus returns device online/offline status.
func (u *DeviceUsecase) GetStatus(ctx context.Context, userID string) (*entity.DeviceStatus, error) {
	status, err := u.statusRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if status == nil {
		return &entity.DeviceStatus{
			UserID:   userID,
			IsOnline: false,
			LastSeen: nil,
		}, nil
	}
	return status, nil
}

// UpdateStatus upserts device online/offline state — called by MQTT handler on device/status topic.
func (u *DeviceUsecase) UpdateStatus(ctx context.Context, userID string, isOnline bool) error {
	existing, err := u.statusRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var status entity.DeviceStatus
	if existing != nil {
		status = *existing
	} else {
		status.UserID = userID
	}

	status.IsOnline = isOnline
	if isOnline {
		now := time.Now().UTC()
		status.LastSeen = &now
	}

	return u.statusRepo.Upsert(ctx, &status)
}
