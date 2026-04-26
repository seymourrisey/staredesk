package http

import (
	"github.com/gin-gonic/gin"
	"github.com/seymourrisey/staredesk/internal/delivery/http/handler"
	"github.com/seymourrisey/staredesk/internal/delivery/http/middleware"
	websocketdelivery "github.com/seymourrisey/staredesk/internal/delivery/websocket"
)

type Router struct {
	authHandler      *handler.AuthHandler
	deviceHandler    *handler.DeviceHandler
	sessionHandler   *handler.SessionHandler
	analyticsHandler *handler.AnalyticsHandler
	sensorHandler    *handler.SensorHandler
	wsHandler        *websocketdelivery.Handler
	jwtSecret        string
}

func NewRouter(
	authHandler *handler.AuthHandler,
	deviceHandler *handler.DeviceHandler,
	sessionHandler *handler.SessionHandler,
	analyticsHandler *handler.AnalyticsHandler,
	sensorHandler *handler.SensorHandler,
	wsHandler *websocketdelivery.Handler,
	jwtSecret string,
) *Router {
	return &Router{
		authHandler:      authHandler,
		deviceHandler:    deviceHandler,
		sessionHandler:   sessionHandler,
		analyticsHandler: analyticsHandler,
		sensorHandler:    sensorHandler,
		wsHandler:        wsHandler,
		jwtSecret:        jwtSecret,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	auth := engine.Group("/auth")
	{
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/logout", r.authHandler.Logout)
	}

	protected := engine.Group("/")
	protected.Use(middleware.AuthMiddleware(r.jwtSecret))
	{
		protected.GET("/auth/me", r.authHandler.Me)

		device := protected.Group("/device")
		{
			device.GET("/config", r.deviceHandler.GetConfig)
			device.PUT("/config", r.deviceHandler.UpdateConfig)
			device.GET("/status", r.deviceHandler.GetStatus)
		}

		sessions := protected.Group("/sessions")
		{
			sessions.GET("", r.sessionHandler.GetAll)
			sessions.GET("/summary", r.sessionHandler.GetSummary)
			sessions.GET("/:id", r.sessionHandler.GetByID)
		}

		analytics := protected.Group("/analytics")
		{
			analytics.GET("/peak-hours", r.analyticsHandler.GetPeakHours)
			analytics.GET("/condition-breakdown", r.analyticsHandler.GetConditionBreakdown)
			analytics.GET("/timeline", r.analyticsHandler.GetTimeline)
		}

		protected.GET("/sensor-logs", r.sensorHandler.GetLogs)
	}

	engine.GET("/ws", r.wsHandler.ServeWS)
}
