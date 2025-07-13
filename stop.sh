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
