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
sleep 30

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
