# Demographic Intelligence API

A high-performance REST API for querying demographic data with advanced filtering, sorting, pagination, and natural language search capabilities.

## Overview

This API provides demographic intelligence for marketing teams, product teams, and growth analysts. It allows rapid querying of demographic profiles with support for complex filters, real-time search, and natural language queries.

## Features

- **Advanced Filtering**: Filter by gender, age group, country, and probability thresholds
- **Sorting**: Sort by age, creation date, or gender probability
- **Pagination**: Efficient pagination with configurable page size (max 50 per page)
- **Natural Language Search**: Parse plain English queries to filter demographics
- **High Performance**: Indexed database queries for sub-second response times
- **CORS Enabled**: Full cross-origin support for web clients

## Technology Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 12+
- **UUID**: UUID v7 for unique identifiers

## API Endpoints

### 1. GET /api/profiles - Get All Profiles

Returns a paginated list of profiles with support for filtering and sorting.

#### Query Parameters

**Filters:**
- `gender` - "male" or "female"
- `age_group` - "child", "teenager", "adult", or "senior"
- `country_id` - ISO 2-letter country code (e.g., "NG", "KE", "BJ")
- `min_age` - Minimum age (integer)
- `max_age` - Maximum age (integer)
- `min_gender_probability` - Minimum gender confidence score (0-1 float)
- `min_country_probability` - Minimum country confidence score (0-1 float)

**Sorting:**
- `sort_by` - "age", "created_at", or "gender_probability" (default: "created_at")
- `order` - "asc" or "desc" (default: "asc")

**Pagination:**
- `page` - Page number (default: 1)
- `limit` - Results per page (default: 10, max: 50)

#### Example Request

```
GET /api/profiles?gender=male&country_id=NG&min_age=25&sort_by=age&order=desc&page=1&limit=10
```

#### Success Response (200)

```json
{
  "status": "success",
  "page": 1,
  "limit": 10,
  "total": 2026,
  "data": [
    {
      "id": "b3f9c1e2-7d4a-4c91-9c2a-1f0a8e5b6d12",
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

### 2. GET /api/profiles/search - Natural Language Query

Search demographics using plain English queries. The parser converts natural language to filters automatically.

#### Query Parameters

- `q` - Plain English query string (required)
- `page` - Page number (default: 1)
- `limit` - Results per page (default: 10, max: 50)

#### Example Requests

```
GET /api/profiles/search?q=young+males+from+nigeria
GET /api/profiles/search?q=females+above+30
GET /api/profiles/search?q=adult+males+from+kenya
GET /api/profiles/search?q=teenagers+from+uganda
```

#### Success Response (200)

Same structure as /api/profiles endpoint.

#### Error Response (400)

```json
{
  "status": "error",
  "message": "Unable to interpret query"
}
```

### 3. GET /health - Health Check

Simple health check endpoint.

#### Response

```json
{"status":"ok"}
```

## Natural Language Parser

### Parsing Approach

The natural language parser uses **rule-based keyword matching** with no AI or LLM dependencies. It operates by:

1. Converting the query to lowercase for case-insensitive matching
2. Scanning for specific keywords in the query string
3. Mapping keywords to database filters
4. Combining multiple detected filters with AND logic

### Supported Keywords

#### Gender Keywords
- "male", "man", "men" → `gender=male`
- "female", "woman", "women" → `gender=female`

#### Age Group Keywords
- "child", "children" → `age_group=child`
- "teenager", "teenagers" → `age_group=teenager`
- "adult", "adults" → `age_group=adult`
- "senior", "seniors", "elderly" → `age_group=senior`

#### Age Qualifiers
- "young" → `min_age=16&max_age=24`
- "middle" → `min_age=25&max_age=50`
- "old" → `min_age=50&max_age=150`

#### Age Constraints
- "above X" or "over X" → `min_age=X`
- "below X" or "under X" → `max_age=X`

#### Supported Countries (24 African Nations)
Queries containing country names are automatically converted to ISO 2-letter codes:

- Nigeria → NG
- Ghana → GH
- Kenya → KE
- Uganda → UG
- Tanzania → TZ
- Angola → AO
- Cameroon → CM
- Senegal → SN
- Mali → ML
- Ivory Coast → CI
- Benin → BJ
- Burkina Faso → BF
- Burundi → BI
- Cape Verde → CV
- Central African Republic → CF
- Chad → TD
- Comoros → KM
- Congo → CG
- Democratic Republic of Congo → CD
- Djibouti → DJ
- Egypt → EG
- Equatorial Guinea → GQ
- Eritrea → ER
- Eswatini → SZ
- Ethiopia → ET
- Gabon → GA
- Gambia → GM
- Guinea → GN
- Guinea-Bissau → GW
- Lesotho → LS
- Liberia → LR
- Libya → LY
- Madagascar → MG
- Malawi → MW
- Mauritania → MR
- Mauritius → MU
- Morocco → MA
- Mozambique → MZ
- Namibia → NA
- Niger → NE
- Rwanda → RW
- São Tomé and Príncipe → ST
- Seychelles → SC
- Sierra Leone → SL
- Somalia → SO
- South Africa → ZA
- South Sudan → SS
- Sudan → SD
- Togo → TG
- Tunisia → TN
- Zambia → ZM
- Zimbabwe → ZW

### Example Query Mappings

| Query | Parsed Filters |
|-------|---|
| "young males from nigeria" | gender=male, min_age=16, max_age=24, country_id=NG |
| "females above 30" | gender=female, min_age=30 |
| "adult males from kenya" | gender=male, age_group=adult, country_id=KE |
| "teenagers from uganda" | age_group=teenager, country_id=UG |
| "male and female teenagers above 17" | age_group=teenager, min_age=17 |
| "people from cameroon" | country_id=CM |

### Parser Logic

The parser operates as follows:

1. **Gender Detection**: Searches for male/female keywords. If "female" appears, gender is set to female (ignores male keywords). Otherwise, if "male" keywords appear, gender is set to male.

2. **Age Group Detection**: First keyword found from the age group list is used (child, teenager, adult, senior).

3. **Age Constraint Detection**: Looks for "above", "over", "below", "under" keywords followed by numbers. Extracts the first numeric value found.

4. **Age Qualifier Detection**: If no explicit age constraint found, looks for qualifiers (young, middle, old) and sets min/max ages accordingly.

5. **Country Detection**: Searches query for country names and converts to ISO codes.

6. **Validation**: If no meaningful filters are detected, returns an error.

7. **Combination**: All detected filters are combined with AND logic for the query.

### Parser Limitations & Edge Cases

#### Not Supported

1. **OR Logic**: The parser does not support OR conditions. All filters are combined with AND.
   - Query: "males or females" → Returns error (contradictory conditions)
   - Workaround: Use two separate queries

2. **Negation**: The parser does not support negative queries.
   - Query: "not from Nigeria" → Not supported
   - Workaround: Use the direct API with explicit filters

3. **Complex Math**: No support for age ranges like "18-25" or "between 20 and 30".
   - Workaround: Use two queries (one for min_age, one for max_age) or use the direct API

4. **Probability Thresholds**: No support for natural language probability queries.
   - Query: "high confidence males" → Not supported
   - Workaround: Use direct API with min_gender_probability

5. **Multiple Countries**: Only the first country mentioned is used.
   - Query: "people from Nigeria or Ghana" → Only Nigeria parsed
   - Workaround: Use two separate queries

6. **Qualifiers with Explicit Ages**: If both a qualifier and explicit age are present, the qualifier is used.
   - Query: "young person age 45" → Uses young (16-24), ignores 45
   - Workaround: Rephrase query or use direct API

7. **Ambiguous Gender**: If both male and female keywords present, the parser prioritizes female.
   - Query: "male and female people" → Searches for females only
   - Workaround: Use two separate queries or use the direct API

#### Supported Edge Cases

1. **Case Insensitivity**: All keywords are case-insensitive.
   - "Young Males From NIGERIA" → Works correctly

2. **Extra Whitespace**: Multiple spaces and tabs are handled.
   - "young    males" → Works correctly

3. **Word Order**: Keywords can appear in any order.
   - "from nigeria young males" → Works correctly

4. **Partial Words**: Keywords are matched within larger words.
   - "female" matches "females"
   - "male" matches "males"

## Database Schema

```sql
CREATE TABLE profiles (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  gender VARCHAR(20) NOT NULL,
  gender_probability FLOAT NOT NULL,
  age INT NOT NULL,
  age_group VARCHAR(50) NOT NULL,
  country_id VARCHAR(2) NOT NULL,
  country_name VARCHAR(255) NOT NULL,
  country_probability FLOAT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  CHECK (gender IN ('male', 'female')),
  CHECK (gender_probability >= 0 AND gender_probability <= 1),
  CHECK (age >= 0 AND age <= 150),
  CHECK (age_group IN ('child', 'teenager', 'adult', 'senior')),
  CHECK (country_probability >= 0 AND country_probability <= 1)
);

-- Indexes for performance
CREATE INDEX idx_profiles_gender ON profiles(gender);
CREATE INDEX idx_profiles_age ON profiles(age);
CREATE INDEX idx_profiles_age_group ON profiles(age_group);
CREATE INDEX idx_profiles_country_id ON profiles(country_id);
CREATE INDEX idx_profiles_gender_probability ON profiles(gender_probability);
CREATE INDEX idx_profiles_country_probability ON profiles(country_probability);
CREATE INDEX idx_profiles_created_at ON profiles(created_at);
```

## Error Responses

### 400 Bad Request
Missing or empty required parameter:
```json
{
  "status": "error",
  "message": "Missing or empty parameter"
}
```

### 422 Unprocessable Entity
Invalid parameter type or value:
```json
{
  "status": "error",
  "message": "Invalid query parameters"
}
```

### 500 Internal Server Error
Server-side processing error:
```json
{
  "status": "error",
  "message": "Internal server error"
}
```

## Installation & Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

### Local Setup

1. **Clone the repository**
```bash
git clone <repo-url>
cd demographic-api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment**
```bash
# Create .env file
export DATABASE_URL="postgres://user:password@localhost:5432/iqea"
export PORT=8080
```

4. **Initialize database**
```bash
go run . -seed -file profiles.json
```

5. **Start the server**
```bash
go run main.go database.go models.go handlers.go nlparser.go seeder.go
```

The server will start on `http://localhost:8080`

## Deployment on Railway

### Steps

1. **Connect your GitHub repository** to Railway
2. **Set environment variables**:
   - `DATABASE_URL`: Your PostgreSQL connection string
   - `PORT`: 8080 (or your preferred port)
3. **Deploy**:
   - Railway automatically detects Go projects and builds them
   - The server starts with `go run .`

### Database Seeding on Railway

1. Create a seed job or migration script
2. Run the seeder before or after deployment:
```bash
go run . seed -file profiles.json
```

## Performance Characteristics

- **Query Time**: Sub-100ms for typical queries with proper indexing
- **Pagination**: Efficient LIMIT/OFFSET with no full table scans
- **Filtering**: O(1) indexed lookups for most filters
- **Sorting**: Database-level sorting for performance

## Testing

### Test Filtering
```bash
curl "http://localhost:8080/api/profiles?gender=male&country_id=NG"
```

### Test Sorting
```bash
curl "http://localhost:8080/api/profiles?sort_by=age&order=desc&limit=5"
```

### Test Pagination
```bash
curl "http://localhost:8080/api/profiles?page=2&limit=20"
```

### Test Natural Language Search
```bash
curl "http://localhost:8080/api/profiles/search?q=young+males+from+nigeria"
curl "http://localhost:8080/api/profiles/search?q=adult+females+from+kenya"
curl "http://localhost:8080/api/profiles/search?q=teenagers+above+17"
```

### Test CORS
```bash
curl -H "Origin: http://example.com" http://localhost:8080/api/profiles
```

## CORS Headers

All endpoints include CORS headers:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

## Timestamps

All timestamps are in UTC ISO 8601 format: `2026-04-01T12:00:00Z`

## UUID Format

All IDs use UUID v7 format for temporal sorting and uniqueness.

## Response Format

All responses follow a standard format with consistent structure for easy client-side handling.

## Contributing

To add new features or fix bugs:

1. Create a feature branch
2. Make changes with tests
3. Submit a pull request

## License

Proprietary - Insighta Labs
