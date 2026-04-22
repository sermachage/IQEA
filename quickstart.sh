#!/bin/bash
# Quick start script for Demographic API

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Demographic Intelligence API${NC}"
echo -e "${BLUE}  Quick Start${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check dependencies
echo -e "${YELLOW}Checking dependencies...${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Go is installed${NC}"

if ! command -v psql &> /dev/null; then
    echo -e "${YELLOW}⚠ PostgreSQL client not found (install for local testing)${NC}"
else
    echo -e "${GREEN}✓ PostgreSQL client installed${NC}"
fi

echo ""
echo -e "${YELLOW}Setting up project...${NC}"

# Download dependencies
echo "Downloading Go dependencies..."
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies ready${NC}"

echo ""
echo -e "${YELLOW}Database Setup${NC}"

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo -e "${YELLOW}DATABASE_URL not set${NC}"
    echo "Using default: postgres://postgres:postgres@localhost:5432/iqea"
    export DATABASE_URL="postgres://postgres:postgres@localhost:5432/iqea"
fi

# Option to generate test data
read -p "Generate 2026 test profiles? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Generating profiles...${NC}"
    go run . generate -count 2026 -output profiles.json
    echo -e "${GREEN}✓ Profiles generated${NC}"
fi

# Option to seed database
read -p "Seed database? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Seeding database...${NC}"
    go run . seed -file profiles.json
    echo -e "${GREEN}✓ Database seeded${NC}"
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Setup Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Start the server:"
echo "   ${BLUE}go run .${NC}"
echo ""
echo "2. Test the API:"
echo "   ${BLUE}curl http://localhost:8080/health${NC}"
echo "   ${BLUE}curl http://localhost:8080/api/profiles${NC}"
echo "   ${BLUE}curl http://localhost:8080/api/profiles/search?q=young+males+from+nigeria${NC}"
echo ""
echo -e "${YELLOW}Documentation:${NC}"
echo "• API Docs: README.md"
echo "• Testing Guide: TESTING.md"
echo "• Build Guide: BUILD.md"
echo "• Deployment: DEPLOYMENT.md"
echo "• Submission: SUBMISSION.md"
echo ""
