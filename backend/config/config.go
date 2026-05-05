package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	JWT  JWTConfig
	MQTT MQTTConfig
	App  AppConfig
}

type DBConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret string
}

type MQTTConfig struct {
	Broker   string
	Port     int
	Username string
	Password string
	ClientID string
	UserID   string
}

type AppConfig struct {
	Port           string
	AllowedOrigins []string
	IsProd         bool
}

func Load() *Config {
	for _, path := range []string{".env", "../.env"} {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	mqttPort, err := strconv.Atoi(getEnv("MQTT_PORT", "8883"))
	if err != nil {
		log.Fatalf("Invalid MQTT_PORT: %v", err)
	}

	return &Config{
		DB: DBConfig{
			DSN: mustGetEnv("DATABASE_URL"),
		},
		JWT: JWTConfig{
			Secret: mustGetEnv("JWT_SECRET"),
		},
		MQTT: MQTTConfig{
			Broker:   mustGetEnv("MQTT_BROKER"),
			Port:     mqttPort,
			Username: mustGetEnv("MQTT_USERNAME"),
			Password: mustGetEnv("MQTT_PASSWORD"),
			ClientID: getEnv("MQTT_CLIENT_ID", "staredesk-backend"),
			UserID:   mustGetEnv("MQTT_USER_ID"),
		},
		App: AppConfig{
			Port:           getEnv("APP_PORT", "8080"),
			AllowedOrigins: parseOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
			IsProd:         getEnv("APP_ENV", "development") == "production",
		},
	}
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return val
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func parseOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	var origins []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}
