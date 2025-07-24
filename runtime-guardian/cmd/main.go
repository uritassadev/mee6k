package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

type RuntimeGuardian struct {
	APIGatewayURL string
	RabbitMQConn  *amqp.Connection
	RabbitMQCh    *amqp.Channel
}

type SecurityEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Source      string                 `json:"source"`
	EventType   string                 `json:"event_type"`
	Severity    string                 `json:"severity"`
	ContainerID string                 `json:"container_id"`
	ProcessName string                 `json:"process_name"`
	Command     string                 `json:"command"`
	Details     map[string]interface{} `json:"details"`
	PolicyName  string                 `json:"policy_name"`
	Action      string                 `json:"action"`
}

type PolicyRule struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Severity    string                 `json:"severity"`
	Conditions  map[string]interface{} `json:"conditions"`
	Action      string                 `json:"action"`
}

func main() {
	// Load environment variables
	godotenv.Load()

	guardian := &RuntimeGuardian{
		APIGatewayURL: getEnv("API_GATEWAY_URL", "http://api-gateway:8080"),
	}

	// Initialize RabbitMQ connection
	if err := guardian.initRabbitMQ(); err != nil {
		log.Printf("Failed to initialize RabbitMQ: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":   "mee6k-runtime-guardian",
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// Runtime monitoring endpoints
	r.GET("/api/v1/events", guardian.getEvents)
	r.GET("/api/v1/policies", guardian.getPolicies)
	r.POST("/api/v1/policies", guardian.createPolicy)
	r.GET("/api/v1/status", guardian.getStatus)

	// Start background monitoring
	go guardian.startMonitoring()

	log.Printf("🛡️ Runtime Guardian starting on port 8081")
	log.Printf("🔗 API Gateway URL: %s", guardian.APIGatewayURL)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (rg *RuntimeGuardian) initRabbitMQ() error {
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://mee6k:rabbitmq_password_123@rabbitmq:5672/")
	
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return err
	}
	
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	
	rg.RabbitMQConn = conn
	rg.RabbitMQCh = ch
	
	// Declare exchange for security events
	return ch.ExchangeDeclare(
		"security_events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (rg *RuntimeGuardian) getEvents(c *gin.Context) {
	// Mock runtime events based on Falco patterns
	events := []SecurityEvent{
		{
			ID:          "evt_001",
			Timestamp:   time.Now().Add(-5 * time.Minute),
			Source:      "runtime_guardian",
			EventType:   "process_execution",
			Severity:    "HIGH",
			ContainerID: "container_123",
			ProcessName: "nc",
			Command:     "nc -l -p 4444",
			Details: map[string]interface{}{
				"pid":       1234,
				"ppid":      1,
				"user":      "root",
				"image":     "nginx:latest",
				"namespace": "default",
			},
			PolicyName: "Suspicious Process Execution",
			Action:     "alert",
		},
		{
			ID:          "evt_002",
			Timestamp:   time.Now().Add(-2 * time.Minute),
			Source:      "runtime_guardian",
			EventType:   "file_access",
			Severity:    "MEDIUM",
			ContainerID: "container_456",
			ProcessName: "cat",
			Command:     "cat /etc/passwd",
			Details: map[string]interface{}{
				"file_path": "/etc/passwd",
				"access_mode": "read",
				"user":      "www-data",
				"image":     "apache:2.4",
			},
			PolicyName: "Sensitive File Access",
			Action:     "alert",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"total":  len(events),
		"limit":  50,
		"offset": 0,
	})
}

func (rg *RuntimeGuardian) getPolicies(c *gin.Context) {
	// Mock security policies based on Falco rules
	policies := []PolicyRule{
		{
			Name:        "Suspicious Process Execution",
			Description: "Detect execution of suspicious processes in containers",
			Enabled:     true,
			Severity:    "HIGH",
			Conditions: map[string]interface{}{
				"processes": []string{"nc", "ncat", "netcat", "nmap", "wget", "curl"},
				"action":    "alert",
			},
			Action: "alert",
		},
		{
			Name:        "Privilege Escalation Detection",
			Description: "Monitor attempts to escalate privileges",
			Enabled:     true,
			Severity:    "CRITICAL",
			Conditions: map[string]interface{}{
				"commands": []string{"sudo", "su", "chmod +s"},
				"action":   "block",
			},
			Action: "block",
		},
		{
			Name:        "Crypto Mining Detection",
			Description: "Detect cryptocurrency mining activities",
			Enabled:     true,
			Severity:    "CRITICAL",
			Conditions: map[string]interface{}{
				"processes": []string{"xmrig", "cpuminer", "cgminer"},
				"network":   []string{"stratum"},
				"action":    "terminate",
			},
			Action: "terminate",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"policies": policies,
		"total":    len(policies),
	})
}

func (rg *RuntimeGuardian) createPolicy(c *gin.Context) {
	var policy PolicyRule
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, this would save to database
	log.Printf("Created new policy: %s", policy.Name)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Policy created successfully",
		"policy":  policy,
	})
}

func (rg *RuntimeGuardian) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":     "runtime_guardian",
		"status":      "active",
		"last_check":  time.Now().UTC(),
		"policies":    3,
		"events_24h":  156,
		"containers":  12,
		"version":     "1.0.0",
	})
}

func (rg *RuntimeGuardian) startMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Printf("🔍 Starting runtime monitoring...")

	for {
		select {
		case <-ticker.C:
			rg.performSecurityCheck()
		}
	}
}

func (rg *RuntimeGuardian) performSecurityCheck() {
	// Simulate security event detection
	event := SecurityEvent{
		ID:          generateEventID(),
		Timestamp:   time.Now().UTC(),
		Source:      "runtime_guardian",
		EventType:   "security_check",
		Severity:    "INFO",
		ContainerID: "monitoring",
		ProcessName: "guardian",
		Command:     "security_scan",
		Details: map[string]interface{}{
			"scan_type":    "periodic",
			"containers":   12,
			"processes":    45,
			"connections":  8,
		},
		PolicyName: "Periodic Security Check",
		Action:     "monitor",
	}

	// Send event to message queue
	if rg.RabbitMQCh != nil {
		rg.publishEvent(event)
	}

	log.Printf("🔍 Security check completed - %d containers monitored", 12)
}

func (rg *RuntimeGuardian) publishEvent(event SecurityEvent) {
	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	err = rg.RabbitMQCh.Publish(
		"security_events",
		"runtime.security.event",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
}

func generateEventID() string {
	return "evt_" + time.Now().Format("20060102150405")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}