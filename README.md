# MEE6K Box 🛡️

**The Ultimate Container Security Platform**

MEE6K Box is a comprehensive security platform that provides real-time runtime protection and vulnerability management for containerized environments.

## 🎨 Platform Design
- **Primary Colors**: Blue (#2563eb), Purple (#7c3aed), Orange (#ea580c)
- **Modern UI**: Clean, intuitive interface with dark/light theme support
- **Responsive**: Works seamlessly across desktop and mobile devices

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Runtime         │    │ Vulnerability   │    │ Alert Engine    │
│ Guardian        │    │ Scanner         │    │                 │
│ (Runtime Sec)   │    │ (Vuln Mgmt)     │    │ (Notifications) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │ MEE6K Box    │
                    │ Dashboard       │
                    └─────────────────┘
```

## 🚀 Components

### Runtime Guardian
- Real-time container behavior monitoring
- Custom security policies
- Threat detection and response
- Behavioral anomaly detection

### Vulnerability Scanner  
- Comprehensive image scanning
- CVE database integration
- Compliance reporting
- Risk assessment and prioritization

### Alert Engine
- Multi-channel notifications (Email, Slack, PagerDuty)
- Severity-based routing
- Incident management integration
- Custom alert rules

### Dashboard
- Unified security overview
- Interactive vulnerability reports
- Real-time threat monitoring
- Executive security summaries

## 🛠️ Technology Stack

- **Backend**: Go (API Gateway, Microservices)
- **Frontend**: React with TypeScript
- **Database**: PostgreSQL + Redis
- **Message Queue**: RabbitMQ
- **Deployment**: Docker Compose
- **Monitoring**: Prometheus + Grafana

## 📋 Project Phases

1. **Phase 1**: Core Infrastructure & Runtime Guardian
2. **Phase 2**: Vulnerability Scanner Integration  
3. **Phase 3**: Alert Engine & Notifications
4. **Phase 4**: Dashboard & UI
5. **Phase 5**: Integration & Testing

## 🚀 Quick Start

```bash
# Clone and setup
git clone <repo-url>
cd mee6k-box/docker

# Start the platform
docker-compose up -d

# Access dashboard
open http://localhost:3000
```

## 📊 Features

- ✅ Real-time security monitoring
- ✅ Vulnerability assessment
- ✅ Policy-based alerting
- ✅ Compliance reporting
- ✅ Incident management
- ✅ Executive dashboards