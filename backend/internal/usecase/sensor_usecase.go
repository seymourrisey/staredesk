package usecase

import (
	"context"
	"time"

	"github.com/seymourrisey/staredesk/internal/entity"
	"github.com/seymourrisey/staredesk/internal/repository"
)

type SensorUsecase struct {
	sensorLogRepo repository.SensorLogRepository
}

func NewSensorUsecase(sensorLogRepo repository.SensorLogRepository) *SensorUsecase {
	return &SensorUsecase{sensorLogRepo: sensorLogRepo}
}

// Save menyimpan satu entri sensor log ke DB.
// logType: "heartbeat" | "condition_change"
func (u *SensorUsecase) Create(ctx context.Context, userID string, payload *SensorPayload, logType string) error {
	log := &entity.SensorLog{
		UserID:      userID,
		DistanceCM:  payload.DistanceCM,
		LDRValue:    payload.LDRValue,
		PIRDetected: payload.PIRDetected,
		Condition:   payload.Condition,
		LogType:     logType,
		RecordedAt:  time.Now(),
	}
	return u.sensorLogRepo.Create(ctx, log)
}

// SensorPayload adalah DTO dari MQTT telemetry yang sudah di-parse.
type SensorPayload struct {
	DistanceCM  *float64 `json:"distance_cm"`
	LDRValue    *int     `json:"ldr_value"`
	PIRDetected bool     `json:"pir_detected"`
	Condition   string   `json:"condition"`
	LogType     string   `json:"log_type"`
}
