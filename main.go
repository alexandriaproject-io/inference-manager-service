// main.go

package main

import (
	"log"
	"inference-manager-service/src/config"
	"inference-manager-service/src/services/nats" 
)

func main() {
	cfg := config.LoadConfig()

	client, err := natsClient.NewNatsClient(cfg)
	if err != nil {
		log.Fatal("Error creating NATS client:", err)
	}
	defer client.Conn.Close()

	log.Printf("Starting NATS client for subject: %s", cfg.NATS_TOPIC)
	err = client.StartNatsClient(cfg.NATS_TOPIC)
	if err != nil {
		log.Fatal("Error starting NATS client for standard messages:", err)
	}

	log.Printf("Starting JetStream client for subject: %s", cfg.NATS_JS_TOPIC)
	err = client.StartJetStreamClient(cfg.NATS_JS_TOPIC)
	if err != nil {
		log.Fatal("Error starting NATS client for JetStream messages:", err)
	}

	// Keep the main goroutine alive
	select {}
}