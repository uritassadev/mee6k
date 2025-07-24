package main

import (
	"log"
	"os"

	"mee6k-box/api-gateway/internal/handlers"
	"mee6k-box/api-gateway/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize services
	dbService, err := services.NewDatabaseService()
	if err != nil {
		log.Fatal("Failed to initialize database service:", err)
	}

	redisService, err := services.NewRedisService()
	if err != nil {
		log.Fatal("Failed to initialize Redis service:", err)
	}

	rabbitService, err := services.NewRabbitMQService()
	if err != nil {
		log.Fatal("Failed to initialize RabbitMQ service:", err)
	}

	// Initialize handlers
	handler := handlers.NewHandler(dbService, redisService, rabbitService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check
	r.GET("/health", handler.HealthCheck)

	// API routes
	api := r.Group("/api/v1")
	{
		// Dashboard routes
		api.GET("/dashboard/stats", handler.GetDashboardStats)
		api.GET("/dashboard/overview", handler.GetSecurityOverview)

		// Runtime security routes
		runtime := api.Group("/runtime")
		{
			runtime.GET("/alerts", handler.GetRuntimeAlerts)
			runtime.GET("/policies", handler.GetSecurityPolicies)
			runtime.POST("/policies", handler.CreateSecurityPolicy)
			runtime.PUT("/policies/:id", handler.UpdateSecurityPolicy)
			runtime.DELETE("/policies/:id", handler.DeleteSecurityPolicy)
		}

		// Vulnerability routes
		vulns := api.Group("/vulnerabilities")
		{
			vulns.GET("/", handler.GetVulnerabilities)
			vulns.GET("/summary", handler.GetVulnerabilitySummary)
			vulns.GET("/reports", handler.GetVulnerabilityReports)
			vulns.POST("/scan", handler.TriggerVulnerabilityScan)
		}

		// Alert routes
		alerts := api.Group("/alerts")
		{
			alerts.GET("/", handler.GetAlerts)
			alerts.POST("/", handler.CreateAlert)
			alerts.PUT("/:id/acknowledge", handler.AcknowledgeAlert)
			alerts.GET("/channels", handler.GetAlertChannels)
			alerts.POST("/channels", handler.CreateAlertChannel)
		}

		// Settings routes
		settings := api.Group("/settings")
		{
			settings.GET("/notifications", handler.GetNotificationSettings)
			settings.PUT("/notifications", handler.UpdateNotificationSettings)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 MEE6K Box API Gateway starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}