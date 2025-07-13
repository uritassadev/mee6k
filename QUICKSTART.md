# Meeseecs Box Quick Start Guide

## 🚀 Getting Started

1. **Start the platform:**
   ```bash
   ./start.sh
   ```

2. **Access the dashboard:**
   - Open http://localhost:3000 in your browser
   - The dashboard uses blue, purple, and orange theme colors

3. **API Documentation:**
   - API Gateway: http://localhost:8080
   - Health Check: http://localhost:8080/health

## 🔧 Management Commands

- **Start services:** `./start.sh`
- **Stop services:** `./stop.sh`
- **View logs:** `./logs.sh <service-name>`

## 📊 Service Ports

- Dashboard: 3000
- API Gateway: 8080
- PostgreSQL: 5432
- Redis: 6379
- RabbitMQ: 5672 (Management: 15672)

## 🛡️ Security Features

### Runtime Guardian
- Real-time container monitoring
- Custom security policies
- Behavioral anomaly detection

### Vulnerability Scanner
- Automated image scanning
- CVE database integration
- Compliance reporting

### Alert Engine
- Multi-channel notifications
- Severity-based routing
- Incident management

## 🎨 Platform Design

The platform uses a modern design with:
- **Primary Blue:** #2563eb
- **Purple Accent:** #7c3aed  
- **Orange Highlights:** #ea580c

## 📝 Next Steps

1. Configure notification channels in the dashboard
2. Set up custom security policies
3. Integrate with your container registry
4. Configure SMTP for email alerts
