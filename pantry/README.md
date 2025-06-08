# Pantry - Recipe Web Crawler

Pantry is an intelligent web crawler designed to extract structured recipe data from popular cooking websites. It supports multiple crawling strategies, respects rate limits, and provides comprehensive recipe data extraction.

## 🚀 Features

- **Multi-Site Support**: Crawls 16+ popular recipe websites with site-specific configurations
- **JSON-LD Support**: Extracts structured recipe data using schema.org standards
- **Rate Limiting**: Domain-specific request limiting and delays for respectful crawling
- **Backup System**: Automatically saves recipe data to JSON files for redundancy
- **Elasticsearch Integration**: Stores and indexes recipes for fast searching
- **Configurable Parameters**: Customizable crawl depth, workers, and delays
- **Debug Mode**: Detailed logging for troubleshooting and development

## 🌐 Supported Recipe Sites

### Food Blogs
- **Pinch of Yum** (pinchofyum.com)
- **Minimalist Baker** (minimalistbaker.com)
- **Cookie and Kate** (cookieandkate.com)
- **Love and Lemons** (loveandlemons.com)
- **Smitten Kitchen** (smittenkitchen.com)
- **Half Baked Harvest** (halfbakedharvest.com)
- **Budget Bytes** (budgetbytes.com)

### Professional Sites
- **Serious Eats** (seriouseats.com)
- **Food Network** (foodnetwork.com)
- **Epicurious** (epicurious.com)
- **AllRecipes** (allrecipes.com)
- **Simply Recipes** (simplyrecipes.com)
- **Delish** (delish.com)

### Specialty Sites
- **101 Cookbooks** (101cookbooks.com)
- **Food52** (food52.com)
- **The Woks of Life** (thewoksoflife.com)

## 🛠️ Installation

### Prerequisites
- Go 1.19+
- Elasticsearch 7.x or 8.x running on localhost:9200

### Setup
```bash
# Clone the repository
git clone https://github.com/your-username/recipe-smith.git
cd recipe-smith/pantry

# Install dependencies
go mod tidy

# Build the crawler
go build -o recipe-crawler *.go
```

### Configuration
The crawler connects to Elasticsearch at `localhost:9200` by default. Ensure Elasticsearch is running before starting the crawler.

## 🎯 Usage

### Basic Commands

#### Crawl All Popular Recipe Sites
```bash
./recipe-crawler recipes
```

#### Crawl Specific Website
```bash
./recipe-crawler index https://pinchofyum.com/recipes
```

#### Test Recipe Extraction
```bash
./recipe-crawler test-url https://pinchofyum.com/easy-chicken-pad-thai
```

#### Delete Recipe Index
```bash
./recipe-crawler delete
```

### Advanced Usage

#### Custom Crawling Parameters
```bash
./recipe-crawler index https://example.com/recipes \
  -workers=20 \      # Number of concurrent workers
  -depth=5 \         # Maximum crawl depth
  -delay=2 \         # Delay between requests (seconds)
  -debug=true        # Enable debug logging
```

#### Available Parameters
- `-workers=N`: Number of concurrent workers (default: 10)
- `-depth=N`: Maximum crawl depth (default: 3)
- `-delay=N`: Delay between requests in seconds (default: 1)
- `-max-requests=N`: Max concurrent requests per domain (default: 5)
- `-debug=true/false`: Enable debug mode (default: false)

## 🔧 Configuration

### Default Settings
```go
maxCrawlDepth        = 3                    // Maximum crawl depth
concurrentWorkers    = 10                   // Number of concurrent workers
crawlDelayPerDomain  = 1 * time.Second      // Delay between requests
maxRequestsPerDomain = 5                    // Max concurrent requests per domain
```

### Rate Limiting
The crawler implements several politeness features:
- **Domain Semaphores**: Limits concurrent requests per domain
- **Request Delays**: Configurable delays between requests
- **Timeout Protection**: 30-minute maximum crawl time
- **Respectful Headers**: Proper User-Agent identification

## 📊 Data Extraction

### Supported Data Fields
- **Basic Info**: Title, name, description, URL, source site
- **Timing**: Prep time, cook time, total time
- **Serving Info**: Servings/yield, calories
- **Recipe Data**: Ingredients list, step-by-step instructions
- **Media**: Recipe images
- **Metadata**: Crawl date, categories

### Extraction Methods
1. **JSON-LD Schema**: Primary method for modern recipe sites
2. **Site-Specific Selectors**: Custom CSS selectors for each supported site
3. **Fallback Parsing**: Generic HTML parsing for unknown sites

### Data Format
```json
{
  "id": "unique-recipe-id",
  "title": "Easy Chicken Pad Thai",
  "name": "Chicken Pad Thai",
  "description": "A quick and easy pad thai recipe...",
  "url": "https://pinchofyum.com/easy-chicken-pad-thai",
  "image": "https://example.com/image.jpg",
  "prep_time": "15 minutes",
  "cook_time": "10 minutes",
  "total_time": "25 minutes",
  "servings": "4",
  "calories": "380",
  "ingredients": "chicken;rice noodles;eggs;bean sprouts;...",
  "instructions": "Heat oil in pan;Add chicken;Scramble eggs;...",
  "source_site": "pinchofyum.com",
  "crawl_date": "2024-01-01T12:00:00Z"
}
```

## 🗂️ File Structure

```
pantry/
├── main.go                 # Main crawler application
├── enhanced-recipe-crawler # Compiled binary
├── recipe_crawler.log     # Crawler logs
├── recipe_backups/        # JSON backup files
├── src/
│   ├── scraper/
│   │   └── scraper.go     # Web scraping logic
│   ├── elasticsearch/
│   │   └── elastic_search.go # Elasticsearch integration
│   ├── logger/
│   │   └── logger.go      # Logging utilities
│   ├── structs/
│   │   └── structs.go     # Data structures
│   └── variables/
│       └── variables.go   # Configuration variables
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## 🧪 Testing

### Test Recipe Extraction
```bash
# Test specific recipe URL
./recipe-crawler test-url https://minimalistbaker.com/simple-almond-butter-cookies/

# Expected output shows:
# - Extracted recipe data
# - URL analysis
# - Data requirements check
# - Found links
```

### Debug Mode
Enable debug mode for detailed logging:
```bash
./recipe-crawler index https://example.com/recipes -debug=true
```

Debug mode provides:
- Detailed extraction logs
- Site configuration matching
- Link analysis results
- Error diagnostics

## 🔍 Adding New Recipe Sites

### 1. Add Site Configuration
Edit `src/scraper/scraper.go` and add a new site to the `recipeSites` array:

```go
{
    Domain: "newsite.com",
    URLPatterns: []string{
        "/recipe/",
        "/recipes/",
    },
    Selectors: SiteSelectors{
        RecipeTitle:       "h1.recipe-title, .recipe-name",
        RecipeDescription: ".recipe-description, .recipe-summary",
        RecipeIngredients: ".ingredients li, .recipe-ingredients li",
        RecipeInstructions: ".instructions li, .directions li",
        RecipeTime:        ".prep-time, .cook-time, .total-time",
        RecipeServings:    ".servings, .yield",
        RecipeLinks:       []string{"a[href*='/recipe/']"},
    },
},
```

### 2. Add Starting URLs
Add the site's recipe index URL to `getPopularRecipeSites()` in `main.go`:

```go
func getPopularRecipeSites() []string {
    return []string{
        // ... existing sites ...
        "https://newsite.com/recipes",
    }
}
```

### 3. Test the New Site
```bash
./recipe-crawler test-url https://newsite.com/some-recipe
```

## 📈 Performance

### Optimization Features
- **Concurrent Processing**: Configurable worker pools
- **Domain Rate Limiting**: Prevents overwhelming servers
- **Efficient Parsing**: Targets specific recipe data
- **Memory Management**: Streaming JSON processing
- **Connection Pooling**: Reuses HTTP connections

### Monitoring
- **Real-time Logs**: Track crawling progress
- **Error Reporting**: Detailed error logs
- **Statistics**: Crawled URL counts and success rates
- **Backup System**: JSON files for data recovery

## ⚠️ Important Notes

### Ethical Crawling
- Respects rate limits and implements delays
- Only extracts publicly available recipe data
- Uses appropriate User-Agent headers
- Recommends implementing robots.txt compliance

### Troubleshooting
- **Connection Issues**: Ensure Elasticsearch is running
- **Slow Crawling**: Reduce workers or increase delays
- **Memory Issues**: Lower concurrent workers
- **Rate Limiting**: Increase delays between requests

### Legal Considerations
- Review terms of service for each website
- Ensure compliance with copyright laws
- Use extracted data responsibly
- Consider reaching out to site owners for permission

## 📋 Logs and Monitoring

### Log Files
- `recipe_crawler.log`: Main application logs
- `pantry.log`: Legacy log file
- Console output: Real-time crawling status

### Backup System
- `recipe_backups/`: JSON files for each extracted recipe
- Automatic backup creation for data redundancy
- File naming: `{recipe-id}_{safe-title}.json`

## 🐛 Troubleshooting

### Common Issues

**Elasticsearch Connection Failed**
```bash
# Check if Elasticsearch is running
curl http://localhost:9200

# Start Elasticsearch with Docker
docker run -d -p 9200:9200 -e "discovery.type=single-node" elasticsearch:8.11.0
```

**No Recipes Extracted**
- Enable debug mode: `-debug=true`
- Test specific URLs with `test-url` command
- Check site-specific selectors in scraper.go

**Rate Limiting Issues**
- Increase delay: `-delay=2`
- Reduce workers: `-workers=5`
- Check server response codes in logs

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure existing tests pass
5. Submit a pull request

### Development Setup
```bash
# Install development dependencies
go mod tidy

# Run tests
go test ./...

# Build and test
go build -o recipe-crawler *.go
./recipe-crawler test-url https://example.com/recipe
```