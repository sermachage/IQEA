# IQEA Quick Deploy to Vercel (from ~)

Use these commands from your home directory:

```bash
cd ~/IQEA
go test ./...
npm i -g vercel
vercel login
vercel env add DATABASE_URL production
vercel env add DATABASE_URL preview
vercel --prod
```

Verify:

```bash
curl https://<your-vercel-domain>/health
curl "https://<your-vercel-domain>/api/profiles?page=1&limit=10"
```

Notes:
- Vercel detects this project as Go via `go.mod` and `vercel.json`.
- Seed production data from your machine/CI using production `DATABASE_URL`:

```bash
cd ~/IQEA
export DATABASE_URL="<your-production-db-url>"
go run . seed -file profiles.json
```

---

# Deployment Guide for Insighta Labs Demographic API on Railway

## Prerequisites

- GitHub account
- Railway account
- PostgreSQL database (can be provisioned by Railway)
- 2026 profiles JSON file

## Deployment Steps

### Step 1: Push to GitHub

1. Initialize or connect your Git repository:
```bash
git init
git add .
git commit -m "Initial commit: Demographic API"
git branch -M main
git remote add origin https://github.com/yourusername/demographic-api.git
git push -u origin main
```

2. Ensure your repository is public or Railway has access

### Step 2: Create Railway Project

1. Go to [Railway.app](https://railway.app)
2. Click "New Project"
3. Select "Deploy from GitHub"
4. Connect your GitHub account and select this repository
5. Railway will automatically detect the Go project

### Step 3: Configure Environment Variables

In your Railway project dashboard:

1. Add a PostgreSQL plugin:
   - Click "Add Service"
   - Select "Database"
   - Choose "PostgreSQL"
   - Railway will automatically provision and set `DATABASE_URL`

2. Add application variables:
   - `PORT`: 8080 (or your preferred port)

### Step 4: Deploy

1. Railway automatically deploys on every push to main branch
2. You can also manually trigger deployment:
   - Click "Deploy" on your Railway dashboard
   - Watch the build logs

### Step 5: Verify Deployment

Once deployed:

1. Get your service URL from the Railway dashboard
2. Test the health endpoint:
```bash
curl https://your-railway-url.railway.app/health
```

Expected response:
```json
{"status":"ok"}
```

### Step 6: Seed the Database

After deployment:

1. Use Railway's CLI to run the seeding job:
```bash
railway run go run . seed -file profiles.json
```

Or create a one-time job:

2. Connect to your Railway PostgreSQL database
3. Download your 2026 profiles JSON file
4. Run locally first to verify:
```bash
export DATABASE_URL="your-railway-database-url"
go run . seed -file profiles.json
```

5. Or add a GitHub Actions workflow to auto-seed on deploy:

Create `.github/workflows/seed.yml`:
```yaml
name: Seed Database

on:
  deployment

jobs:
  seed:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Seed database
        run: go run . seed -file profiles.json
        env:
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
```

### Step 7: Test Endpoints

Test from multiple networks:

```bash
# Health check
curl https://your-railway-url.railway.app/health

# Get all profiles
curl https://your-railway-url.railway.app/api/profiles

# Filter profiles
curl "https://your-railway-url.railway.app/api/profiles?gender=male&country_id=NG"

# Search with natural language
curl "https://your-railway-url.railway.app/api/profiles/search?q=young+males+from+nigeria"

# Test pagination
curl "https://your-railway-url.railway.app/api/profiles?page=2&limit=20"

# Test sorting
curl "https://your-railway-url.railway.app/api/profiles?sort_by=age&order=desc"
```

## Troubleshooting

### Application won't start
- Check logs: `railway logs`
- Verify `DATABASE_URL` is set correctly
- Ensure PostgreSQL service is running
- Check for Go compilation errors

### Database connection issues
- Verify `DATABASE_URL` format: `postgres://user:password@host:port/dbname`
- Check PostgreSQL service status
- Ensure IP whitelist includes Railway's egress IPs

### No profiles in database
- Check if seeding ran successfully
- Verify profiles.json is properly formatted
- Check for constraint violations (duplicate names)
- View seeding logs: `railway logs`

### CORS errors from client
- Verify `Access-Control-Allow-Origin: *` header is present
- Check browser console for exact error
- Test with `curl` to verify server is responding

### Performance issues
- Check query logs for slow queries
- Verify indexes are created:
  - idx_profiles_gender
  - idx_profiles_age
  - idx_profiles_age_group
  - idx_profiles_country_id
  - idx_profiles_gender_probability
  - idx_profiles_country_probability
  - idx_profiles_created_at
- Consider optimizing query filters
- Check database connection pool settings

## Monitoring

1. **Logs**: View in Railway dashboard under "Logs" tab
2. **Health Endpoint**: Continuously poll `/health` for uptime monitoring
3. **Metrics**: Railway provides CPU, Memory, and Network metrics
4. **Alerts**: Set up alerts in Railway dashboard for failures

## Scaling

1. **Database**: Railway can auto-scale PostgreSQL
2. **Application**: Increase replicas in railway.json:
```json
{
  "deploy": {
    "numReplicas": 2
  }
}
```
3. **Caching**: Consider adding Redis for caching frequent queries
4. **Load Balancing**: Railway handles automatic load balancing

## Security

1. **Environment Variables**: Never commit `.env` files
2. **Database URL**: Keep DATABASE_URL secret in Railway secrets
3. **CORS**: Currently allows all origins (`*`). Restrict in production if needed
4. **Input Validation**: All parameters are validated server-side
5. **SQL Injection**: Using parameterized queries to prevent injection

## Rollback

To rollback to a previous version:

1. In Railway dashboard, navigate to "Deployments"
2. Find the previous successful deployment
3. Click the three dots and select "Redeploy"
4. Confirm the rollback

## Continuous Deployment

1. Every push to main branch automatically deploys
2. Failed builds prevent deployment
3. Successful builds go live immediately
4. Zero-downtime deployments with Railway

## Cost Considerations

- **PostgreSQL**: First 10GB free, then pay-as-you-go
- **Application**: $5 minimum per month
- **Database egress**: Charged by Railway
- **Traffic**: Generally included in plan

## Accessing Railway Logs

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Link to project
railway link

# View logs
railway logs

# View env vars
railway variables
```

## Next Steps

1. Monitor performance and user feedback
2. Optimize queries based on usage patterns
3. Consider implementing caching for frequent queries
4. Add authentication if needed
5. Set up automated backups
