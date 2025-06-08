# Recipe Smith

![Recipe Smith Banner](https://github.com/n0tt1m/code-assistant/raw/main/images/the-girls3.png)

### A comprehensive recipe collection and search system to help you discover what to cook for dinner tonight.

Recipe Smith is a full-stack application that crawls popular recipe websites, extracts structured recipe data, and provides a searchable interface to help you find the perfect recipe. The system consists of three main components: Pantry (web crawler), Sous (web API), and Sauté (Flutter mobile app).

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Pantry      │────│ Elasticsearch   │────│      Sous       │
│ (Web Crawler)   │    │   (Database)    │    │   (Web API)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                               ┌─────────────────┐
                                               │     Sauté       │
                                               │ (Flutter App)   │
                                               └─────────────────┘
```

## 🚀 Features

- **Intelligent Web Crawler**: Automatically crawls 16+ popular recipe websites
- **Structured Data Extraction**: Supports JSON-LD schema and site-specific parsing
- **Elasticsearch Integration**: Fast, scalable recipe search and storage
- **REST API**: Clean API for recipe retrieval and search
- **Mobile App**: Flutter-based mobile interface
- **Rate Limiting**: Respectful crawling with domain-specific limits
- **Recipe Backup**: JSON file backups for data redundancy

## 📱 Components

### Pantry (Web Crawler)
- **Language**: Go
- **Purpose**: Crawls recipe websites and extracts structured data
- **Features**: Multi-site support, JSON-LD parsing, rate limiting, backup system

### Sous (Web API)
- **Language**: Go
- **Purpose**: Provides REST API for recipe search and retrieval
- **Database**: Elasticsearch for fast full-text search

### Sauté (Mobile App)
- **Language**: Dart/Flutter
- **Purpose**: Cross-platform mobile app for recipe browsing
- **Features**: Recipe search, favorites, shopping lists

## 🛠️ Installation

### Prerequisites
- Go 1.19+
- Elasticsearch 7.x or 8.x
- Flutter 3.x (for mobile app)
- Docker (optional)

### Quick Start with Docker
```bash
# Clone the repository
git clone https://github.com/your-username/recipe-smith.git
cd recipe-smith

# Start Elasticsearch
docker-compose up -d

# Build and run the crawler
cd pantry
go build -o recipe-crawler *.go
./recipe-crawler recipes

# Start the API server
cd ../sous
go run *.go

# Run the mobile app
cd ../sauté
flutter run
```

### Manual Installation

#### 1. Setup Elasticsearch
```bash
# Using Docker
docker run -d --name elasticsearch \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  elasticsearch:8.11.0

# Or install locally following Elasticsearch documentation
```

#### 2. Build Pantry (Crawler)
```bash
cd pantry
go mod tidy
go build -o recipe-crawler *.go
```

#### 3. Build Sous (API Server)
```bash
cd sous
go mod tidy
go build -o api-server *.go
```

#### 4. Setup Sauté (Mobile App)
```bash
cd sauté
flutter pub get
flutter run
```

## 🎯 Usage

### Crawling Recipes

#### Crawl Popular Recipe Sites
```bash
cd pantry
./recipe-crawler recipes
```

#### Crawl Specific Site
```bash
./recipe-crawler index https://pinchofyum.com/recipes
```

#### Custom Crawling Parameters
```bash
./recipe-crawler index https://example.com/recipes \
  -workers=20 -depth=5 -delay=2 -debug=true
```

#### Test URL Extraction
```bash
./recipe-crawler test-url https://pinchofyum.com/easy-chicken-pad-thai
```

### API Server
```bash
cd sous
./api-server
# Server starts on http://localhost:8080
```

### Mobile App
```bash
cd sauté
flutter run
```

## 🌐 Supported Recipe Sites

The crawler supports 16+ popular recipe websites:

- **Food Blogs**: Pinch of Yum, Minimalist Baker, Cookie and Kate, Love and Lemons
- **Professional Sites**: Serious Eats, Food Network, Epicurious, AllRecipes
- **Specialty Sites**: Half Baked Harvest, Budget Bytes, 101 Cookbooks, Food52
- **International**: The Woks of Life (Chinese cuisine)
- **General**: Delish, Simply Recipes, Smitten Kitchen

## 🔧 Configuration

### Crawler Settings
```go
// pantry/main.go
maxCrawlDepth        = 3          // Maximum crawl depth
concurrentWorkers    = 10         // Number of concurrent workers
crawlDelayPerDomain  = 1 * time.Second  // Delay between requests
maxRequestsPerDomain = 5          // Max concurrent requests per domain
```

### API Configuration
```go
// sous/main.go
port = ":8080"  // API server port
```

## 📊 API Endpoints

### Recipe Search
```http
GET /search?q=chicken&limit=10&offset=0
```

### Get Recipe by ID
```http
GET /recipe/:id
```

### Get All Recipes
```http
GET /recipes?limit=20&offset=0
```

### Recipe Statistics
```http
GET /stats
```

## 🗂️ Data Structure

### Recipe Schema
```json
{
  "id": "unique-id",
  "title": "Recipe Title",
  "name": "Recipe Name",
  "description": "Recipe description",
  "url": "source-url",
  "image": "image-url",
  "prep_time": "15 minutes",
  "cook_time": "30 minutes",
  "total_time": "45 minutes",
  "servings": "4",
  "calories": "320",
  "ingredients": "ingredient1;ingredient2;ingredient3",
  "instructions": "step1;step2;step3",
  "source_site": "pinchofyum.com",
  "crawl_date": "2024-01-01T00:00:00Z"
}
```

## 🔒 Rate Limiting & Ethics

The crawler implements several politeness features:
- **Domain-specific rate limiting**: Maximum 5 concurrent requests per domain
- **Request delays**: 1-second delay between requests to same domain
- **Respectful User-Agent**: Identifies as a recipe collection bot
- **Robots.txt compliance**: (Recommended to implement)
- **Content respect**: Only extracts publicly available recipe data

## 🧪 Testing

### Test Crawler
```bash
cd pantry
go test ./...
```

### Test API
```bash
cd sous
go test ./...
```

### Test Mobile App
```bash
cd sauté
flutter test
```

## 📝 Development

### Adding New Recipe Sites

1. Add site configuration to `pantry/src/scraper/scraper.go`:
```go
{
    Domain: "newsite.com",
    URLPatterns: []string{"/recipe/", "/recipes/"},
    Selectors: SiteSelectors{
        RecipeTitle:       "h1.recipe-title",
        RecipeDescription: ".recipe-description",
        RecipeIngredients: ".ingredients li",
        RecipeInstructions: ".instructions li",
        RecipeTime:        ".recipe-time",
        RecipeServings:    ".servings",
        RecipeLinks:       []string{"a[href*='/recipe/']"},
    },
}
```

2. Add the site URL to `getPopularRecipeSites()` in `main.go`

### Project Structure
```
recipe-smith/
├── pantry/           # Web crawler (Go)
│   ├── main.go
│   ├── src/
│   │   ├── scraper/
│   │   ├── elasticsearch/
│   │   ├── logger/
│   │   └── structs/
│   └── recipe_backups/
├── sous/             # API server (Go)
│   ├── main.go
│   └── src/
├── sauté/            # Mobile app (Flutter)
│   ├── lib/
│   ├── android/
│   └── ios/
└── docker-compose.yml
```

## 🚧 Roadmap

- [ ] Implement robots.txt compliance
- [ ] Add recipe recommendation engine
- [ ] Support for video recipes
- [ ] Nutritional analysis integration
- [ ] User accounts and favorites
- [ ] Meal planning features
- [ ] Shopping list integration
- [ ] Recipe scaling calculator

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is for educational and personal use. Please respect the terms of service of the websites you crawl and ensure you have permission to scrape their content. The authors are not responsible for any misuse of this software.

## 🙏 Acknowledgments

- Thanks to all the amazing food bloggers who share their recipes
- Elasticsearch for providing excellent search capabilities
- The Go and Flutter communities for excellent documentation
- All contributors who help improve this project