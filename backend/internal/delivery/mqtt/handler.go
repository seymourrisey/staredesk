package mqttdelivery

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	websocketdelivery "github.com/seymourrisey/staredesk/internal/delivery/websocket"
)

// MakeTelemetryHandler mengembalikan handler untuk topic device/telemetry.
// Payload MQTT di-forward langsung ke WebSocket broadcast.
func MakeTelemetryHandler(hub *websocketdelivery.Hub) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		log.Printf("[MQTT] telemetry | topic: %s | payload: %s", msg.Topic(), string(msg.Payload()))
		hub.Broadcast <- msg.Payload()
	}
}

// MakeStatusHandler mengembalikan handler untuk topic device/status.
func MakeStatusHandler(hub *websocketdelivery.Hub) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		log.Printf("[MQTT] status    | topic: %s | payload: %s", msg.Topic(), string(msg.Payload()))
		hub.Broadcast <- msg.Payload()
	}
}

// MakeConfigAckHandler mengembalikan handler untuk topic device/config/ack.
func MakeConfigAckHandler(hub *websocketdelivery.Hub) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		log.Printf("[MQTT] config/ack | topic: %s | payload: %s", msg.Topic(), string(msg.Payload()))
		hub.Broadcast <- msg.Payload()
	}
}
