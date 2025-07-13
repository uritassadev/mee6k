-- Meeseecs Box Database Initialization
-- This script sets up the initial database structure and sample data

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Insert sample security policies
INSERT INTO security_policies (name, description, enabled, severity, rules, created_at, updated_at) VALUES
('Suspicious Process Execution', 'Detect execution of suspicious processes in containers', true, 'HIGH', 
 '{"processes": ["nc", "ncat", "netcat", "nmap", "wget", "curl"], "action": "alert"}', NOW(), NOW()),

('Privilege Escalation Detection', 'Monitor attempts to escalate privileges', true, 'CRITICAL',
 '{"commands": ["sudo", "su", "chmod +s"], "action": "block"}', NOW(), NOW()),

('Crypto Mining Detection', 'Detect cryptocurrency mining activities', true, 'CRITICAL',
 '{"processes": ["xmrig", "cpuminer", "cgminer"], "network": ["stratum"], "action": "terminate"}', NOW(), NOW()),

('Sensitive File Access', 'Monitor access to sensitive system files', true, 'MEDIUM',
 '{"files": ["/etc/passwd", "/etc/shadow", "/root/.ssh/", "/home/*/.ssh/"], "action": "alert"}', NOW(), NOW()),

('Network Anomaly Detection', 'Detect unusual network connections', true, 'HIGH',
 '{"ports": [22, 23, 3389], "protocols": ["ssh", "telnet", "rdp"], "action": "alert"}', NOW(), NOW());

-- Insert sample alert channels
INSERT INTO alert_channels (name, type, config, enabled, severities, created_at, updated_at) VALUES
('Email Alerts', 'email', 
 '{"smtp_host": "smtp.example.com", "smtp_port": 587, "username": "alerts@meeseecs.com", "recipients": ["security@company.com"]}',
 true, '["CRITICAL", "HIGH"]', NOW(), NOW()),

('Slack Security Channel', 'slack',
 '{"webhook_url": "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK", "channel": "#security-alerts"}',
 false, '["CRITICAL", "HIGH", "MEDIUM"]', NOW(), NOW()),

('PagerDuty Integration', 'pagerduty',
 '{"integration_key": "YOUR_PAGERDUTY_INTEGRATION_KEY", "service_id": "PXXXXXX"}',
 false, '["CRITICAL"]', NOW(), NOW());

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_alerts_severity ON alerts(severity);
CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status);
CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON alerts(created_at);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_severity ON vulnerabilities(severity);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_cve ON vulnerabilities(cve);
CREATE INDEX IF NOT EXISTS idx_runtime_events_created_at ON runtime_events(created_at);
CREATE INDEX IF NOT EXISTS idx_runtime_events_severity ON runtime_events(severity);