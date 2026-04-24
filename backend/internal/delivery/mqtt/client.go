package mqttdelivery

import (
	"log"

	"github.com/seymourrisey/staredesk/internal/infrastructure/broker"
)

// Subscriber memegang referensi ke MQTTClient dan userID.
type Subscriber struct {
	client *broker.MQTTClient
	userID string
}

// NewSubscriber membuat Subscriber baru.
func NewSubscriber(client *broker.MQTTClient, userID string) *Subscriber {
	return &Subscriber{
		client: client,
		userID: userID,
	}
}

// SubscribeAll mendaftarkan handler ke semua topic yang perlu di-listen.
func (s *Subscriber) SubscribeAll() error {
	topics := []struct {
		topic   string
		qos     byte
		handler func(_ interface{}, msg interface{})
	}{}
	_ = topics

	// telemetry — QoS 0
	t := TopicTelemetry(s.userID)
	if err := s.client.Subscribe(t, QoS0, HandleTelemetry); err != nil {
		return err
	}
	log.Printf("[MQTT] Subscribed: %s", t)

	// device/status — QoS 1
	t = TopicStatus(s.userID)
	if err := s.client.Subscribe(t, QoS1, HandleStatus); err != nil {
		return err
	}
	log.Printf("[MQTT] Subscribed: %s", t)

	// device/config/ack — QoS 1
	t = TopicConfigAck(s.userID)
	if err := s.client.Subscribe(t, QoS1, HandleConfigAck); err != nil {
		return err
	}
	log.Printf("[MQTT] Subscribed: %s", t)

	return nil
}
