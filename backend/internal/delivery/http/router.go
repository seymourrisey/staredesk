package http

import (
	"github.com/gin-gonic/gin"
	"github.com/seymourrisey/staredesk/internal/delivery/http/handler"
	"github.com/seymourrisey/staredesk/internal/delivery/http/middleware"
)

type Router struct {
	authHandler *handler.AuthHandler
	jwtSecret   string
}

func NewRouter(authHandler *handler.AuthHandler, jwtSecret string) *Router {
	return &Router{authHandler: authHandler, jwtSecret: jwtSecret}
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
	}
}
