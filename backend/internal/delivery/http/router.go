package http

import (
	"github.com/gin-gonic/gin"
	"github.com/seymourrisey/staredesk/internal/delivery/http/handler"
	"github.com/seymourrisey/staredesk/internal/delivery/http/middleware"
	websocketdelivery "github.com/seymourrisey/staredesk/internal/delivery/websocket"
)

type Router struct {
	authHandler   *handler.AuthHandler
	deviceHandler *handler.DeviceHandler
	wsHandler     *websocketdelivery.Handler
	jwtSecret     string
}

func NewRouter(
	authHandler *handler.AuthHandler,
	deviceHandler *handler.DeviceHandler,
	wsHandler *websocketdelivery.Handler,
	jwtSecret string,
) *Router {
	return &Router{
		authHandler:   authHandler,
		deviceHandler: deviceHandler,
		wsHandler:     wsHandler,
		jwtSecret:     jwtSecret,
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
	}

	// /ws tidak pakai AuthMiddleware — JWT divalidasi langsung di wsHandler
	// sebelum upgrade HTTP → WebSocket
	engine.GET("/ws", r.wsHandler.ServeWS)
}
