package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/seymourrisey/staredesk/config"
	deliveryhttp "github.com/seymourrisey/staredesk/internal/delivery/http"
	"github.com/seymourrisey/staredesk/internal/delivery/http/handler"
	mqttdelivery "github.com/seymourrisey/staredesk/internal/delivery/mqtt"
	websocketdelivery "github.com/seymourrisey/staredesk/internal/delivery/websocket"
	"github.com/seymourrisey/staredesk/internal/infrastructure/broker"
	"github.com/seymourrisey/staredesk/internal/infrastructure/postgres"
	"github.com/seymourrisey/staredesk/internal/usecase"
)

func main() {
	cfg := config.Load()

	// Database
	db, err := postgres.NewPool(&cfg.DB)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	sessionRepo := postgres.NewSessionPostgres(db)
	sensorLogRepo := postgres.NewSensorLogPostgres(db)
	deviceConfigRepo := postgres.NewDeviceConfigPostgres(db)
	deviceStatusRepo := postgres.NewDeviceStatusPostgres(db)
	analyticsRepo := postgres.NewAnalyticsPostgres(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret)
	deviceUsecase := usecase.NewDeviceUsecase(deviceConfigRepo, deviceStatusRepo)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, deviceConfigRepo)
	sensorUsecase := usecase.NewSensorUsecase(sensorLogRepo)
	analyticsUsecase := usecase.NewAnalyticsUsecase(analyticsRepo)

	// WebSocket Hub
	hub := websocketdelivery.NewHub()
	go hub.Run()
	wsHandler := websocketdelivery.NewHandler(hub, cfg.JWT.Secret)

	// MQTT
	mqttClient, err := broker.NewMQTTClient(
		cfg.MQTT.Broker,
		cfg.MQTT.ClientID,
		cfg.MQTT.Username,
		cfg.MQTT.Password,
		cfg.MQTT.Port,
	)
	if err != nil {
		log.Fatalf("MQTT connection failed: %v", err)
	}
	defer mqttClient.Disconnect()

	subscriber := mqttdelivery.NewSubscriber(mqttClient, cfg.MQTT.UserID, hub, sessionUsecase, sensorUsecase, deviceUsecase)
	if err := subscriber.SubscribeAll(); err != nil {
		log.Fatalf("MQTT subscribe failed: %v", err)
	}

	// HTTP Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	deviceHandler := handler.NewDeviceHandler(deviceUsecase, subscriber)
	sessionHandler := handler.NewSessionHandler(sessionUsecase)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsUsecase)
	sensorHandler := handler.NewSensorHandler(sensorUsecase)

	// HTTP Router
	engine := gin.Default()
	log.Printf("CORS allowed origins: %s", cfg.App.AllowedOrigins)
	router := deliveryhttp.NewRouter(authHandler, deviceHandler, sessionHandler, analyticsHandler, sensorHandler, wsHandler, cfg.JWT.Secret, cfg.App.AllowedOrigins)
	router.Setup(engine)

	log.Printf("StareDesk backend starting on port %s", cfg.App.Port)
	if err := engine.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
