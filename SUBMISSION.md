# Submission Guide

## Pre-Submission Checklist

- [ ] All code pushed to GitHub
- [ ] README complete with parsing approach and limitations documented
- [ ] Database schema implemented with UUID v7 and all required fields
- [ ] All 2026 profiles seeded into database
- [ ] GET /api/profiles endpoint working with all filters
- [ ] GET /api/profiles/search endpoint working with natural language parsing
- [ ] Pagination works (page, limit with max 50)
- [ ] Sorting works (age, created_at, gender_probability)
- [ ] Natural language parser handles all required mappings
- [ ] CORS headers present (Access-Control-Allow-Origin: *)
- [ ] All timestamps in UTC ISO 8601 format
- [ ] All IDs are UUID v7 format
- [ ] Error responses follow specification
- [ ] Tested from multiple networks
- [ ] Performance acceptable (< 200ms for typical queries)

## Submission Steps

### Step 1: Local Testing (24 hours before submission)

```bash
# Setup
export DATABASE_URL="postgres://user:password@localhost:5432/iqea"
go mod download

# Generate test data if needed
go run . generate -count 2026 -output profiles.json

# Seed database
go run . seed -file profiles.json

# Verify seeding
# In psql:
# SELECT COUNT(*) FROM profiles;
# Should return: 2026

# Start server
go run .

# Test all endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/profiles
curl "http://localhost:8080/api/profiles?gender=male&country_id=NG"
curl "http://localhost:8080/api/profiles/search?q=young+males+from+nigeria"

# Test from another terminal/machine
curl http://YOUR_IP:8080/health
```

### Step 2: Deploy to Railway (48 hours before submission)

```bash
# 1. Push to GitHub
git add .
git commit -m "Complete demographic API with NLP search"
git push origin main

# 2. Login to Railway
railway login

# 3. Link to your Railway project
railway link

# 4. Check database is set up
railway variables
# DATABASE_URL should be present

# 5. Deploy
railway up

# 6. Get your production URL
railway open
# Copy the URL (e.g., https://demographic-api-prod.railway.app)

# 7. Seed production database
railway run go run . seed -file profiles.json

# 8. Verify production deployment
curl https://your-railway-url.railway.app/health
curl https://your-railway-url.railway.app/api/profiles
```

### Step 3: Final Testing (24 hours before submission)

**Test from your local machine:**
```bash
PROD_URL="https://your-railway-url.railway.app"

# Health check
curl $PROD_URL/health

# Get all profiles
curl $PROD_URL/api/profiles

# Test filtering
curl "$PROD_URL/api/profiles?gender=male&country_id=NG"
curl "$PROD_URL/api/profiles?age_group=adult&min_age=25"

# Test sorting
curl "$PROD_URL/api/profiles?sort_by=age&order=desc&limit=5"

# Test pagination
curl "$PROD_URL/api/profiles?page=1&limit=10"
curl "$PROD_URL/api/profiles?page=2&limit=10"

# Test natural language search
curl "$PROD_URL/api/profiles/search?q=young+males+from+nigeria"
curl "$PROD_URL/api/profiles/search?q=females+above+30"
curl "$PROD_URL/api/profiles/search?q=adult+males+from+kenya"
curl "$PROD_URL/api/profiles/search?q=teenagers+from+uganda"

# Test error cases
curl "$PROD_URL/api/profiles?gender=invalid"  # Should be 422
curl "$PROD_URL/api/profiles?limit=100"       # Should be 422
curl "$PROD_URL/api/profiles/search?q=xyz+qwerty+asdf"  # Should be 400
curl "$PROD_URL/api/profiles/search"          # Should be 400

# Test CORS
curl -H "Origin: http://example.com" $PROD_URL/api/profiles
# Check for Access-Control-Allow-Origin: * header
```

**Test from different network (mobile hotspot, office, friend's network):**
```bash
# Repeat same tests to verify accessibility
curl https://your-railway-url.railway.app/api/profiles
```

### Step 4: Verify Documentation

Check these files exist and are complete:
- [ ] **README.md** - Contains:
  - Overview of API
  - All endpoints documented with examples
  - Natural language parser explanation
  - Supported keywords and country codes
  - Limitations and edge cases
  - Example queries and their mappings
  - Error responses
  - Database schema
  - Setup instructions

- [ ] **DEPLOYMENT.md** - Contains:
  - Railway deployment steps
  - Environment configuration
  - Troubleshooting guide
  - Monitoring and scaling info

- [ ] **TESTING.md** - Contains:
  - All test scenarios
  - Example curl commands for each endpoint
  - Error case testing
  - Performance testing
  - Load testing

- [ ] **BUILD.md** - Contains:
  - Build instructions
  - Docker deployment
  - Local development setup
  - Debugging guide

### Step 5: Code Review Checklist

```bash
# Check code quality
go fmt ./...
go vet ./...

# Run tests
go test -v
go test -cover

# Build release binary
go build -ldflags="-s -w" -o demographic-api

# Check for security issues
go list -m all | grep -i security
```

### Step 6: Final Documentation Check

Verify README includes all required information:

1. **Parsing Approach** ✓
   - Explanation of rule-based keyword matching
   - No AI/LLM explanation
   - How keywords map to filters
   - Combination logic (AND)

2. **Supported Keywords** ✓
   - Gender: male, female, man, woman, men, women
   - Age groups: child, teenager, adult, senior
   - Age qualifiers: young (16-24), middle (25-50), old (50+)
   - Age constraints: above/over, below/under
   - Countries: All 24+ African nations

3. **Limitations** ✓
   - OR logic not supported
   - Negation not supported
   - Complex math expressions not supported
   - Probability thresholds not supported
   - Multiple countries not supported
   - Qualifiers vs explicit ages (qualifier wins)
   - Ambiguous gender (female wins)

4. **Example Queries** ✓
   - "young males from nigeria" → gender=male, min_age=16, max_age=24, country_id=NG
   - "females above 30" → gender=female, min_age=30
   - "adult males from kenya" → gender=male, age_group=adult, country_id=KE
   - And more

### Step 7: Prepare Submission

Create a summary document with:

1. **GitHub Repository URL**
   ```
   https://github.com/yourusername/demographic-api
   ```

2. **Production Server URL**
   ```
   https://your-railway-url.railway.app
   ```

3. **Key Statistics**
   - Total profiles seeded: 2026
   - Database records: [actual count from production]
   - Response time (average): [measure with curl]
   - Database indexes: 7 (gender, age, age_group, country_id, gender_probability, country_probability, created_at)

4. **Tested Endpoints**
   - GET /api/profiles - ✓
   - GET /api/profiles/search - ✓
   - GET /health - ✓
   - OPTIONS /api/profiles - ✓ (CORS)

5. **Features Implemented**
   - Basic filtering ✓
   - Combined filters ✓
   - Sorting ✓
   - Pagination ✓
   - Natural language parsing ✓
   - Error handling ✓
   - CORS support ✓

## Submission Format

Submit the following to the grading platform:

1. **GitHub Repository Link**
   - Ensure repository is public
   - Code is clean and well-organized
   - Commit history is meaningful

2. **Production Server URL**
   - Server must be live and accessible
   - All endpoints responding correctly
   - Database is seeded

3. **README.md Content**
   - Natural language parsing explanation
   - Limitations documentation
   - Example queries and mappings

4. **Test Results Summary**
   - All endpoints tested ✓
   - Pagination verified ✓
   - Sorting verified ✓
   - Natural language search tested ✓
   - Error handling verified ✓
   - CORS headers confirmed ✓
   - Tested from multiple networks ✓

## Scoring Breakdown (100 points total)

| Component | Points | Notes |
|-----------|--------|-------|
| Filtering Logic | 20 | Gender, age_group, country_id, age ranges |
| Combined Filters | 15 | AND logic for multiple filters |
| Pagination | 15 | page, limit (max 50) |
| Sorting | 10 | age, created_at, gender_probability |
| Natural Language Parsing | 20 | Rule-based parsing with keyword mapping |
| README Explanation | 10 | Parsing approach + limitations |
| Query Validation | 5 | Proper error responses |
| Performance | 5 | Sub-200ms queries, proper indexing |
| **TOTAL** | **100** | |

## Common Issues Before Submission

### "Database connection refused"
```bash
# Check Railway database is configured
railway variables

# Check DATABASE_URL format
echo $DATABASE_URL
```

### "No profiles in database"
```bash
# Verify seeding ran
railway run psql -c "SELECT COUNT(*) FROM profiles;"

# If empty, seed again
railway run go run . seed -file profiles.json
```

### "CORS errors in testing"
```bash
# Verify header present
curl -H "Origin: http://example.com" \
     https://your-railway-url.railway.app/api/profiles -i
```

### "Slow response times"
```bash
# Check indexes exist
railway run psql -c "\di profiles*"

# Check query performance
railway run psql -c "EXPLAIN ANALYZE SELECT * FROM profiles LIMIT 10;"
```

## Last-Minute Checklist (1 hour before deadline)

- [ ] Server is running and responding
- [ ] Database has 2026 profiles
- [ ] Health endpoint returns 200
- [ ] /api/profiles returns data
- [ ] /api/profiles/search works
- [ ] CORS headers present
- [ ] Error responses correct
- [ ] README complete
- [ ] No console errors
- [ ] Performance is acceptable
- [ ] All tests passing locally
- [ ] GitHub repository is public
- [ ] No sensitive data in repository
- [ ] All documentation is present
- [ ] Production URL is accessible

## Support

If you encounter issues:

1. Check the [TESTING.md](TESTING.md) for endpoint tests
2. Review [DEPLOYMENT.md](DEPLOYMENT.md) for troubleshooting
3. Check [BUILD.md](BUILD.md) for build/development issues
4. Review database logs: `railway logs`
5. Test locally first before assuming production issue

## Success Criteria

Your submission will be successful if:

1. ✓ Server is live and publicly accessible on Railway
2. ✓ All 2026 profiles are seeded in the database
3. ✓ All API endpoints respond correctly
4. ✓ Natural language parsing works for example queries
5. ✓ Filtering, sorting, and pagination work correctly
6. ✓ CORS headers are present
7. ✓ Error responses follow specification
8. ✓ README documents parsing approach and limitations
9. ✓ Response times are acceptable (< 200ms)
10. ✓ Database is properly indexed
