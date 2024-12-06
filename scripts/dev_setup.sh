#!/bin/bash
set -e

echo "Setting up development environment..."
docker-compose up -d
./scripts/migrate.sh
psql $DATABASE_URL -f scripts/seed_data.sql
echo "Development environment setup completed."
