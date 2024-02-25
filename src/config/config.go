// config.go

package config

import (
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	NatsServerURL string
	NatsUser      string
	NatsPass      string
	NATS_TOPIC string
	NATS_JS_TOPIC string
	EXECUTING_TIME int
}

// LoadConfig reads configuration from environment variables or .env file
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	return &Config{
		NatsServerURL: getEnv("NATS_SERVER_URL", "nats://localhost:4222"),
		NatsUser:      getEnv("NATS_USER", ""),
		NatsPass:      getEnv("NATS_PASS", ""),
		NATS_TOPIC: getEnv("NATS_TOPIC", ""),
		NATS_JS_TOPIC: getEnv("NATS_JS_TOPIC", ""),
		EXECUTING_TIME: getEnvAsInt("EXECUTING_TIME", 1),
	}
}

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	} else {
		log.Printf("Warning: Unable to convert environment variable %s=%s to int, using default %d\n", key, valueStr, defaultValue)
		return defaultValue
	}
}
