# API Testing Guide

## Quick Start Testing

### 1. Local Testing Setup

```bash
# Install dependencies
go mod download

# Start the server
go run .

# Server will start on http://localhost:8080
```

### 2. Generate Test Data (if needed)

```go
// Generate 2026 test profiles
go run testdata.go -generate profiles.json 2026
```

### 3. Seed Database

```bash
# From the root directory
go run . seed -file profiles.json
```

## API Endpoint Tests

### Test 1: Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

### Test 2: Get All Profiles (No Filters)

```bash
curl http://localhost:8080/api/profiles
```

Expected response:
```json
{
  "status": "success",
  "page": 1,
  "limit": 10,
  "total": 2026,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Emmanuel Okafor",
      "gender": "male",
      "gender_probability": 0.99,
      "age": 34,
      "age_group": "adult",
      "country_id": "NG",
      "country_name": "Nigeria",
      "country_probability": 0.85,
      "created_at": "2026-04-01T12:00:00Z"
    }
  ]
}
```

### Test 3: Filter by Gender

```bash
curl "http://localhost:8080/api/profiles?gender=male"
curl "http://localhost:8080/api/profiles?gender=female"
```

### Test 4: Filter by Country

```bash
curl "http://localhost:8080/api/profiles?country_id=NG"
curl "http://localhost:8080/api/profiles?country_id=KE"
curl "http://localhost:8080/api/profiles?country_id=GH"
```

### Test 5: Filter by Age Group

```bash
curl "http://localhost:8080/api/profiles?age_group=adult"
curl "http://localhost:8080/api/profiles?age_group=teenager"
curl "http://localhost:8080/api/profiles?age_group=child"
curl "http://localhost:8080/api/profiles?age_group=senior"
```

### Test 6: Filter by Age Range

```bash
# People aged 25-40
curl "http://localhost:8080/api/profiles?min_age=25&max_age=40"

# People over 30
curl "http://localhost:8080/api/profiles?min_age=30"

# People under 25
curl "http://localhost:8080/api/profiles?max_age=25"
```

### Test 7: Combined Filters (AND logic)

```bash
# Males from Nigeria
curl "http://localhost:8080/api/profiles?gender=male&country_id=NG"

# Adult females from Kenya
curl "http://localhost:8080/api/profiles?gender=female&age_group=adult&country_id=KE"

# Females over 25 from Uganda
curl "http://localhost:8080/api/profiles?gender=female&min_age=25&country_id=UG"

# Adult males above 30 from Nigeria with high confidence
curl "http://localhost:8080/api/profiles?gender=male&age_group=adult&min_age=30&country_id=NG&min_gender_probability=0.9"
```

### Test 8: Sorting

```bash
# Sort by age ascending
curl "http://localhost:8080/api/profiles?sort_by=age&order=asc"

# Sort by age descending (oldest first)
curl "http://localhost:8080/api/profiles?sort_by=age&order=desc"

# Sort by creation date
curl "http://localhost:8080/api/profiles?sort_by=created_at&order=asc"

# Sort by gender probability
curl "http://localhost:8080/api/profiles?sort_by=gender_probability&order=desc"
```

### Test 9: Pagination

```bash
# First page (default)
curl "http://localhost:8080/api/profiles?page=1&limit=10"

# Second page
curl "http://localhost:8080/api/profiles?page=2&limit=10"

# 20 results per page
curl "http://localhost:8080/api/profiles?page=1&limit=20"

# Maximum allowed (50 per page)
curl "http://localhost:8080/api/profiles?page=1&limit=50"

# Last page calculation: total / limit
curl "http://localhost:8080/api/profiles?page=41&limit=50"
```

### Test 10: Natural Language Search

```bash
# Young males from Nigeria
curl "http://localhost:8080/api/profiles/search?q=young+males+from+nigeria"

# Females above 30
curl "http://localhost:8080/api/profiles/search?q=females+above+30"

# Adult males from Kenya
curl "http://localhost:8080/api/profiles/search?q=adult+males+from+kenya"

# Teenagers from Uganda
curl "http://localhost:8080/api/profiles/search?q=teenagers+from+uganda"

# People from Cameroon
curl "http://localhost:8080/api/profiles/search?q=people+from+cameroon"

# Young people
curl "http://localhost:8080/api/profiles/search?q=young+people"

# Males above 25
curl "http://localhost:8080/api/profiles/search?q=males+above+25"

# Female teenagers
curl "http://localhost:8080/api/profiles/search?q=female+teenagers"

# Senior citizens from South Africa
curl "http://localhost:8080/api/profiles/search?q=senior+citizens+from+south+africa"
```

### Test 11: Search with Pagination

```bash
# First page of search results
curl "http://localhost:8080/api/profiles/search?q=adult+males&page=1&limit=10"

# Second page
curl "http://localhost:8080/api/profiles/search?q=adult+males&page=2&limit=20"
```

### Test 12: Error Cases

#### Invalid Gender
```bash
curl "http://localhost:8080/api/profiles?gender=unknown"
# Expected: 422 Unprocessable Entity
# {"status":"error","message":"Invalid gender value"}
```

#### Invalid Age Group
```bash
curl "http://localhost:8080/api/profiles?age_group=invalid"
# Expected: 422 Unprocessable Entity
# {"status":"error","message":"Invalid age_group value"}
```

#### Invalid Limit
```bash
curl "http://localhost:8080/api/profiles?limit=100"
# Expected: 422 Unprocessable Entity
# {"status":"error","message":"Invalid limit value (must be 1-50)"}
```

#### Invalid Sort By
```bash
curl "http://localhost:8080/api/profiles?sort_by=invalid"
# Expected: 422 Unprocessable Entity
# {"status":"error","message":"Invalid sort_by value"}
```

#### Invalid Order
```bash
curl "http://localhost:8080/api/profiles?order=invalid"
# Expected: 422 Unprocessable Entity
# {"status":"error","message":"Invalid order value"}
```

#### Unparseable Query
```bash
curl "http://localhost:8080/api/profiles/search?q=xyz+qwerty+asdf"
# Expected: 400 Bad Request
# {"status":"error","message":"Unable to interpret query"}
```

#### Missing Search Query
```bash
curl "http://localhost:8080/api/profiles/search"
# Expected: 400 Bad Request
# {"status":"error","message":"Missing or empty parameter"}
```

### Test 13: CORS Headers

```bash
# Test with Origin header
curl -H "Origin: http://example.com" \
     -H "Access-Control-Request-Method: GET" \
     -H "Access-Control-Request-Headers: content-type" \
     -X OPTIONS \
     http://localhost:8080/api/profiles

# Check response headers
# Access-Control-Allow-Origin: *
# Access-Control-Allow-Methods: GET, POST, OPTIONS
# Access-Control-Allow-Headers: Content-Type
```

## Automated Testing

### Run Unit Tests

```bash
go test -v
```

### Run Tests with Coverage

```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Load Testing

```bash
# Using Apache Bench
ab -n 1000 -c 10 http://localhost:8080/api/profiles

# Using wrk (requires installation)
wrk -t12 -c400 -d30s http://localhost:8080/api/profiles
```

## Cross-Network Testing

### Test from Different Machines

1. **Find your IP address**:
```bash
# Linux/Mac
ifconfig | grep "inet " | grep -v 127.0.0.1

# Windows
ipconfig
```

2. **Test from another machine on the same network**:
```bash
curl http://YOUR_IP:8080/api/profiles
```

3. **Test from Docker container**:
```bash
docker run -it --network host curlimages/curl http://localhost:8080/api/profiles
```

## Response Time Testing

### Measure Query Performance

```bash
# Using curl with timing
curl -w "Total: %{time_total}s, Connect: %{time_connect}s, Starttransfer: %{time_starttransfer}s\n" \
     -o /dev/null -s \
     "http://localhost:8080/api/profiles?gender=male&country_id=NG&min_age=25"
```

### Expected Performance Metrics
- Simple query: < 100ms
- Complex filters: < 150ms
- Large limit: < 200ms
- Pagination queries: < 100ms

## Testing Checklist

- [ ] Health endpoint returns 200
- [ ] GET /api/profiles returns all profiles
- [ ] Gender filter works (male/female)
- [ ] Age group filter works (child/teenager/adult/senior)
- [ ] Country ID filter works
- [ ] Min/max age filters work
- [ ] Probability threshold filters work
- [ ] Sort by age works (asc/desc)
- [ ] Sort by created_at works
- [ ] Sort by gender_probability works
- [ ] Pagination works with different page sizes
- [ ] Page limit cap of 50 is enforced
- [ ] Multiple filters combine with AND logic
- [ ] Natural language query "young males from nigeria" works
- [ ] Natural language query "females above 30" works
- [ ] Natural language query "adult males from kenya" works
- [ ] Invalid queries return error
- [ ] CORS headers present in response
- [ ] Timestamps are in UTC ISO 8601 format
- [ ] IDs are valid UUIDs
- [ ] Response structure matches specification
- [ ] Error responses have correct structure
- [ ] All status codes correct (200, 400, 422, 500)

## Database Verification

```sql
-- Connect to database
psql -U postgres -d iqea

-- Check profile count
SELECT COUNT(*) FROM profiles;

-- Verify schema
\d profiles

-- Check indexes
\di profiles*

-- Sample query performance
EXPLAIN ANALYZE 
SELECT * FROM profiles 
WHERE gender = 'male' AND country_id = 'NG' 
LIMIT 10;
```

## Troubleshooting Tests

### Query returns no results
- Verify data is seeded: `SELECT COUNT(*) FROM profiles;`
- Check filter values: `SELECT DISTINCT gender FROM profiles;`
- Verify country codes: `SELECT DISTINCT country_id FROM profiles;`

### Slow response times
- Check database indexes: `\di` in psql
- Monitor system resources: `top` or Task Manager
- Profile the query: `EXPLAIN ANALYZE`

### CORS errors in browser
- Verify header is present: `curl -i http://localhost:8080/api/profiles`
- Check browser console for full error
- Test with OPTIONS request

### 422 errors
- Validate all query parameter types
- Check for case sensitivity issues
- Verify enum values are exact
