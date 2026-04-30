package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seymourrisey/staredesk/internal/usecase"
)

type AnalyticsHandler struct {
	analyticsUsecase *usecase.AnalyticsUsecase
}

func NewAnalyticsHandler(analyticsUsecase *usecase.AnalyticsUsecase) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsUsecase: analyticsUsecase}
}

// GET /analytics/peak-hours?range=week|month
func (h *AnalyticsHandler) GetPeakHours(c *gin.Context) {
	userID, _ := c.Get("user_id")
	rangeParam := c.DefaultQuery("range", "week")

	result, err := h.analyticsUsecase.GetPeakHours(c.Request.Context(), userID.(string), rangeParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"range":   result.Range,
		"entries": result.Entries,
	})
}

// GET /analytics/condition-breakdown?range=today|week|month
func (h *AnalyticsHandler) GetConditionBreakdown(c *gin.Context) {
	userID, _ := c.Get("user_id")
	rangeParam := c.DefaultQuery("range", "today")

	result, err := h.analyticsUsecase.GetConditionBreakdown(c.Request.Context(), userID.(string), rangeParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"range":   result.Range,
		"entries": result.Entries,
	})
}

// GET /analytics/timeline?date=YYYY-MM-DD
func (h *AnalyticsHandler) GetTimeline(c *gin.Context) {
	userID, _ := c.Get("user_id")
	dateParam := c.Query("date")

	result, err := h.analyticsUsecase.GetTimeline(c.Request.Context(), userID.(string), dateParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":    result.Date,
		"entries": result.Entries,
	})
}
