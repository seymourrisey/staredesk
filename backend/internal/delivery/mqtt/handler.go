package mqttdelivery

import (
	"context"
	"encoding/json"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	websocketdelivery "github.com/seymourrisey/staredesk/internal/delivery/websocket"
	"github.com/seymourrisey/staredesk/internal/usecase"
)

// mqttTelemetryPayload adalah struktur JSON yang dikirim ESP32.
type mqttTelemetryPayload struct {
	DistanceCM  *float64 `json:"distance_cm"`
	LDRValue    *int     `json:"ldr_value"`
	PIRDetected bool     `json:"pir_detected"`
	Condition   string   `json:"condition"`
	LogType     string   `json:"log_type"` // "heartbeat" | "condition_change"
	Timestamp   string   `json:"timestamp"`
}

// mqttStatusPayload adalah struktur JSON device/status dari ESP32.
type mqttStatusPayload struct {
	IsOnline bool `json:"is_online"`
}

// wsPayload adalah struktur yang di-broadcast ke WebSocket clients.
type wsPayload struct {
	Type      string         `json:"type"`
	Timestamp string         `json:"timestamp"`
	Device    wsDeviceInfo   `json:"device"`
	Sensors   *wsSensorInfo  `json:"sensors,omitempty"`
	Condition string         `json:"condition,omitempty"`
	Session   *wsSessionInfo `json:"session,omitempty"`
}

type wsDeviceInfo struct {
	IsOnline bool   `json:"is_online"`
	LastSeen string `json:"last_seen"`
}

type wsSensorInfo struct {
	DistanceCM  *float64 `json:"distance_cm"`
	LDRValue    *int     `json:"ldr_value"`
	PIRDetected bool     `json:"pir_detected"`
}

type wsSessionInfo struct {
	IsActive  bool   `json:"is_active"`
	StartedAt string `json:"started_at,omitempty"`
}

// MakeTelemetryHandler returns MQTT message handler untuk topic device/telemetry.
func MakeTelemetryHandler(
	hub *websocketdelivery.Hub,
	sessionUC *usecase.SessionUsecase,
	sensorUC *usecase.SensorUsecase,
	userID string,
) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var payload mqttTelemetryPayload
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Printf("[mqtt] telemetry parse error: %v", err)
			return
		}

		ctx := context.Background()
		now := time.Now().UTC()

		// 1. Process session logic — terima event
		event, err := sessionUC.ProcessCondition(ctx, userID, payload.Condition, now)
		if err != nil {
			log.Printf("[mqtt] session process error: %v", err)
		}

		// 2. Broadcast session_start jika ada
		if event == usecase.SessionEventStart {
			wsSessionStart := map[string]interface{}{
				"type":       "session_start",
				"timestamp":  now.Format(time.RFC3339),
				"started_at": now.Format(time.RFC3339),
			}
			if data, err := json.Marshal(wsSessionStart); err == nil {
				hub.Broadcast <- data
				log.Printf("[mqtt] session_start broadcasted")
			}
		}

		// 3. Broadcast session_end jika ada
		if event == usecase.SessionEventEnd {
			wsSessionEnd := map[string]interface{}{
				"type":      "session_end",
				"timestamp": now.Format(time.RFC3339),
			}
			if data, err := json.Marshal(wsSessionEnd); err == nil {
				hub.Broadcast <- data
				log.Printf("[mqtt] session_end broadcasted")
			}
		}

		// 4. Persist sensor log
		sensorPayload := &usecase.SensorPayload{
			DistanceCM:  payload.DistanceCM,
			LDRValue:    payload.LDRValue,
			PIRDetected: payload.PIRDetected,
			Condition:   payload.Condition,
			LogType:     payload.LogType,
		}
		logType := payload.LogType
		if logType == "" {
			logType = "heartbeat"
		}
		if err := sensorUC.Create(ctx, userID, sensorPayload, logType); err != nil {
			log.Printf("[mqtt] sensor log error: %v", err)
		}

		sessionInfo := &wsSessionInfo{
			IsActive: sessionUC.IsSessionActive(),
		}
		if sessionInfo.IsActive {
			if t := sessionUC.ActiveSessionStartedAt(); t != nil {
				sessionInfo.StartedAt = t.Format(time.RFC3339)
			}
		}

		// 5. Build & broadcast telemetry WebSocket payload
		wsMsg := wsPayload{
			Type:      "telemetry",
			Timestamp: now.Format(time.RFC3339),
			Device: wsDeviceInfo{
				IsOnline: true,
				LastSeen: now.Format(time.RFC3339),
			},
			Sensors: &wsSensorInfo{
				DistanceCM:  payload.DistanceCM,
				LDRValue:    payload.LDRValue,
				PIRDetected: payload.PIRDetected,
			},
			Condition: payload.Condition,
			Session:   sessionInfo,
		}

		if payload.LogType == "condition_change" {
			wsMsg.Type = "condition_change"
		}

		data, err := json.Marshal(wsMsg)
		if err != nil {
			log.Printf("[mqtt] ws marshal error: %v", err)
			return
		}
		hub.Broadcast <- data
	}
}

// MakeStatusHandler returns MQTT message handler untuk topic device/status.
func MakeStatusHandler(
	hub *websocketdelivery.Hub,
	deviceUC *usecase.DeviceUsecase,
	sessionUC *usecase.SessionUsecase,
	userID string,
) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var payload mqttStatusPayload
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Printf("[mqtt] status parse error: %v", err)
			return
		}

		ctx := context.Background()
		now := time.Now().UTC()

		// Update device_status di DB
		if err := deviceUC.UpdateStatus(ctx, userID, payload.IsOnline); err != nil {
			log.Printf("[mqtt] update device status error: %v", err)
		}

		// Device offline => force end active session
		if !payload.IsOnline && sessionUC.IsSessionActive() {
			log.Printf("[mqtt] device offline — force ending active session")
			if err := sessionUC.ForceEndSession(ctx, userID, now); err != nil {
				log.Printf("[mqtt] force end session error: %v", err)
			} else {
				// Broadcast session_end
				wsSessionEnd := map[string]interface{}{
					"type":      "session_end",
					"timestamp": now.Format(time.RFC3339),
				}
				if data, err := json.Marshal(wsSessionEnd); err == nil {
					hub.Broadcast <- data
				}
			}
		}

		wsMsg := wsPayload{
			Type:      "device_status",
			Timestamp: now.Format(time.RFC3339),
			Device: wsDeviceInfo{
				IsOnline: payload.IsOnline,
				LastSeen: now.Format(time.RFC3339),
			},
		}

		data, err := json.Marshal(wsMsg)
		if err != nil {
			log.Printf("[mqtt] ws marshal error: %v", err)
			return
		}
		hub.Broadcast <- data
	}
}

// MakeConfigAckHandler returns MQTT message handler untuk topic device/config/ack.
func MakeConfigAckHandler(
	hub *websocketdelivery.Hub,
	deviceUC *usecase.DeviceUsecase,
	userID string,
) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		// Update config_ack = true di DB
		ctx := context.Background()
		if err := deviceUC.SetConfigAck(ctx, userID, true); err != nil {
			log.Printf("[mqtt] set config ack error: %v", err)
		}

		now := time.Now().UTC().Format(time.RFC3339)
		wsMsg := map[string]string{
			"type":      "config_ack",
			"timestamp": now,
		}

		data, err := json.Marshal(wsMsg)
		if err != nil {
			log.Printf("[mqtt] ws marshal error: %v", err)
			return
		}
		hub.Broadcast <- data
	}
}
