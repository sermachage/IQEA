# Local Development & Deployment Guide

## Overview

This guide covers local development and deployment for the Demographic Intelligence API. For production deployment to cloud platforms, please configure them separately through their respective web interfaces.

## Local Development

### Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Git

### Initial Setup

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/demographic-api.git
cd demographic-api
```

2. **Install Go dependencies**
```bash
go mod download
go mod tidy
```

3. **Create `.env` file**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. **Set up PostgreSQL locally**

```bash
# Option 1: PostgreSQL on your system
createdb iqea

# Option 2: Using Docker container
docker run --name iqea-postgres \
  -e POSTGRES_DB=iqea \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres:15-alpine
```

5. **Generate test data**
```bash
go run . generate -count 2026 -output profiles.json
```

6. **Seed the database**
```bash
go run . seed -file profiles.json
```

7. **Start the development server**
```bash
go run .
```

The server will start on `http://localhost:8080`

### Environment Variables

The following environment variables are required:

- `DATABASE_URL` - PostgreSQL connection string (format: `postgres://user:password@host:port/dbname`)
- `PORT` - Server port (default: 8080)
- `PROFILE_SEED_FILE` - Optional path to seed file (default: profiles.json)

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestName ./...

# Run with verbose output
go test -v ./...
```

### Building a Binary

```bash
# Standard build
go build -o demographic-api

# Optimized for size
go build -ldflags="-s -w" -o demographic-api

# With version information
go build -ldflags="-s -w -X main.version=1.0.0" -o demographic-api
```

### Running the Binary Locally

```bash
# Set environment variables
export DATABASE_URL="postgres://user:password@localhost:5432/iqea"
export PORT=8080

# Run
./demographic-api
```

## Testing Local Endpoints

Once the server is running, test the API:

```bash
# Health check
curl http://localhost:8080/health

# Get all profiles
curl http://localhost:8080/api/profiles

# Filter profiles
curl "http://localhost:8080/api/profiles?gender=male&country_id=NG"

# Search with natural language
curl "http://localhost:8080/api/profiles/search?q=young+males+from+nigeria"

# Test pagination
curl "http://localhost:8080/api/profiles?page=2&limit=20"

# Test sorting
curl "http://localhost:8080/api/profiles?sort_by=age&order=desc"
```

## Troubleshooting

### Application won't start
- Verify Go version: `go version` (should be 1.21+)
- Check PostgreSQL is running: `psql -U postgres -h localhost`
- Verify `.env` file exists and has correct `DATABASE_URL`
- Check for compilation errors: `go build`

### Database connection issues
- Verify `DATABASE_URL` format: `postgres://user:password@host:port/dbname`
- Test connection: `psql "your-connection-string"`
- Ensure PostgreSQL service is running
- Check port 5432 is not already in use

### No profiles in database
- Run seed command: `go run . seed -file profiles.json`
- Verify `profiles.json` file exists and is properly formatted
- Check for database permission issues: `psql -l`

### Port already in use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process (if needed)
kill -9 <PID>

# Or use a different port
export PORT=3000
go run .
```

## Development Workflow

1. Make changes to source files
2. Run tests: `go test ./...`
3. Restart server: `go run .`
4. Test endpoints with curl or Postman
5. Commit changes following Conventional Commits format
6. Push to your branch

## Performance Tips for Local Development

- Use indexed queries for better performance
- Batch database operations when possible
- Monitor resource usage: `top`, `htop`
- Check database query logs for slow queries
- Use connection pooling (configured in application)

## Cloud Deployment

To deploy to cloud platforms (Vercel, Railway, AWS, etc.):

1. Configure platform-specific deployment files/settings through their web interfaces
2. Set environment variables in your cloud platform dashboard
3. Configure database connections (PostgreSQL) through the platform
4. Test endpoints after deployment
5. Monitor logs and metrics through the platform's dashboard

**Note**: All cloud platform configuration files have been removed. Configure your desired platform through their web interfaces and follow their deployment guides.
