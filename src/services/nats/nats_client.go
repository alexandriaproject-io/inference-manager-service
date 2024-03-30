// nats_client.go

package natsClient

import (
	"log"
	"time"
	"github.com/nats-io/nats.go"
	"inference-manager-service/src/config" // Adjust the import path based on your project structure
)

// Client struct holds the NATS and JetStream connection.
type Client struct {
	Conn *nats.Conn
	Js   nats.JetStreamContext
	Config *config.Config
}

// NewNatsClient creates a new NATS client with a connection to the server, including JetStream.
func NewNatsClient(cfg *config.Config) (*Client, error) {
	opts := []nats.Option{nats.Name("NATS Inference Manager Client")}

	if cfg.NatsUser != "" && cfg.NatsPass != "" {
		opts = append(opts, nats.UserInfo(cfg.NatsUser, cfg.NatsPass))
	}

	conn, err := nats.Connect(cfg.NatsServerURL, opts...)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}

	return &Client{Conn: conn, Js: js, Config: cfg}, nil
}

// StartNatsClient starts the NATS client and subscribes to a subject for standard NATS messages.
func (c *Client) StartNatsClient(subject string) error {
	_, err := c.Conn.Subscribe(subject, func(m *nats.Msg) {
		log.Printf("Received a standard message on %s: %s", m.Subject, string(m.Data))
		// Simulate task processing
		time.Sleep(time.Duration(c.Config.EXECUTING_TIME) * time.Second)
		// Respond to message
		if err := c.Publish(m.Reply, []byte("Processed: "+string(m.Data))); err != nil {
			log.Println("Failed to publish standard response (dropping message):", err)
		}
	})
	if err != nil {
		return err
	}

	log.Printf("Listening for standard messages on [%s]", subject)
	return nil
}

// StartJetStreamClient starts the NATS client for JetStream and subscribes to a subject.
func (c *Client) StartJetStreamClient(subject string) error {
	_, err := c.Js.Subscribe(subject, func(m *nats.Msg) {
		log.Printf("Received a JetStream message on %s: %s", m.Subject, string(m.Data))
		// Simulate task processing
		time.Sleep(time.Duration(c.Config.EXECUTING_TIME) * time.Second)
		// Respond to message
		if _, err := c.Js.Publish(m.Reply, []byte("Processed: "+string(m.Data))); err != nil {
			log.Println("Failed to publish JetStream response, retrying:", err)
			// Implement retry logic as needed
		}
	}, nats.Durable("myDurableConsumer"), nats.ManualAck())

	if err != nil {
		return err
	}

	log.Printf("Listening for JetStream messages on [%s]", subject)
	return nil
}

// Publish sends a message to a subject using standard NATS.
func (c *Client) Publish(subject string, message []byte) error {
	return c.Conn.Publish(subject, message)
}

// PublishJetStream sends a message to a subject using JetStream.
func (c *Client) PublishJetStream(subject string, message []byte) error {
	_, err := c.Js.Publish(subject, message)
	return err
}