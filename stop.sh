#!/bin/bash

echo "🛑 Stopping Meeseecs Box Security Platform"
echo "=========================================="

cd docker-compose

# Use docker compose (v2) or docker-compose (v1)
if command -v docker-compose &> /dev/null; then
    docker-compose down
else
    docker compose down
fi

echo "✅ All services stopped"
