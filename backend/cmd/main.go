package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/seymourrisey/staredesk/config"
	deliveryhttp "github.com/seymourrisey/staredesk/internal/delivery/http"
	"github.com/seymourrisey/staredesk/internal/delivery/http/handler"
	"github.com/seymourrisey/staredesk/internal/infrastructure/postgres"
	"github.com/seymourrisey/staredesk/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewPool(&cfg.DB)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	log.Printf("Connected to database")

	// Repositories
	userRepo := postgres.NewUserRepository(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret)

	// Handlers
	authHandler := handler.NewAuthHandler(authUsecase)

	// Router
	engine := gin.Default()
	router := deliveryhttp.NewRouter(authHandler, cfg.JWT.Secret)
	router.Setup(engine)

	log.Printf("StareDesk backend starting on port %s", cfg.App.Port)
	if err := engine.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
