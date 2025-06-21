#!/bin/bash

# Database Seeder Script
# Pastikan GOOGLE_PROJECT_ID sudah diset di environment

echo "🚀 Starting Database Seeder..."
echo ""

# Check if GOOGLE_PROJECT_ID is set
if [ -z "$GOOGLE_PROJECT_ID" ]; then
    echo "❌ Error: GOOGLE_PROJECT_ID environment variable is not set"
    echo "Please set it first:"
    echo "export GOOGLE_PROJECT_ID=\"your-project-id\""
    exit 1
fi

echo "✅ Using Google Project ID: $GOOGLE_PROJECT_ID"
echo ""

# Change to the eventual directory
cd "$(dirname "$0")/.."

# Run the seeder
echo "🌱 Running database seeder..."
go run cmd/seeder/main.go

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Database seeding completed successfully!"
    echo ""
    echo "📊 Summary:"
    echo "   - Cars: 5,000 units"
    echo "   - Hotel Rooms: 1,500 rooms"
    echo "   - Train Seats: 20,000 seats"
    echo "   - Total: 26,500 records"
else
    echo ""
    echo "❌ Database seeding failed!"
    exit 1
fi 