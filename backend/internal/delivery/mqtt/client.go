package mqttdelivery

import (
	"encoding/json"
	"fmt"
	"log"

	websocketdelivery "github.com/seymourrisey/staredesk/internal/delivery/websocket"
	"github.com/seymourrisey/staredesk/internal/infrastructure/broker"
	"github.com/seymourrisey/staredesk/internal/usecase"
)

// Subscriber memegang referensi ke MQTTClient, userID, WebSocket Hub, dan usecase.
type Subscriber struct {
	client    *broker.MQTTClient
	userID    string
	hub       *websocketdelivery.Hub
	sessionUC *usecase.SessionUsecase
	sensorUC  *usecase.SensorUsecase
	deviceUC  *usecase.DeviceUsecase
}

// NewSubscriber membuat Subscriber baru.
func NewSubscriber(
	client *broker.MQTTClient,
	userID string,
	hub *websocketdelivery.Hub,
	sessionUC *usecase.SessionUsecase,
	sensorUC *usecase.SensorUsecase,
	deviceUC *usecase.DeviceUsecase,
) *Subscriber {
	return &Subscriber{
		client:    client,
		userID:    userID,
		hub:       hub,
		sessionUC: sessionUC,
		sensorUC:  sensorUC,
		deviceUC:  deviceUC,
	}
}

// SubscribeAll mendaftarkan handler ke semua topic yang perlu di-listen.
func (s *Subscriber) SubscribeAll() error {
	// telemetry — QoS 0
	t := TopicTelemetry(s.userID)
	if err := s.client.Subscribe(t, QoS0, MakeTelemetryHandler(s.hub, s.sessionUC, s.sensorUC, s.userID)); err != nil {
		return err
	}
	log.Printf("[MQTT] Subscribed: %s", t)

	// device/status — QoS 1
	t = TopicStatus(s.userID)
	if err := s.client.Subscribe(t, QoS1, MakeStatusHandler(s.hub, s.deviceUC, s.userID)); err != nil {
		return err
	}
	log.Printf("[MQTT] Subscribed: %s", t)

	// device/config/ack — QoS 1
	t = TopicConfigAck(s.userID)
	if err := s.client.Subscribe(t, QoS1, MakeConfigAckHandler(s.hub, s.deviceUC, s.userID)); err != nil {
		return err
	}
	log.Printf("[MQTT] Subscribed: %s", t)

	return nil
}

// configMQTTPayload adalah struktur JSON yang dikirim ke ESP32 saat config update.
type configMQTTPayload struct {
	DistanceMinCM      int `json:"distance_min_cm"`
	DistanceMaxCM      int `json:"distance_max_cm"`
	LDRThreshold       int `json:"ldr_threshold"`
	AwayTimeoutMinutes int `json:"away_timeout_minutes"`
}

// PublishConfig mengirim threshold config terbaru ke ESP32 via MQTT.
func (s *Subscriber) PublishConfig(config *usecase.DeviceConfigPayload) error {
	payload := configMQTTPayload{
		DistanceMinCM:      config.DistanceMinCM,
		DistanceMaxCM:      config.DistanceMaxCM,
		LDRThreshold:       config.LDRThreshold,
		AwayTimeoutMinutes: config.AwayTimeoutMinutes,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal config payload: %w", err)
	}
	topic := TopicConfig(s.userID)
	if err := s.client.Publish(topic, QoS1, false, data); err != nil {
		return fmt.Errorf("failed to publish config: %w", err)
	}
	log.Printf("[MQTT] Config published to %s", topic)
	return nil
}
