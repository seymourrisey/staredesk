package mqttdelivery

import (
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
}

// NewSubscriber membuat Subscriber baru.
func NewSubscriber(
	client *broker.MQTTClient,
	userID string,
	hub *websocketdelivery.Hub,
	sessionUC *usecase.SessionUsecase,
	sensorUC *usecase.SensorUsecase,
) *Subscriber {
	return &Subscriber{
		client:    client,
		userID:    userID,
		hub:       hub,
		sessionUC: sessionUC,
		sensorUC:  sensorUC,
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
	if err := s.client.Subscribe(t, QoS1, MakeStatusHandler(s.hub)); err != nil {
		return err
	}
	log.Printf("[MQTT] Subscribed: %s", t)

	// device/config/ack — QoS 1
	t = TopicConfigAck(s.userID)
	if err := s.client.Subscribe(t, QoS1, MakeConfigAckHandler(s.hub)); err != nil {
		return err
	}
	log.Printf("[MQTT] Subscribed: %s", t)

	return nil
}
