package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type Message struct {
	Type      string      `json:"type"`
	Source    string      `json:"source"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func NewRabbitMQService() (*RabbitMQService, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://meeseecs:rabbitmq_password_123@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	service := &RabbitMQService{
		conn:    conn,
		channel: ch,
	}

	// Declare exchanges and queues
	if err := service.setupExchangesAndQueues(); err != nil {
		return nil, fmt.Errorf("failed to setup exchanges and queues: %w", err)
	}

	return service, nil
}

func (s *RabbitMQService) setupExchangesAndQueues() error {
	// Declare exchanges
	exchanges := []string{
		"meeseecs.alerts",
		"meeseecs.vulnerabilities",
		"meeseecs.runtime",
		"meeseecs.notifications",
	}

	for _, exchange := range exchanges {
		err := s.channel.ExchangeDeclare(
			exchange, // name
			"topic",  // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", exchange, err)
		}
	}

	// Declare queues
	queues := map[string]string{
		"runtime.alerts":        "meeseecs.alerts",
		"vulnerability.scans":   "meeseecs.vulnerabilities",
		"runtime.events":        "meeseecs.runtime",
		"notification.email":    "meeseecs.notifications",
		"notification.slack":    "meeseecs.notifications",
		"notification.webhook":  "meeseecs.notifications",
	}

	for queueName, exchange := range queues {
		_, err := s.channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}

		// Bind queue to exchange
		err = s.channel.QueueBind(
			queueName,    // queue name
			queueName,    // routing key
			exchange,     // exchange
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, exchange, err)
		}
	}

	return nil
}

// Publish messages
func (s *RabbitMQService) PublishAlert(alert interface{}) error {
	message := Message{
		Type:   "alert",
		Source: "api-gateway",
		Data:   alert,
	}
	return s.publish("meeseecs.alerts", "runtime.alerts", message)
}

func (s *RabbitMQService) PublishVulnerability(vuln interface{}) error {
	message := Message{
		Type:   "vulnerability",
		Source: "api-gateway",
		Data:   vuln,
	}
	return s.publish("meeseecs.vulnerabilities", "vulnerability.scans", message)
}

func (s *RabbitMQService) PublishNotification(notificationType string, data interface{}) error {
	message := Message{
		Type:   "notification",
		Source: "api-gateway",
		Data:   data,
	}
	
	routingKey := fmt.Sprintf("notification.%s", notificationType)
	return s.publish("meeseecs.notifications", routingKey, message)
}

func (s *RabbitMQService) PublishRuntimeEvent(event interface{}) error {
	message := Message{
		Type:   "runtime_event",
		Source: "api-gateway",
		Data:   event,
	}
	return s.publish("meeseecs.runtime", "runtime.events", message)
}

func (s *RabbitMQService) publish(exchange, routingKey string, message Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = s.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to %s with routing key %s", exchange, routingKey)
	return nil
}

// Consume messages
func (s *RabbitMQService) ConsumeAlerts(handler func([]byte) error) error {
	return s.consume("runtime.alerts", handler)
}

func (s *RabbitMQService) ConsumeVulnerabilities(handler func([]byte) error) error {
	return s.consume("vulnerability.scans", handler)
}

func (s *RabbitMQService) ConsumeRuntimeEvents(handler func([]byte) error) error {
	return s.consume("runtime.events", handler)
}

func (s *RabbitMQService) consume(queueName string, handler func([]byte) error) error {
	msgs, err := s.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("Error handling message from %s: %v", queueName, err)
			}
		}
	}()

	return nil
}

func (s *RabbitMQService) Close() error {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}