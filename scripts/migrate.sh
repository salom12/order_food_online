#!/bin/bash
set -e

echo "Running migrations..."
psql $DATABASE_URL -f migrations/001_create_products_table.sql
psql $DATABASE_URL -f migrations/002_create_orders_table.sql
echo "Migrations completed."
