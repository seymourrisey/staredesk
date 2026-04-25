package mqttdelivery

import "fmt"

// QoS constants
const (
	QoS0 = byte(0)
	QoS1 = byte(1)
)

// Topic helpers — format topic dengan user_id
func TopicTelemetry(userID string) string {
	return fmt.Sprintf("study/%s/device/telemetry", userID)
}

func TopicStatus(userID string) string {
	return fmt.Sprintf("study/%s/device/status", userID)
}

func TopicConfig(userID string) string {
	return fmt.Sprintf("study/%s/device/config", userID)
}

func TopicConfigAck(userID string) string {
	return fmt.Sprintf("study/%s/device/config/ack", userID)
}
