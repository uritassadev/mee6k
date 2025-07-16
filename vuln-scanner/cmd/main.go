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

type VulnScanner struct {
	APIGatewayURL string
	RabbitMQConn  *amqp.Connection
	RabbitMQCh    *amqp.Channel
}

type ScanResult struct {
	ID           string        `json:"id"`
	ImageName    string        `json:"image_name"`
	ImageTag     string        `json:"image_tag"`
	ScanDate     time.Time     `json:"scan_date"`
	Status       string        `json:"status"`
	Critical     int           `json:"critical"`
	High         int           `json:"high"`
	Medium       int           `json:"medium"`
	Low          int           `json:"low"`
	Total        int           `json:"total"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type Vulnerability struct {
	ID          string    `json:"id"`
	CVE         string    `json:"cve"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Score       float64   `json:"score"`
	Package     string    `json:"package"`
	Version     string    `json:"version"`
	FixedIn     string    `json:"fixed_in"`
	PublishedAt time.Time `json:"published_at"`
}

type ScanRequest struct {
	ImageName string `json:"image_name" binding:"required"`
	ImageTag  string `json:"image_tag"`
}

func main() {
	// Load environment variables
	godotenv.Load()

	scanner := &VulnScanner{
		APIGatewayURL: getEnv("API_GATEWAY_URL", "http://api-gateway:8080"),
	}

	// Initialize RabbitMQ connection
	if err := scanner.initRabbitMQ(); err != nil {
		log.Printf("Failed to initialize RabbitMQ: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":   "meeseecs-vuln-scanner",
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// Vulnerability scanning endpoints
	r.GET("/api/v1/scans", scanner.getScans)
	r.POST("/api/v1/scan", scanner.startScan)
	r.GET("/api/v1/scan/:id", scanner.getScanResult)
	r.GET("/api/v1/status", scanner.getStatus)
	r.GET("/api/v1/vulnerabilities", scanner.getVulnerabilities)

	// Start background scanning
	go scanner.startPeriodicScanning()

	log.Printf("🔍 Vulnerability Scanner starting on port 8082")
	log.Printf("🔗 API Gateway URL: %s", scanner.APIGatewayURL)

	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (vs *VulnScanner) initRabbitMQ() error {
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://meeseecs:rabbitmq_password_123@rabbitmq:5672/")
	
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return err
	}
	
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	
	vs.RabbitMQConn = conn
	vs.RabbitMQCh = ch
	
	// Declare exchange for scan results
	return ch.ExchangeDeclare(
		"scan_results",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (vs *VulnScanner) getScans(c *gin.Context) {
	// Mock scan results based on Trivy patterns
	scans := []ScanResult{
		{
			ID:        "scan_001",
			ImageName: "nginx",
			ImageTag:  "latest",
			ScanDate:  time.Now().Add(-2 * time.Hour),
			Status:    "completed",
			Critical:  2,
			High:      5,
			Medium:    12,
			Low:       8,
			Total:     27,
		},
		{
			ID:        "scan_002",
			ImageName: "redis",
			ImageTag:  "alpine",
			ScanDate:  time.Now().Add(-4 * time.Hour),
			Status:    "completed",
			Critical:  0,
			High:      1,
			Medium:    3,
			Low:       2,
			Total:     6,
		},
		{
			ID:        "scan_003",
			ImageName: "postgres",
			ImageTag:  "15-alpine",
			ScanDate:  time.Now().Add(-1 * time.Hour),
			Status:    "running",
			Critical:  0,
			High:      0,
			Medium:    0,
			Low:       0,
			Total:     0,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"scans": scans,
		"total": len(scans),
		"limit": 50,
		"offset": 0,
	})
}

func (vs *VulnScanner) startScan(c *gin.Context) {
	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate scan ID
	scanID := "scan_" + time.Now().Format("20060102150405")
	
	// Start scan in background
	go vs.performScan(scanID, req.ImageName, req.ImageTag)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Scan started successfully",
		"scan_id": scanID,
		"image":   req.ImageName + ":" + req.ImageTag,
		"status":  "running",
	})
}

func (vs *VulnScanner) getScanResult(c *gin.Context) {
	scanID := c.Param("id")
	
	// Mock detailed scan result
	result := ScanResult{
		ID:        scanID,
		ImageName: "nginx",
		ImageTag:  "latest",
		ScanDate:  time.Now().Add(-30 * time.Minute),
		Status:    "completed",
		Critical:  2,
		High:      5,
		Medium:    12,
		Low:       8,
		Total:     27,
		Vulnerabilities: []Vulnerability{
			{
				ID:          "vuln_001",
				CVE:         "CVE-2023-1234",
				Title:       "Buffer Overflow in libssl",
				Description: "A buffer overflow vulnerability in OpenSSL library",
				Severity:    "CRITICAL",
				Score:       9.8,
				Package:     "libssl1.1",
				Version:     "1.1.1f-1ubuntu2.16",
				FixedIn:     "1.1.1f-1ubuntu2.17",
				PublishedAt: time.Now().Add(-30 * 24 * time.Hour),
			},
			{
				ID:          "vuln_002",
				CVE:         "CVE-2023-5678",
				Title:       "Information Disclosure in nginx",
				Description: "Information disclosure vulnerability in nginx server",
				Severity:    "HIGH",
				Score:       7.5,
				Package:     "nginx",
				Version:     "1.18.0-6ubuntu14.3",
				FixedIn:     "1.18.0-6ubuntu14.4",
				PublishedAt: time.Now().Add(-15 * 24 * time.Hour),
			},
		},
	}

	c.JSON(http.StatusOK, result)
}

func (vs *VulnScanner) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":      "vuln_scanner",
		"status":       "active",
		"last_scan":    time.Now().Add(-30 * time.Minute).UTC(),
		"total_scans":  156,
		"active_scans": 2,
		"images":       45,
		"version":      "1.0.0",
	})
}

func (vs *VulnScanner) getVulnerabilities(c *gin.Context) {
	// Mock vulnerability database
	vulns := []Vulnerability{
		{
			ID:          "vuln_001",
			CVE:         "CVE-2023-1234",
			Title:       "Buffer Overflow in libssl",
			Description: "A buffer overflow vulnerability in OpenSSL library",
			Severity:    "CRITICAL",
			Score:       9.8,
			Package:     "libssl1.1",
			Version:     "1.1.1f-1ubuntu2.16",
			FixedIn:     "1.1.1f-1ubuntu2.17",
			PublishedAt: time.Now().Add(-30 * 24 * time.Hour),
		},
		{
			ID:          "vuln_002",
			CVE:         "CVE-2023-5678",
			Title:       "Information Disclosure in nginx",
			Description: "Information disclosure vulnerability in nginx server",
			Severity:    "HIGH",
			Score:       7.5,
			Package:     "nginx",
			Version:     "1.18.0-6ubuntu14.3",
			FixedIn:     "1.18.0-6ubuntu14.4",
			PublishedAt: time.Now().Add(-15 * 24 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"vulnerabilities": vulns,
		"total":          len(vulns),
		"critical":       1,
		"high":           1,
		"medium":         0,
		"low":            0,
	})
}

func (vs *VulnScanner) startPeriodicScanning() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	log.Printf("🔍 Starting periodic vulnerability scanning...")

	for {
		select {
		case <-ticker.C:
			vs.performPeriodicScan()
		}
	}
}

func (vs *VulnScanner) performScan(scanID, imageName, imageTag string) {
	log.Printf("🔍 Starting scan %s for %s:%s", scanID, imageName, imageTag)
	
	// Simulate scan duration
	time.Sleep(30 * time.Second)
	
	// Create scan result
	result := ScanResult{
		ID:        scanID,
		ImageName: imageName,
		ImageTag:  imageTag,
		ScanDate:  time.Now().UTC(),
		Status:    "completed",
		Critical:  2,
		High:      3,
		Medium:    8,
		Low:       5,
		Total:     18,
	}

	// Publish scan result
	if vs.RabbitMQCh != nil {
		vs.publishScanResult(result)
	}

	log.Printf("✅ Scan %s completed - %d vulnerabilities found", scanID, result.Total)
}

func (vs *VulnScanner) performPeriodicScan() {
	images := []string{"nginx:latest", "redis:alpine", "postgres:15-alpine"}
	
	for _, image := range images {
		scanID := "periodic_" + time.Now().Format("20060102150405")
		parts := splitString(image, ":")
		imageName := parts[0]
		imageTag := "latest"
		if len(parts) > 1 {
			imageTag = parts[1]
		}
		
		go vs.performScan(scanID, imageName, imageTag)
		time.Sleep(5 * time.Second) // Stagger scans
	}
}

func splitString(s, sep string) []string {
	result := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if i < len(s)-len(sep)+1 && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func (vs *VulnScanner) publishScanResult(result ScanResult) {
	body, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to marshal scan result: %v", err)
		return
	}

	err = vs.RabbitMQCh.Publish(
		"scan_results",
		"vuln.scan.completed",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		log.Printf("Failed to publish scan result: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}