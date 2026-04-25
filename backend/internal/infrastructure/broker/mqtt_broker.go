package broker

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient wraps paho MQTT client dengan reconnect logic.
type MQTTClient struct {
	client mqtt.Client
}

// NewMQTTClient membuat koneksi ke HiveMQ Cloud dengan TLS.
// broker contoh: "abc123.s1.eu.hivemq.cloud"
// port: 8883
func NewMQTTClient(broker, clientID, username, password string, port int) (*MQTTClient, error) {
	brokerURL := fmt.Sprintf("ssl://%s:%d", broker, port)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)

	// TLS — HiveMQ Cloud pakai certificate resmi, tidak perlu custom CA
	opts.SetTLSConfig(&tls.Config{
		InsecureSkipVerify: false,
	})

	// Reconnect otomatis
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		log.Println("[MQTT] Connected to broker")
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("[MQTT] Connection lost: %v — reconnecting...", err)
	})
	opts.SetReconnectingHandler(func(_ mqtt.Client, _ *mqtt.ClientOptions) {
		log.Println("[MQTT] Reconnecting...")
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return nil, fmt.Errorf("MQTT connect timeout")
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("MQTT connect error: %w", err)
	}

	return &MQTTClient{client: client}, nil
}

// Subscribe mendaftarkan handler ke satu topic.
func (m *MQTTClient) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	token := m.client.Subscribe(topic, qos, handler)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("subscribe timeout: %s", topic)
	}
	return token.Error()
}

// Publish mengirim pesan ke satu topic.
func (m *MQTTClient) Publish(topic string, qos byte, retained bool, payload []byte) error {
	token := m.client.Publish(topic, qos, retained, payload)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("publish timeout: %s", topic)
	}
	return token.Error()
}

// Disconnect menutup koneksi dengan grace period.
func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
	log.Println("[MQTT] Disconnected")
}
