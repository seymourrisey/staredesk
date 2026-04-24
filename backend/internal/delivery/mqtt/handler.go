package mqttdelivery

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HandleTelemetry dipanggil saat ada pesan masuk di topic device/telemetry.
// Untuk sekarang: log raw payload. Business logic ditambah di room berikutnya.
func HandleTelemetry(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] telemetry | topic: %s | payload: %s", msg.Topic(), string(msg.Payload()))
}

// HandleStatus dipanggil saat ada pesan masuk di topic device/status.
func HandleStatus(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] status    | topic: %s | payload: %s", msg.Topic(), string(msg.Payload()))
}

// HandleConfigAck dipanggil saat ESP32 konfirmasi config diterima.
func HandleConfigAck(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] config/ack | topic: %s | payload: %s", msg.Topic(), string(msg.Payload()))
}
