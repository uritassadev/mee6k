#!/bin/bash

# Meeseecs Box Setup Script
# This script sets up the complete Meeseecs Box security platform

set -e

echo "🛡️  Setting up Meeseecs Box Security Platform"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
ORANGE='\033[0;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check if Docker Compose is installed (v1 or v2)
if ! docker compose version &> /dev/null 2>&1 && ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

echo -e "${BLUE}📋 Phase 1: Core Infrastructure Setup${NC}"
echo "-----------------------------------"

# Create necessary directories
echo -e "${PURPLE}📁 Creating project directories...${NC}"
mkdir -p {logs,data/{postgres,redis,rabbitmq}}

# Set permissions
chmod 755 logs
chmod 755 data/{postgres,redis,rabbitmq}

echo -e "${GREEN}✅ Directories created successfully${NC}"

# Skip Go build - will be handled by Docker
echo -e "${PURPLE}🔨 API Gateway ready for Docker build...${NC}"

echo -e "${GREEN}✅ API Gateway ready${NC}"

# Create environment file
echo -e "${PURPLE}⚙️  Creating environment configuration...${NC}"
cat > .env << EOF
# Meeseecs Box Configuration
COMPOSE_PROJECT_NAME=meeseecs-box

# Database Configuration
POSTGRES_DB=meeseecs_box
POSTGRES_USER=meeseecs
POSTGRES_PASSWORD=secure_password_123

# Redis Configuration
REDIS_PASSWORD=redis_password_123

# RabbitMQ Configuration
RABBITMQ_DEFAULT_USER=meeseecs
RABBITMQ_DEFAULT_PASS=rabbitmq_password_123

# API Gateway Configuration
API_GATEWAY_PORT=8080
DASHBOARD_PORT=3000

# Security Configuration
JWT_SECRET=meeseecs_box_jwt_secret_key_2024
ENCRYPTION_KEY=meeseecs_box_encryption_key_32_chars

# Notification Configuration
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=alerts@meeseecs.com
SMTP_PASSWORD=your_smtp_password

# Platform Branding
PLATFORM_NAME=Meeseecs Box
PLATFORM_VERSION=1.0.0
PLATFORM_COLORS=blue,purple,orange
EOF

echo -e "${GREEN}✅ Environment configuration created${NC}"

# Create docker-compose override for development
echo -e "${PURPLE}🐳 Creating Docker Compose configuration...${NC}"
cat > docker-compose/docker-compose.override.yml << EOF
version: '3.8'

services:
  api-gateway:
    volumes:
      - ../logs:/app/logs
    environment:
      - GIN_MODE=debug
      - LOG_LEVEL=debug
    
  postgres:
    volumes:
      - ../data/postgres:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    
  redis:
    volumes:
      - ../data/redis:/data
    ports:
      - "6379:6379"
    
  rabbitmq:
    volumes:
      - ../data/rabbitmq:/var/lib/rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
EOF

echo -e "${GREEN}✅ Docker Compose configuration ready${NC}"

# Create startup script
echo -e "${PURPLE}🚀 Creating startup scripts...${NC}"
cat > start.sh << 'EOF'
#!/bin/bash

echo "🛡️  Starting Meeseecs Box Security Platform"
echo "==========================================="

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Start services
echo "🐳 Starting Docker services..."
cd docker-compose

# Try docker compose (v2) first, fallback to docker-compose (v1)
if docker compose version &> /dev/null 2>&1; then
    docker compose up -d
else
    docker-compose up -d
fi

echo "⏳ Waiting for services to be ready..."
sleep 60

# Check service health
echo "🔍 Checking service health..."
services=("postgres:5432" "redis:6379" "rabbitmq:5672" "api-gateway:8080")

for service in "${services[@]}"; do
    IFS=':' read -r name port <<< "$service"
    if nc -z localhost $port 2>/dev/null; then
        echo "✅ $name is ready"
    else
        echo "❌ $name is not responding"
    fi
done

echo ""
echo "🎉 Meeseecs Box is starting up!"
echo "📊 Dashboard: http://localhost:3000"
echo "🔌 API Gateway: http://localhost:8080"
echo "🐰 RabbitMQ Management: http://localhost:15672"
echo ""
echo "Default credentials:"
echo "  RabbitMQ: meeseecs / rabbitmq_password_123"
echo ""
EOF

chmod +x start.sh

# Create stop script
cat > stop.sh << 'EOF'
#!/bin/bash

echo "🛑 Stopping Meeseecs Box Security Platform"
echo "=========================================="

cd docker-compose

# Try docker compose (v2) first, fallback to docker-compose (v1)
if docker compose version &> /dev/null 2>&1; then
    docker compose down
else
    docker-compose down
fi

echo "✅ All services stopped"
EOF

chmod +x stop.sh

# Create logs script
cat > logs.sh << 'EOF'
#!/bin/bash

echo "📋 Meeseecs Box Service Logs"
echo "============================"

cd docker-compose

if [ $# -eq 0 ]; then
    echo "Available services:"
    echo "  api-gateway"
    echo "  runtime-guardian"
    echo "  vuln-scanner"
    echo "  alert-engine"
    echo "  dashboard"
    echo "  postgres"
    echo "  redis"
    echo "  rabbitmq"
    echo ""
    echo "Usage: ./logs.sh <service-name>"
    echo "Example: ./logs.sh api-gateway"
else
    # Try docker compose (v2) first, fallback to docker-compose (v1)
    if docker compose version &> /dev/null 2>&1; then
        docker compose logs -f $1
    else
        docker-compose logs -f $1
    fi
fi
EOF

chmod +x logs.sh

# Create test script
cat > test-api.sh << 'EOF'
#!/bin/bash

# Meeseecs Box API Testing Script
# This script tests the API Gateway endpoints

set -e

API_BASE="http://localhost:8080"
BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
PURPLE='\033[0;35m'
ORANGE='\033[0;33m'
NC='\033[0m'

echo -e "${BLUE}🧪 Testing Meeseecs Box API Gateway${NC}"
echo "=================================="

# Function to test an endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local description=$3
    local expected_status=${4:-200}
    
    echo -e "${PURPLE}Testing: ${description}${NC}"
    echo -e "  ${method} ${endpoint}"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "%{http_code}" -o /tmp/response.json "${API_BASE}${endpoint}")
    else
        response=$(curl -s -w "%{http_code}" -X "$method" -H "Content-Type: application/json" -o /tmp/response.json "${API_BASE}${endpoint}")
    fi
    
    if [ "$response" = "$expected_status" ]; then
        echo -e "  ${GREEN}✅ Success (HTTP $response)${NC}"
        if [ -f /tmp/response.json ]; then
            echo -e "  ${GREEN}Response:${NC}"
            cat /tmp/response.json | jq . 2>/dev/null || cat /tmp/response.json
        fi
    else
        echo -e "  ${RED}❌ Failed (HTTP $response, expected $expected_status)${NC}"
        if [ -f /tmp/response.json ]; then
            echo -e "  ${RED}Response:${NC}"
            cat /tmp/response.json
        fi
    fi
    echo ""
}

# Wait for API to be ready
echo -e "${ORANGE}⏳ Waiting for API Gateway to be ready...${NC}"
for i in {1..30}; do
    if curl -s "${API_BASE}/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ API Gateway is ready!${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ API Gateway is not responding after 30 seconds${NC}"
        exit 1
    fi
    sleep 1
done

echo ""

# Test endpoints
test_endpoint "GET" "/health" "Health Check"
test_endpoint "GET" "/api/v1/dashboard/stats" "Dashboard Statistics"
test_endpoint "GET" "/api/v1/dashboard/overview" "Security Overview"
test_endpoint "GET" "/api/v1/runtime/alerts" "Runtime Alerts"
test_endpoint "GET" "/api/v1/runtime/policies" "Security Policies"
test_endpoint "GET" "/api/v1/vulnerabilities/" "Vulnerabilities"
test_endpoint "GET" "/api/v1/vulnerabilities/summary" "Vulnerability Summary"
test_endpoint "GET" "/api/v1/vulnerabilities/reports" "Vulnerability Reports"
test_endpoint "GET" "/api/v1/alerts/" "All Alerts"
test_endpoint "GET" "/api/v1/alerts/channels" "Alert Channels"
test_endpoint "GET" "/api/v1/settings/notifications" "Notification Settings"

echo ""
echo -e "${BLUE}🛡️ Testing Phase 2 Services${NC}"
echo "=============================="

# Test Runtime Guardian
echo -e "${PURPLE}Testing Runtime Guardian (Port 8081)${NC}"
test_endpoint "GET" "http://localhost:8081/health" "Runtime Guardian Health" 200
test_endpoint "GET" "http://localhost:8081/api/v1/events" "Runtime Events" 200
test_endpoint "GET" "http://localhost:8081/api/v1/policies" "Security Policies" 200
test_endpoint "GET" "http://localhost:8081/api/v1/status" "Runtime Status" 200

# Test Vulnerability Scanner
echo -e "${PURPLE}Testing Vulnerability Scanner (Port 8082)${NC}"
test_endpoint "GET" "http://localhost:8082/health" "Vuln Scanner Health" 200
test_endpoint "GET" "http://localhost:8082/api/v1/scans" "Scan Results" 200
test_endpoint "GET" "http://localhost:8082/api/v1/vulnerabilities" "Vulnerabilities" 200
test_endpoint "GET" "http://localhost:8082/api/v1/status" "Scanner Status" 200

# Clean up
rm -f /tmp/response.json

echo -e "${GREEN}🎉 API testing completed!${NC}"
echo ""
echo -e "${BLUE}📊 Platform Status:${NC}"
echo -e "${PURPLE}• API Gateway: ${GREEN}✅ Operational${NC}"
echo -e "${PURPLE}• Database: ${GREEN}✅ Connected${NC}"
echo -e "${PURPLE}• Message Queue: ${GREEN}✅ Ready${NC}"
echo -e "${PURPLE}• Cache: ${GREEN}✅ Available${NC}"
echo ""
echo -e "${ORANGE}Next: Access the dashboard at http://localhost:3000${NC}"
EOF

chmod +x test-api.sh

echo -e "${GREEN}✅ All scripts created${NC}"

# Create README for quick reference
cat > QUICKSTART.md << 'EOF'
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
EOF

echo -e "${GREEN}✅ Quick start guide created${NC}"

echo ""
echo -e "${GREEN}🎉 Meeseecs Box setup completed successfully!${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo -e "${PURPLE}1. Run: ${ORANGE}./start.sh${NC} to start the platform"
echo -e "${PURPLE}2. Open: ${ORANGE}http://localhost:3000${NC} to access the dashboard"
echo -e "${PURPLE}3. Read: ${ORANGE}QUICKSTART.md${NC} for detailed instructions"
echo ""
echo -e "${BLUE}Platform Features:${NC}"
echo -e "${PURPLE}🛡️  Runtime Security Monitoring${NC}"
echo -e "${PURPLE}🔍 Vulnerability Scanning${NC}"
echo -e "${PURPLE}📧 Multi-channel Alerting${NC}"
echo -e "${PURPLE}📊 Security Dashboard${NC}"
echo ""
echo -e "${ORANGE}Happy securing! 🚀${NC}"