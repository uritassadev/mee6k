#!/bin/bash

echo "🛑 Stopping Meeseecs Box Security Platform"
echo "=========================================="

cd docker-compose

docker compose down

echo "✅ All services stopped"
