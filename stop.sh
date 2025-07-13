#!/bin/bash

echo "🛑 Stopping Meeseecs Box Security Platform"
echo "=========================================="

cd docker-compose

# Use docker compose (v2) or docker-compose (v1)
if docker compose version &> /dev/null; then
    docker compose down
elif command -v docker-compose &> /dev/null; then
    docker-compose down
else
    echo "❌ Docker Compose not found"
    exit 1
fi

echo "✅ All services stopped"
