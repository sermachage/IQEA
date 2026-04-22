# Build & Run Guide

## Local Development

### Prerequisites
- Go 1.21 or later
- PostgreSQL 12 or later
- Git

### Setup Steps

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

4. **Configure PostgreSQL**
```bash
# Option 1: Local PostgreSQL
createdb iqea

# Option 2: Docker PostgreSQL
docker run --name iqea-postgres \
  -e POSTGRES_DB=iqea \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres:15-alpine
```

5. **Generate test data (if you don't have 2026 profiles)**
```bash
go run . generate -count 2026 -output profiles.json
```

6. **Seed the database**
```bash
go run . seed -file profiles.json
```

7. **Start the server**
```bash
go run .
```

The server will start on `http://localhost:8080`

## Building for Production

### Build Binary

```bash
# Standard build
go build -o demographic-api

# Optimized for size
go build -ldflags="-s -w" -o demographic-api

# With version information
go build -ldflags="-s -w -X main.version=1.0.0" -o demographic-api
```

### Cross-platform builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o demographic-api-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o demographic-api-darwin

# Windows
GOOS=windows GOARCH=amd64 go build -o demographic-api.exe
```

### Run Binary

```bash
# Set environment variables
export DATABASE_URL="postgres://user:password@localhost:5432/iqea"
export PORT=8080

# Run
./demographic-api
```

## Docker Deployment

### Build Docker Image

```bash
# Build
docker build -t demographic-api:latest .

# Tag
docker tag demographic-api:latest demographic-api:v1.0.0
```

### Run Docker Container

```bash
# With environment file
docker run --env-file .env \
  -p 8080:8080 \
  --name demographic-api \
  demographic-api:latest

# With environment variables
docker run -e DATABASE_URL="postgres://user:password@host:5432/iqea" \
  -e PORT=8080 \
  -p 8080:8080 \
  --name demographic-api \
  demographic-api:latest

# With volume mount for profiles
docker run -e DATABASE_URL="postgres://user:password@host:5432/iqea" \
  -v $(pwd)/profiles.json:/app/profiles.json \
  -p 8080:8080 \
  demographic-api:latest
```

### Docker Compose (with PostgreSQL)

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: iqea
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/iqea
      PORT: 8080
    depends_on:
      - postgres
    command: sh -c "go run . seed -file profiles.json && ./demographic-api"

volumes:
  postgres_data:
```

Then run:
```bash
docker-compose up
```

## Railway Deployment

### Using Railway CLI

1. **Install Railway CLI**
```bash
npm install -g @railway/cli
```

2. **Login**
```bash
railway login
```

3. **Create new project**
```bash
railway init
```

4. **Add PostgreSQL**
```bash
railway add
# Select PostgreSQL
```

5. **Deploy**
```bash
railway up
```

6. **Seed database**
```bash
railway run go run . seed -file profiles.json
```

### Using GitHub Integration

1. Push code to GitHub
2. Connect GitHub to Railway
3. Railway auto-deploys on push
4. Add environment variables in Railway dashboard
5. Seed database manually or via GitHub Actions

## Testing Commands

### Unit Tests
```bash
go test -v
go test -cover
go test -coverprofile=coverage.out && go tool cover -html=coverage.out
```

### API Testing
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test profiles endpoint
curl http://localhost:8080/api/profiles

# Test with filters
curl "http://localhost:8080/api/profiles?gender=male&country_id=NG"

# Test search
curl "http://localhost:8080/api/profiles/search?q=young+males+from+nigeria"
```

### Load Testing
```bash
# Apache Bench
ab -n 1000 -c 100 http://localhost:8080/api/profiles

# wrk
wrk -t12 -c400 -d30s http://localhost:8080/api/profiles

# Using Go's internal testing
go test -bench=. -benchmem
```

## Debugging

### Enable Verbose Logging
```bash
# Run with debug output
go run -v .
```

### Database Connection Issues
```bash
# Test PostgreSQL connection
psql -U postgres -d iqea -h localhost

# Check environment variables
env | grep DATABASE_URL
```

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
export PORT=8081
go run .
```

## Useful Commands

### View Database
```bash
# Connect to database
psql -U postgres -d iqea

# List tables
\dt

# Count profiles
SELECT COUNT(*) FROM profiles;

# View profile schema
\d profiles

# View indexes
\di profiles*

# Sample query
SELECT * FROM profiles LIMIT 5;
```

### Clean Build
```bash
go clean
rm -f demographic-api
go build
```

### Format Code
```bash
go fmt ./...
```

### Vet Code
```bash
go vet ./...
```

### Dependencies
```bash
# List direct dependencies
go list -m all

# Download all dependencies
go mod download

# Verify integrity
go mod verify

# Clean unused
go mod tidy
```

## Environment Variables

Required:
- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - Server port (default: 8080)

Optional:
- `PROFILE_SEED_FILE` - Path to profiles.json

Example:
```
DATABASE_URL=postgres://postgres:password@localhost:5432/iqea
PORT=8080
PROFILE_SEED_FILE=./profiles.json
```

## Common Issues & Solutions

### Database Connection Refused
```bash
# Check PostgreSQL is running
psql -U postgres -c "SELECT 1"

# Verify DATABASE_URL format
# Format: postgres://username:password@host:port/dbname
```

### Out of Memory During Seeding
```bash
# Reduce batch size by modifying seeder.go
# Or seed in chunks using multiple files
```

### Port Already in Use
```bash
# Use a different port
export PORT=9000
go run .
```

### Slow Queries
```sql
-- Check indexes exist
SELECT * FROM pg_indexes WHERE tablename = 'profiles';

-- Check query plan
EXPLAIN ANALYZE SELECT * FROM profiles WHERE gender = 'male' LIMIT 10;
```

## Performance Optimization

### Database Tuning
```sql
-- Increase buffer cache
ALTER SYSTEM SET shared_buffers = '256MB';

-- Increase work memory
ALTER SYSTEM SET work_mem = '64MB';

-- Increase effective cache size
ALTER SYSTEM SET effective_cache_size = '2GB';

-- Reload configuration
SELECT pg_reload_conf();
```

### Connection Pooling
Consider using PgBouncer for connection pooling in production:
```
[databases]
iqea = host=localhost port=5432 dbname=iqea

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 25
```

## Monitoring

### Health Checks
```bash
# Simple health check
curl http://localhost:8080/health

# With response time
curl -w "%{time_total}s" http://localhost:8080/health
```

### Logs
```bash
# Save logs to file
go run . 2>&1 | tee app.log

# Filter logs
grep "error" app.log
```

## Deployment Checklist

- [ ] Code pushed to GitHub
- [ ] All tests passing locally
- [ ] `.env` file configured with production database
- [ ] Database migrations run
- [ ] Profiles seeded (2026 records)
- [ ] Health endpoint tested
- [ ] All API endpoints tested from multiple networks
- [ ] CORS headers verified
- [ ] Error responses validated
- [ ] Performance tested
- [ ] Backups configured
- [ ] Monitoring set up
- [ ] Documentation reviewed
