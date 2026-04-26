package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seymourrisey/staredesk/internal/usecase"
)

type SessionHandler struct {
	sessionUsecase *usecase.SessionUsecase
}

func NewSessionHandler(sessionUsecase *usecase.SessionUsecase) *SessionHandler {
	return &SessionHandler{sessionUsecase: sessionUsecase}
}

// GET /sessions?limit=20&offset=0
func (h *SessionHandler) GetAll(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result, err := h.sessionUsecase.GetAll(c.Request.Context(), userID.(string), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": result.Sessions,
		"total":    result.Total,
		"limit":    result.Limit,
		"offset":   result.Offset,
	})
}

// GET /sessions/summary?range=today|week|month
func (h *SessionHandler) GetSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	rangeParam := c.DefaultQuery("range", "today")

	result, err := h.sessionUsecase.GetSummary(c.Request.Context(), userID.(string), rangeParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"range":         result.Range,
		"total_sec":     result.TotalSec,
		"session_count": result.SessionCount,
		"sessions":      result.Sessions,
	})
}

// GET /sessions/:id
func (h *SessionHandler) GetByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	session, err := h.sessionUsecase.GetByID(c.Request.Context(), id, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}
