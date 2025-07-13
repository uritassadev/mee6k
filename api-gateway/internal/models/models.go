package models

import (
	"time"
	"gorm.io/gorm"
)

// Alert represents a security alert in the system
type Alert struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Source      string    `json:"source" gorm:"not null"` // runtime-guardian, vuln-scanner
	Type        string    `json:"type" gorm:"not null"`   // runtime, vulnerability
	Severity    string    `json:"severity" gorm:"not null"` // CRITICAL, HIGH, MEDIUM, LOW
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"` // container, image, pod name
	Namespace   string    `json:"namespace"`
	Status      string    `json:"status" gorm:"default:'open'"` // open, acknowledged, resolved
	Metadata    string    `json:"metadata" gorm:"type:jsonb"` // Additional context as JSON
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CVE             string    `json:"cve" gorm:"uniqueIndex"`
	Severity        string    `json:"severity" gorm:"not null"`
	Score           float64   `json:"score"`
	Title           string    `json:"title" gorm:"not null"`
	Description     string    `json:"description"`
	Package         string    `json:"package"`
	InstalledVersion string   `json:"installed_version"`
	FixedVersion    string    `json:"fixed_version"`
	ImageName       string    `json:"image_name"`
	ImageTag        string    `json:"image_tag"`
	Status          string    `json:"status" gorm:"default:'open'"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SecurityPolicy represents a runtime security policy
type SecurityPolicy struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled" gorm:"default:true"`
	Severity    string    `json:"severity" gorm:"not null"`
	Rules       string    `json:"rules" gorm:"type:jsonb"` // Policy rules as JSON
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AlertChannel represents notification channels
type AlertChannel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Type        string    `json:"type" gorm:"not null"` // email, slack, pagerduty, webhook
	Config      string    `json:"config" gorm:"type:jsonb"` // Channel configuration as JSON
	Enabled     bool      `json:"enabled" gorm:"default:true"`
	Severities  string    `json:"severities" gorm:"type:jsonb"` // Array of severities to notify
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ScanResult represents vulnerability scan results
type ScanResult struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	ImageName        string    `json:"image_name" gorm:"not null"`
	ImageTag         string    `json:"image_tag" gorm:"not null"`
	ScanType         string    `json:"scan_type" gorm:"not null"` // vulnerability, secret, config
	Status           string    `json:"status" gorm:"not null"` // running, completed, failed
	TotalVulns       int       `json:"total_vulns"`
	CriticalVulns    int       `json:"critical_vulns"`
	HighVulns        int       `json:"high_vulns"`
	MediumVulns      int       `json:"medium_vulns"`
	LowVulns         int       `json:"low_vulns"`
	Results          string    `json:"results" gorm:"type:jsonb"` // Detailed results as JSON
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RuntimeEvent represents runtime security events
type RuntimeEvent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	EventType   string    `json:"event_type" gorm:"not null"`
	Severity    string    `json:"severity" gorm:"not null"`
	Source      string    `json:"source" gorm:"not null"` // container, process, network
	Message     string    `json:"message" gorm:"not null"`
	ContainerID string    `json:"container_id"`
	ImageName   string    `json:"image_name"`
	ProcessName string    `json:"process_name"`
	Command     string    `json:"command"`
	Metadata    string    `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time `json:"created_at"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalAlerts        int64 `json:"total_alerts"`
	CriticalAlerts     int64 `json:"critical_alerts"`
	HighAlerts         int64 `json:"high_alerts"`
	MediumAlerts       int64 `json:"medium_alerts"`
	LowAlerts          int64 `json:"low_alerts"`
	TotalVulns         int64 `json:"total_vulnerabilities"`
	CriticalVulns      int64 `json:"critical_vulnerabilities"`
	HighVulns          int64 `json:"high_vulnerabilities"`
	ActivePolicies     int64 `json:"active_policies"`
	ScannedImages      int64 `json:"scanned_images"`
	RuntimeEvents24h   int64 `json:"runtime_events_24h"`
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Alert{},
		&Vulnerability{},
		&SecurityPolicy{},
		&AlertChannel{},
		&ScanResult{},
		&RuntimeEvent{},
	)
}