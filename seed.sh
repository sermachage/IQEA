#!/bin/bash
# Database seeding script for Insighta Labs Demographic API

# Configuration
DATABASE_URL=${DATABASE_URL:-"postgres://user:password@localhost:5432/iqea"}
PROFILE_FILE=${1:-"profiles.json"}

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Insighta Labs Database Seeder ===${NC}"
echo ""

# Check if profile file exists
if [ ! -f "$PROFILE_FILE" ]; then
    echo -e "${RED}Error: Profile file '$PROFILE_FILE' not found${NC}"
    exit 1
fi

echo -e "${BLUE}Seeding database with profiles from: $PROFILE_FILE${NC}"
echo "Database URL: $DATABASE_URL"
echo ""

# Run the Go seeder
go run . seed -file "$PROFILE_FILE"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Database seeding completed successfully!${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}✗ Database seeding failed${NC}"
    exit 1
fi
