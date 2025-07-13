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
    # Use docker compose (v2) or docker-compose (v1)
    if docker compose version &> /dev/null; then
        docker compose logs -f $1
    elif command -v docker-compose &> /dev/null; then
        docker-compose logs -f $1
    else
        echo "❌ Docker Compose not found"
        exit 1
    fi
fi
