package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mqttdelivery "github.com/seymourrisey/staredesk/internal/delivery/mqtt"
	"github.com/seymourrisey/staredesk/internal/usecase"
)

type DeviceHandler struct {
	deviceUC *usecase.DeviceUsecase
	mqttSub  *mqttdelivery.Subscriber
}

// DeviceConfigPayload digunakan untuk transfer data config ke MQTT layer.
type DeviceConfigPayload struct {
	DistanceMinCM      int
	DistanceMaxCM      int
	LDRThreshold       int
	AwayTimeoutMinutes int
}

func NewDeviceHandler(deviceUC *usecase.DeviceUsecase, mqttSub *mqttdelivery.Subscriber) *DeviceHandler {
	return &DeviceHandler{
		deviceUC: deviceUC,
		mqttSub:  mqttSub,
	}
}

// GET /device/config
func (h *DeviceHandler) GetConfig(c *gin.Context) {
	userID := c.GetString("user_id")

	config, err := h.deviceUC.GetConfig(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                   config.ID,
		"distance_min_cm":      config.DistanceMinCM,
		"distance_max_cm":      config.DistanceMaxCM,
		"ldr_threshold":        config.LDRThreshold,
		"away_timeout_minutes": config.AwayTimeoutMinutes,
		"config_ack":           config.ConfigAck,
		"updated_at":           config.UpdatedAt,
	})
}

// PUT /device/config
func (h *DeviceHandler) UpdateConfig(c *gin.Context) {
	userID := c.GetString("user_id")

	var body struct {
		DistanceMinCM      int `json:"distance_min_cm" binding:"required,min=1"`
		DistanceMaxCM      int `json:"distance_max_cm" binding:"required,min=1"`
		LDRThreshold       int `json:"ldr_threshold" binding:"required,min=0"`
		AwayTimeoutMinutes int `json:"away_timeout_minutes" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.deviceUC.UpdateConfig(
		c.Request.Context(),
		userID,
		body.DistanceMinCM,
		body.DistanceMaxCM,
		body.LDRThreshold,
		body.AwayTimeoutMinutes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update config"})
		return
	}

	// Publish config baru ke ESP32 via MQTT
	mqttPayload := &usecase.DeviceConfigPayload{
		DistanceMinCM:      config.DistanceMinCM,
		DistanceMaxCM:      config.DistanceMaxCM,
		LDRThreshold:       config.LDRThreshold,
		AwayTimeoutMinutes: config.AwayTimeoutMinutes,
	}
	if err := h.mqttSub.PublishConfig(mqttPayload); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":                   config.ID,
			"distance_min_cm":      config.DistanceMinCM,
			"distance_max_cm":      config.DistanceMaxCM,
			"ldr_threshold":        config.LDRThreshold,
			"away_timeout_minutes": config.AwayTimeoutMinutes,
			"config_ack":           config.ConfigAck,
			"updated_at":           config.UpdatedAt,
			"mqtt_publish":         "failed — config saved but ESP32 not notified",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                   config.ID,
		"distance_min_cm":      config.DistanceMinCM,
		"distance_max_cm":      config.DistanceMaxCM,
		"ldr_threshold":        config.LDRThreshold,
		"away_timeout_minutes": config.AwayTimeoutMinutes,
		"config_ack":           config.ConfigAck,
		"updated_at":           config.UpdatedAt,
		"mqtt_publish":         "ok",
	})
}

// GET /device/status
func (h *DeviceHandler) GetStatus(c *gin.Context) {
	userID := c.GetString("user_id")

	status, err := h.deviceUC.GetStatus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_online": status.IsOnline,
		"last_seen": status.LastSeen,
	})
}
