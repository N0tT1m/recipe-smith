# Sous - Recipe API Server

Sous is a RESTful API server that provides access to the recipe database collected by the Pantry crawler. Built with Go and Elasticsearch, it offers fast search capabilities and comprehensive recipe data management.

## 🚀 Features

- **Fast Recipe Search**: Full-text search powered by Elasticsearch
- **RESTful API**: Clean, well-documented endpoints
- **Recipe Management**: CRUD operations for recipe data
- **Advanced Filtering**: Search by ingredients, cuisine, cooking time, etc.
- **Statistics API**: Recipe collection analytics
- **JSON Response**: Structured data in JSON format
- **CORS Support**: Cross-origin requests for web applications
- **Error Handling**: Comprehensive error responses

## 🛠️ Installation

### Prerequisites
- Go 1.19+
- Elasticsearch 7.x or 8.x
- Recipe data (crawled by Pantry)

### Setup
```bash
# Clone the repository
git clone https://github.com/your-username/recipe-smith.git
cd recipe-smith/sous

# Install dependencies
go mod tidy

# Build the server
go build -o api-server *.go

# Run the server
./api-server
```

The server will start on `http://localhost:8080` by default.

## 📊 API Endpoints

### Recipe Search

#### Search Recipes
```http
GET /search?q={query}&limit={limit}&offset={offset}
```

**Parameters:**
- `q` (string): Search query (recipe name, ingredients, etc.)
- `limit` (int): Number of results to return (default: 20, max: 100)
- `offset` (int): Number of results to skip (default: 0)

**Example:**
```bash
curl "http://localhost:8080/search?q=chicken%20pasta&limit=10&offset=0"
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "total_hits": 150,
    "pages": [
      {
        "id": "recipe-123",
        "title": "Creamy Chicken Pasta",
        "name": "Chicken Alfredo Pasta",
        "description": "A rich and creamy pasta dish...",
        "url": "https://example.com/chicken-pasta",
        "image": "https://example.com/image.jpg",
        "prep_time": "15 minutes",
        "cook_time": "20 minutes",
        "total_time": "35 minutes",
        "servings": "4",
        "calories": "450",
        "ingredients": "chicken breast;pasta;heavy cream;parmesan",
        "instructions": "Cook pasta;Season chicken;Make sauce;Combine",
        "source_site": "example.com",
        "crawl_date": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

### Recipe Management

#### Get Recipe by ID
```http
GET /recipe/{id}
```

**Example:**
```bash
curl "http://localhost:8080/recipe/recipe-123"
```

#### Get All Recipes
```http
GET /recipes?limit={limit}&offset={offset}
```

**Parameters:**
- `limit` (int): Number of results to return (default: 20)
- `offset` (int): Number of results to skip (default: 0)

**Example:**
```bash
curl "http://localhost:8080/recipes?limit=50&offset=100"
```

### Advanced Search

#### Search by Ingredients
```http
GET /search/ingredients?ingredients={ingredient1,ingredient2}&limit={limit}
```

**Example:**
```bash
curl "http://localhost:8080/search/ingredients?ingredients=chicken,garlic,onion&limit=10"
```

#### Search by Cuisine Type
```http
GET /search/cuisine?type={cuisine}&limit={limit}
```

**Example:**
```bash
curl "http://localhost:8080/search/cuisine?type=italian&limit=10"
```

#### Search by Cooking Time
```http
GET /search/time?max_time={minutes}&limit={limit}
```

**Example:**
```bash
curl "http://localhost:8080/search/time?max_time=30&limit=10"
```

### Statistics

#### Get Recipe Statistics
```http
GET /stats
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "total_recipes": 15420,
    "recipes_by_site": {
      "pinchofyum.com": 1250,
      "minimalistbaker.com": 980,
      "seriouseats.com": 1100
    },
    "avg_ingredients_count": 8.5,
    "avg_instructions_count": 6.2,
    "most_common_categories": [
      {"category": "dinner", "count": 5200},
      {"category": "vegetarian", "count": 3100}
    ],
    "last_crawled": "2024-01-01T15:30:00Z"
  }
}
```

#### Get Site Statistics
```http
GET /stats/sites
```

### Health Check

#### Server Health
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "elasticsearch": "connected",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 🔧 Configuration

### Environment Variables
```bash
export PORT=8080                    # Server port
export ELASTICSEARCH_URL=http://localhost:9200  # Elasticsearch URL
export LOG_LEVEL=info               # Log level (debug, info, warn, error)
export CORS_ENABLED=true            # Enable CORS
export MAX_RESULTS_PER_PAGE=100     # Maximum results per page
```

### Default Configuration
```go
// Default settings in main.go
port := ":8080"
elasticsearchURL := "http://localhost:9200"
maxResultsPerPage := 100
requestTimeout := 30 * time.Second
```

## 🗂️ Data Structures

### Recipe Response Format
```go
type Page struct {
    ID           string    `json:"id"`
    Title        string    `json:"title"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    URL          string    `json:"url"`
    Image        string    `json:"image"`
    PrepTime     string    `json:"prep_time"`
    CookTime     string    `json:"cook_time"`
    TotalTime    string    `json:"total_time"`
    Servings     string    `json:"servings"`
    Calories     string    `json:"calories"`
    Ingredients  string    `json:"ingredients"`  // Semicolon-separated
    Instructions string    `json:"instructions"` // Semicolon-separated
    SourceSite   string    `json:"source_site"`
    CrawlDate    time.Time `json:"crawl_date"`
}
```

### Search Response Format
```go
type SearchResult struct {
    TotalHits int64  `json:"total_hits"`
    Pages     []Page `json:"pages"`
}

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}
```

## 🔍 Search Features

### Full-Text Search
- Searches across recipe titles, descriptions, and ingredients
- Supports partial matching and fuzzy search
- Elasticsearch-powered relevance scoring

### Advanced Filters
- **Ingredients**: Find recipes containing specific ingredients
- **Cooking Time**: Filter by maximum cooking time
- **Servings**: Filter by number of servings
- **Source Site**: Filter by recipe source website
- **Date Range**: Filter by crawl date

### Search Examples

#### Find Quick Recipes
```bash
curl "http://localhost:8080/search/time?max_time=15&q=easy"
```

#### Vegetarian Pasta Recipes
```bash
curl "http://localhost:8080/search?q=pasta%20vegetarian&limit=20"
```

#### Recipes with Specific Ingredients
```bash
curl "http://localhost:8080/search/ingredients?ingredients=chicken,broccoli,cheese"
```

## 🚀 Performance

### Optimization Features
- **Elasticsearch Indexing**: Fast full-text search
- **Connection Pooling**: Efficient database connections
- **Response Caching**: Configurable cache headers
- **Pagination**: Efficient large dataset handling
- **Query Optimization**: Optimized Elasticsearch queries

### Performance Tips
- Use pagination for large result sets
- Implement client-side caching for frequently accessed data
- Use specific search terms for better performance
- Consider implementing rate limiting for production use

## 🧪 Testing

### Manual Testing
```bash
# Start the server
./api-server

# Test search endpoint
curl "http://localhost:8080/search?q=chicken"

# Test recipe retrieval
curl "http://localhost:8080/recipes?limit=5"

# Test statistics
curl "http://localhost:8080/stats"
```

### Automated Testing
```bash
# Run unit tests
go test ./...

# Run integration tests (requires Elasticsearch)
go test -tags=integration ./...
```

## 🐛 Error Handling

### HTTP Status Codes
- `200 OK`: Successful request
- `400 Bad Request`: Invalid parameters
- `404 Not Found`: Recipe not found
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Elasticsearch unavailable

### Error Response Format
```json
{
  "status": "error",
  "message": "Recipe not found",
  "data": null
}
```

### Common Errors

#### Elasticsearch Connection Error
```json
{
  "status": "error",
  "message": "Database connection failed",
  "data": null
}
```

#### Invalid Search Parameters
```json
{
  "status": "error",
  "message": "Invalid limit parameter: must be between 1 and 100",
  "data": null
}
```

## 🔒 Security Considerations

### Current Implementation
- Input validation and sanitization
- Error message sanitization
- CORS configuration for web applications

### Recommended Enhancements
- **Rate Limiting**: Implement request rate limiting
- **API Authentication**: Add API key or JWT authentication
- **HTTPS**: Use HTTPS in production
- **Input Validation**: Enhanced parameter validation
- **Logging**: Security event logging

## 🚀 Deployment

### Local Development
```bash
# Start Elasticsearch
docker run -d -p 9200:9200 -e "discovery.type=single-node" elasticsearch:8.11.0

# Start the API server
go run *.go
```

### Production Deployment
```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server *.go

# Run with production settings
export PORT=8080
export ELASTICSEARCH_URL=http://your-elasticsearch:9200
./api-server
```

### Docker Deployment
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-server *.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-server .
EXPOSE 8080
CMD ["./api-server"]
```

## 📝 Development

### Project Structure
```
sous/
├── main.go                 # Main server application
├── api-server             # Compiled binary
├── src/
│   ├── elasticsearch/
│   │   └── elastic_search.go # Elasticsearch client
│   ├── logger/
│   │   └── logger.go      # Logging utilities
│   └── structs/
│       └── structs.go     # Data structures
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

### Adding New Endpoints

1. **Define Handler Function**
```go
func newEndpointHandler(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

2. **Add Route**
```go
http.HandleFunc("/new-endpoint", newEndpointHandler)
```

3. **Add Documentation**
Update this README with the new endpoint details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests for new functionality
5. Update documentation
6. Submit a pull request

### Development Guidelines
- Follow Go best practices
- Add comprehensive error handling
- Include unit tests
- Update API documentation
- Use consistent code formatting

## 📋 Monitoring

### Logging
- Request/response logging
- Error logging with stack traces
- Performance metrics logging
- Elasticsearch query logging

### Health Checks
- `/health` endpoint for monitoring
- Elasticsearch connectivity checks
- Response time monitoring
- Error rate tracking

## 🎯 Future Enhancements

- [ ] Recipe recommendation engine
- [ ] User favorites and ratings
- [ ] Recipe nutrition analysis
- [ ] Advanced search filters
- [ ] Recipe similarity matching
- [ ] Meal planning endpoints
- [ ] Shopping list generation
- [ ] Recipe scaling calculator