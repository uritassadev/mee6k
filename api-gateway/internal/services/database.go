package services

import (
	"fmt"
	"os"

	"meeseecs-box/api-gateway/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseService struct {
	DB *gorm.DB
}

func NewDatabaseService() (*DatabaseService, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://meeseecs:secure_password_123@localhost:5432/meeseecs_box?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := models.AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DatabaseService{DB: db}, nil
}

// Alert operations
func (s *DatabaseService) CreateAlert(alert *models.Alert) error {
	return s.DB.Create(alert).Error
}

func (s *DatabaseService) GetAlerts(limit, offset int, severity string) ([]models.Alert, error) {
	var alerts []models.Alert
	query := s.DB.Order("created_at DESC")
	
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	
	err := query.Limit(limit).Offset(offset).Find(&alerts).Error
	return alerts, err
}

func (s *DatabaseService) GetAlertByID(id uint) (*models.Alert, error) {
	var alert models.Alert
	err := s.DB.First(&alert, id).Error
	return &alert, err
}

func (s *DatabaseService) UpdateAlert(alert *models.Alert) error {
	return s.DB.Save(alert).Error
}

// Vulnerability operations
func (s *DatabaseService) CreateVulnerability(vuln *models.Vulnerability) error {
	return s.DB.Create(vuln).Error
}

func (s *DatabaseService) GetVulnerabilities(limit, offset int, severity string) ([]models.Vulnerability, error) {
	var vulns []models.Vulnerability
	query := s.DB.Order("score DESC")
	
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	
	err := query.Limit(limit).Offset(offset).Find(&vulns).Error
	return vulns, err
}

// Security Policy operations
func (s *DatabaseService) CreateSecurityPolicy(policy *models.SecurityPolicy) error {
	return s.DB.Create(policy).Error
}

func (s *DatabaseService) GetSecurityPolicies() ([]models.SecurityPolicy, error) {
	var policies []models.SecurityPolicy
	err := s.DB.Find(&policies).Error
	return policies, err
}

func (s *DatabaseService) GetSecurityPolicyByID(id uint) (*models.SecurityPolicy, error) {
	var policy models.SecurityPolicy
	err := s.DB.First(&policy, id).Error
	return &policy, err
}

func (s *DatabaseService) UpdateSecurityPolicy(policy *models.SecurityPolicy) error {
	return s.DB.Save(policy).Error
}

func (s *DatabaseService) DeleteSecurityPolicy(id uint) error {
	return s.DB.Delete(&models.SecurityPolicy{}, id).Error
}

// Alert Channel operations
func (s *DatabaseService) CreateAlertChannel(channel *models.AlertChannel) error {
	return s.DB.Create(channel).Error
}

func (s *DatabaseService) GetAlertChannels() ([]models.AlertChannel, error) {
	var channels []models.AlertChannel
	err := s.DB.Find(&channels).Error
	return channels, err
}

// Dashboard statistics
func (s *DatabaseService) GetDashboardStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}
	
	// Count alerts by severity
	s.DB.Model(&models.Alert{}).Where("status != 'resolved'").Count(&stats.TotalAlerts)
	s.DB.Model(&models.Alert{}).Where("severity = 'CRITICAL' AND status != 'resolved'").Count(&stats.CriticalAlerts)
	s.DB.Model(&models.Alert{}).Where("severity = 'HIGH' AND status != 'resolved'").Count(&stats.HighAlerts)
	s.DB.Model(&models.Alert{}).Where("severity = 'MEDIUM' AND status != 'resolved'").Count(&stats.MediumAlerts)
	s.DB.Model(&models.Alert{}).Where("severity = 'LOW' AND status != 'resolved'").Count(&stats.LowAlerts)
	
	// Count vulnerabilities by severity
	s.DB.Model(&models.Vulnerability{}).Where("status != 'resolved'").Count(&stats.TotalVulns)
	s.DB.Model(&models.Vulnerability{}).Where("severity = 'CRITICAL' AND status != 'resolved'").Count(&stats.CriticalVulns)
	s.DB.Model(&models.Vulnerability{}).Where("severity = 'HIGH' AND status != 'resolved'").Count(&stats.HighVulns)
	
	// Count active policies
	s.DB.Model(&models.SecurityPolicy{}).Where("enabled = true").Count(&stats.ActivePolicies)
	
	// Count scanned images
	s.DB.Model(&models.ScanResult{}).Select("DISTINCT image_name").Count(&stats.ScannedImages)
	
	// Count runtime events in last 24h
	s.DB.Model(&models.RuntimeEvent{}).Where("created_at > NOW() - INTERVAL '24 hours'").Count(&stats.RuntimeEvents24h)
	
	return stats, nil
}