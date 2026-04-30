package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seymourrisey/staredesk/internal/usecase"
)

type SensorHandler struct {
	sensorUsecase *usecase.SensorUsecase
}

func NewSensorHandler(sensorUsecase *usecase.SensorUsecase) *SensorHandler {
	return &SensorHandler{sensorUsecase: sensorUsecase}
}

// GET /sensor-logs?from=&to=&limit=
func (h *SensorHandler) GetLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	var from, to time.Time
	var err error

	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' format, use RFC3339"})
			return
		}
	} else {
		from = time.Now().Add(-24 * time.Hour)
	}

	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' format, use RFC3339"})
			return
		}
	} else {
		to = time.Now()
	}

	logs, err := h.sensorUsecase.GetLogs(c.Request.Context(), userID.(string), from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"from":  from.UTC().Format(time.RFC3339),
		"to":    to.UTC().Format(time.RFC3339),
		"limit": limit,
		"logs":  logs,
	})
}
