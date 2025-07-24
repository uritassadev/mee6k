package handlers

import (
	"net/http"
	"strconv"
	"time"

	"mee6k-box/api-gateway/internal/models"
	"mee6k-box/api-gateway/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db     *services.DatabaseService
	redis  *services.RedisService
	rabbit *services.RabbitMQService
}

func NewHandler(db *services.DatabaseService, redis *services.RedisService, rabbit *services.RabbitMQService) *Handler {
	return &Handler{
		db:     db,
		redis:  redis,
		rabbit: rabbit,
	}
}

// Health check
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "mee6k-box-api-gateway",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
	})
}

// Dashboard handlers
func (h *Handler) GetDashboardStats(c *gin.Context) {
	// Try to get from cache first
	var stats models.DashboardStats
	if err := h.redis.Get("dashboard:stats", &stats); err == nil {
		c.JSON(http.StatusOK, stats)
		return
	}

	// Get from database
	dbStats, err := h.db.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}

	// Cache for 5 minutes
	h.redis.Set("dashboard:stats", dbStats, 5*time.Minute)

	c.JSON(http.StatusOK, dbStats)
}

func (h *Handler) GetSecurityOverview(c *gin.Context) {
	overview := gin.H{
		"platform": gin.H{
			"name":    "MEE6K Box",
			"version": "1.0.0",
			"status":  "operational",
		},
		"components": gin.H{
			"runtime_guardian": gin.H{
				"status": "active",
				"last_check": time.Now().Add(-2 * time.Minute),
			},
			"vuln_scanner": gin.H{
				"status": "active", 
				"last_scan": time.Now().Add(-30 * time.Minute),
			},
			"alert_engine": gin.H{
				"status": "active",
				"processed_alerts": 156,
			},
		},
		"security_posture": gin.H{
			"risk_level": "MEDIUM",
			"compliance_score": 85,
			"last_assessment": time.Now().Add(-1 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, overview)
}

// Runtime security handlers
func (h *Handler) GetRuntimeAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	severity := c.Query("severity")

	alerts, err := h.db.GetAlerts(limit, offset, severity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get runtime alerts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"total":  len(alerts),
		"limit":  limit,
		"offset": offset,
	})
}

func (h *Handler) GetSecurityPolicies(c *gin.Context) {
	policies, err := h.db.GetSecurityPolicies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get security policies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"policies": policies,
		"total":    len(policies),
	})
}

func (h *Handler) CreateSecurityPolicy(c *gin.Context) {
	var policy models.SecurityPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.CreateSecurityPolicy(&policy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create security policy"})
		return
	}

	// Publish policy update to runtime guardian
	h.rabbit.PublishRuntimeEvent(gin.H{
		"type": "policy_created",
		"policy": policy,
	})

	c.JSON(http.StatusCreated, policy)
}

func (h *Handler) UpdateSecurityPolicy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID"})
		return
	}

	policy, err := h.db.GetSecurityPolicyByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
		return
	}

	if err := c.ShouldBindJSON(policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateSecurityPolicy(policy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update security policy"})
		return
	}

	// Publish policy update
	h.rabbit.PublishRuntimeEvent(gin.H{
		"type": "policy_updated",
		"policy": policy,
	})

	c.JSON(http.StatusOK, policy)
}

func (h *Handler) DeleteSecurityPolicy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID"})
		return
	}

	if err := h.db.DeleteSecurityPolicy(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete security policy"})
		return
	}

	// Publish policy deletion
	h.rabbit.PublishRuntimeEvent(gin.H{
		"type": "policy_deleted",
		"policy_id": id,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Policy deleted successfully"})
}

// Vulnerability handlers
func (h *Handler) GetVulnerabilities(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	severity := c.Query("severity")

	vulns, err := h.db.GetVulnerabilities(limit, offset, severity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vulnerabilities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vulnerabilities": vulns,
		"total":          len(vulns),
		"limit":          limit,
		"offset":         offset,
	})
}

func (h *Handler) GetVulnerabilitySummary(c *gin.Context) {
	stats, err := h.db.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vulnerability summary"})
		return
	}

	summary := gin.H{
		"total_vulnerabilities": stats.TotalVulns,
		"by_severity": gin.H{
			"critical": stats.CriticalVulns,
			"high":     stats.HighVulns,
		},
		"scanned_images": stats.ScannedImages,
		"last_scan":     time.Now().Add(-30 * time.Minute),
	}

	c.JSON(http.StatusOK, summary)
}

func (h *Handler) GetVulnerabilityReports(c *gin.Context) {
	// This would typically fetch detailed scan reports
	reports := []gin.H{
		{
			"id":         1,
			"image":      "nginx:latest",
			"scan_date":  time.Now().Add(-2 * time.Hour),
			"status":     "completed",
			"critical":   2,
			"high":       5,
			"medium":     12,
			"low":        8,
		},
		{
			"id":         2,
			"image":      "redis:alpine",
			"scan_date":  time.Now().Add(-4 * time.Hour),
			"status":     "completed",
			"critical":   0,
			"high":       1,
			"medium":     3,
			"low":        2,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"total":   len(reports),
	})
}

func (h *Handler) TriggerVulnerabilityScan(c *gin.Context) {
	var request struct {
		ImageName string `json:"image_name" binding:"required"`
		ImageTag  string `json:"image_tag"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.ImageTag == "" {
		request.ImageTag = "latest"
	}

	// Publish scan request
	scanRequest := gin.H{
		"type":       "vulnerability_scan",
		"image_name": request.ImageName,
		"image_tag":  request.ImageTag,
		"timestamp":  time.Now(),
	}

	if err := h.rabbit.PublishVulnerability(scanRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to trigger scan"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Vulnerability scan triggered",
		"image":   request.ImageName + ":" + request.ImageTag,
		"status":  "queued",
	})
}

// Alert handlers
func (h *Handler) GetAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	severity := c.Query("severity")

	alerts, err := h.db.GetAlerts(limit, offset, severity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alerts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"total":  len(alerts),
	})
}

func (h *Handler) CreateAlert(c *gin.Context) {
	var alert models.Alert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.CreateAlert(&alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert"})
		return
	}

	// Publish alert for notifications
	h.rabbit.PublishAlert(alert)

	c.JSON(http.StatusCreated, alert)
}

func (h *Handler) AcknowledgeAlert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	alert, err := h.db.GetAlertByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	alert.Status = "acknowledged"
	if err := h.db.UpdateAlert(alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acknowledge alert"})
		return
	}

	c.JSON(http.StatusOK, alert)
}

func (h *Handler) GetAlertChannels(c *gin.Context) {
	channels, err := h.db.GetAlertChannels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alert channels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"channels": channels,
		"total":    len(channels),
	})
}

func (h *Handler) CreateAlertChannel(c *gin.Context) {
	var channel models.AlertChannel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.CreateAlertChannel(&channel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert channel"})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

// Settings handlers
func (h *Handler) GetNotificationSettings(c *gin.Context) {
	settings := gin.H{
		"email": gin.H{
			"enabled": true,
			"smtp_host": "smtp.example.com",
			"smtp_port": 587,
		},
		"slack": gin.H{
			"enabled": false,
			"webhook_url": "",
		},
		"pagerduty": gin.H{
			"enabled": false,
			"integration_key": "",
		},
	}

	c.JSON(http.StatusOK, settings)
}

func (h *Handler) UpdateNotificationSettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cache the settings
	h.redis.Set("notification:settings", settings, 0)

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification settings updated",
		"settings": settings,
	})
}