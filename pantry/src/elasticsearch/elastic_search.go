package elasticsearch

// Elastic search client
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"reflect"
	"regexp"
	"search-engine-indexer/src/logger"
	"search-engine-indexer/src/structs"
	"strings"
	"time"

	elastic "github.com/olivere/elastic/v7"
)

// Update the IndexMapping in elasticsearch/elasticsearch.go
const (
	IndexName    = "recipes"
	IndexMapping = `{
        "settings":{
            "number_of_shards":1,
            "number_of_replicas":0,
            "analysis": {
                "analyzer": {
                    "recipe_analyzer": {
                        "type": "custom",
                        "tokenizer": "standard",
                        "char_filter": ["html_strip"],
                        "filter": ["lowercase", "asciifolding", "stop", "snowball"]
                    }
                }
            }
        },
        "mappings":{
            "properties":{
                "title": {
                    "type": "text",
                    "analyzer": "recipe_analyzer",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 256
                        }
                    }
                },
                "description": {
                    "type": "text",
                    "analyzer": "recipe_analyzer"
                },
                "body": {
                    "type": "text",
                    "analyzer": "recipe_analyzer"
                },
                "url": {
                    "type": "text",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 2048
                        }
                    }
                },
                "image": {
                    "type": "text",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 2048
                        }
                    }
                },
                "name": {
                    "type": "text",
                    "analyzer": "recipe_analyzer",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 256
                        }
                    }
                },
                "prep_time": {
                    "type": "text"
                },
                "cook_time": {
                    "type": "text"
                },
                "total_time": {
                    "type": "text"
                },
                "calories": {
                    "type": "text"
                },
                "servings": {
                    "type": "text"
                },
                "ingredients": {
                    "type": "text",
                    "analyzer": "recipe_analyzer"
                },
                "instructions": {
                    "type": "text",
                    "analyzer": "recipe_analyzer"
                },
                "source_site": {
                    "type": "keyword"
                },
                "crawl_date": {
                    "type": "date"
                }
            }
        }
    }`
)

var client *elastic.Client

// NewElasticSearchClient returns an elastic seach client
func NewElasticSearchClient() *elastic.Client {
	var err error
	connected := false
	retries := 0

	// Custom retry strategy for docker-compose initialization
	for connected == false {
		// Create a new elastic client
		client, err = elastic.NewClient(
			elastic.SetURL("http://192.168.1.78:9200"), elastic.SetSniff(false))
		if err != nil {
			// log.Fatal(err)
			if retries == 5 {
				log.Fatal(err)
			}
			fmt.Println("Elasticsearch isn't ready for connection", 5-retries, "less")
			retries++
			time.Sleep(3 * time.Second)
		} else {
			connected = true
		}
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://192.168.1.78:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	return client
}

// ExistsIndex checks if the given index exists or not
func ExistsIndex(i string) bool {
	// Check if index exists
	exists, err := client.IndexExists(i).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return exists
}

// CreateIndex creates a new index
func CreateIndex(i string) {
	createIndex, err := client.CreateIndex(IndexName).
		Body(IndexMapping).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if !createIndex.Acknowledged {
		log.Println("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}
}

// DeleteIndex in the indexName constant
func DeleteIndex() {
	ctx := context.Background()
	deleteIndex, err := client.DeleteIndex(IndexName).Do(ctx)
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
	if !deleteIndex.Acknowledged {
		log.Println("DeleteIndex was not acknowledged. Check that timeout value is correct.")
	}
	fmt.Println("Index", IndexName, "deleted")
}

// ExistingPage return a boolean and a page if the title is already
// stored in the database
func ExistingPage(title string) (bool, structs.Page) {
	var exists bool
	var p structs.Page

	ctx := context.Background()

	// Normalize the title before searching
	normalizedTitle := normalizeTitle(title)

	logger.WriteInfo(normalizedTitle)

	// Use BoolQuery with MatchPhraseQuery for exact matching
	q := elastic.NewBoolQuery().
		Must(elastic.NewMatchPhraseQuery("title", normalizedTitle))

	result, err := client.Search().
		Index(IndexName).
		Query(q).
		Size(1).
		Do(ctx)

	if err != nil {
		logger.WriteError("Failed to search for page:", err)
		return false, p
	}

	var ttyp structs.Page
	for _, item := range result.Each(reflect.TypeOf(ttyp)) {
		page := item.(structs.Page)
		// Compare normalized titles
		if normalizeTitle(page.Title) == normalizedTitle {
			exists, p = true, page
			return exists, p
		}
	}

	return exists, p
}

// normalizeTitle removes extra whitespace and converts to lowercase for consistent matching
func normalizeTitle(title string) string {
	// Remove extra whitespace and convert to lowercase
	normalized := strings.TrimSpace(strings.ToLower(title))

	// Split into words and rejoin to ensure consistent spacing
	words := strings.Fields(normalized)
	return strings.Join(words, " ")
}

// Replace isValidDelishURL with a more generic function
func isValidRecipeURL(rawURL string) bool {
	// List of allowed recipe domains
	allowedDomains := []string{
		"www.delish.com",
		"www.allrecipes.com",
		"www.foodnetwork.com",
		"www.epicurious.com",
		"www.simplyrecipes.com",
		// Add more sites as needed
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Check if domain is in allowed list
	domainValid := false
	for _, domain := range allowedDomains {
		if parsedURL.Host == domain {
			domainValid = true
			break
		}
	}

	if !domainValid {
		return false
	}

	// Map of domain-specific validation patterns
	validPathPatterns := map[string][]string{
		"www.delish.com":        {"/cooking/recipe-ideas/", "/everyday-cooking/quick-and-easy/"},
		"www.allrecipes.com":    {"/recipe/", "/recipes/"},
		"www.foodnetwork.com":   {"/recipes/", "/recipe/"},
		"www.epicurious.com":    {"/recipes/", "/recipe/"},
		"www.simplyrecipes.com": {"/recipes/"},
		// Add more patterns as needed
	}

	// Check domain-specific path patterns
	patterns, exists := validPathPatterns[parsedURL.Host]
	if !exists {
		return false
	}

	path := parsedURL.Path
	for _, pattern := range patterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	return false
}

// Add this function to elasticsearch/elasticsearch.go
func normalizeRecipeTitle(title string) string {
	// Convert to lowercase
	normalized := strings.ToLower(title)

	// Remove common recipe title patterns
	normalized = strings.ReplaceAll(normalized, "recipe", "")
	normalized = strings.ReplaceAll(normalized, "best", "")
	normalized = strings.ReplaceAll(normalized, "easy", "")
	normalized = strings.ReplaceAll(normalized, "homemade", "")

	// Remove special characters and extra spaces
	normalized = regexp.MustCompile(`[^\w\s]`).ReplaceAllString(normalized, " ")
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")

	return strings.TrimSpace(normalized)
}

// Improve the existingTitle function
func existingTitle(client *elastic.Client, indexName string, title string) bool {
	ctx := context.Background()

	// Normalize the title for comparison
	normalizedTitle := normalizeRecipeTitle(title)

	// Use fuzzy matching for better duplicate detection
	q := elastic.NewBoolQuery().
		Should(
			elastic.NewMatchPhraseQuery("title", title).Boost(2),
			elastic.NewMatchQuery("title", normalizedTitle).Fuzziness("2"),
		).
		MinimumShouldMatch("1")

	// Execute the search
	result, err := client.Search().
		Index(indexName).
		Query(q).
		Size(5).
		Do(ctx)

	if err != nil {
		logger.WriteWarning(fmt.Sprintf("Failed to check for existing title: %v", err))
		return false
	}

	// Check if any results exceed our similarity threshold
	if result.TotalHits() > 0 {
		// Log the potential duplicates
		for _, hit := range result.Hits.Hits {
			var page structs.Page
			if err := json.Unmarshal(hit.Source, &page); err == nil {
				similarity := calculateTitleSimilarity(title, page.Title)
				if similarity > 0.7 { // Threshold for considering it a duplicate
					logger.WriteInfo(fmt.Sprintf("Found duplicate title (%.2f similarity): %s", similarity, page.Title))
					return true
				}
			}
		}
	}

	return false
}

// Add a function to calculate title similarity
func calculateTitleSimilarity(title1, title2 string) float64 {
	// Normalize both titles
	norm1 := normalizeRecipeTitle(title1)
	norm2 := normalizeRecipeTitle(title2)

	// Simple Levenshtein distance calculation
	// (In a real implementation, you might want to use a proper string similarity library)
	distance := levenshteinDistance(norm1, norm2)
	maxLength := math.Max(float64(len(norm1)), float64(len(norm2)))

	if maxLength == 0 {
		return 1.0
	}

	return 1.0 - float64(distance)/maxLength
}

// Levenshtein distance implementation
func levenshteinDistance(s1, s2 string) int {
	// Implementation of the algorithm (simplified version)
	// In a real app, use a library for this

	// Simple implementation for example purposes
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create 2D slice for calculations
	d := make([][]int, len(s1)+1)
	for i := range d {
		d[i] = make([]int, len(s2)+1)
	}

	// Initialize base cases
	for i := 0; i <= len(s1); i++ {
		d[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		d[0][j] = j
	}

	// Calculate distances
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			d[i][j] = min(d[i-1][j]+1, min(d[i][j-1]+1, d[i-1][j-1]+cost))
		}
	}

	return d[len(s1)][len(s2)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// existingURL checks if a URL already exists in the database
func existingURL(client *elastic.Client, indexName string, urlToCheck string) bool {
	ctx := context.Background()

	// Create a bool query to check for exact URL match
	q := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("url.keyword", urlToCheck))

	// Execute the search
	result, err := client.Search().
		Index(indexName).
		Query(q).
		Size(25).
		Do(ctx)

	if err != nil {
		logger.WriteWarning(fmt.Sprintf("Failed to check for existing URL: %v", err))
		return false
	}

	// Log the search results for debugging
	hits := result.TotalHits()
	if hits > 0 {
		logger.WriteInfo(fmt.Sprintf("Found %d existing entries for URL: %s", hits, urlToCheck))
		// Log the first matching document
		var page structs.Page
		if err := json.Unmarshal(result.Hits.Hits[0].Source, &page); err == nil {
			logger.WriteInfo(fmt.Sprintf("Existing page details - ID: %s, Title: %s", page.ID, page.Title))
		}
		return true
	}

	return false
}

// Update the CreatePage function to handle recipe fields
func CreatePage(p structs.Page) bool {
	ctx := context.Background()

	// Convert relative URL to absolute if needed
	if strings.HasPrefix(p.URL, "/") {
		// Extract domain from context or use a default one
		domain := extractDomainFromURL(p.URL)
		if domain == "" {
			domain = "https://www.example.com" // Fallback
		}
		p.URL = domain + p.URL
	}

	// Validate URL format for any supported recipe site
	if !isValidRecipeURL(p.URL) {
		logger.WriteWarning(fmt.Sprintf("Invalid recipe URL format: %s", p.URL))
		return false
	}

	// Set source site from URL
	p.SourceSite = extractSourceSite(p.URL)

	// Set crawl date to current time
	p.CrawlDate = time.Now()

	// If name isn't set, use the title
	if p.Name == "" {
		p.Name = p.Title
	}

	// Check if URL already exists
	if existingURL(client, IndexName, p.URL) {
		logger.WriteWarning(fmt.Sprintf("URL already exists in database: %s", p.URL))
		return false
	}

	// Check if the title is a duplicate (using improved detection)
	if existingTitle(client, IndexName, p.Title) {
		logger.WriteWarning(fmt.Sprintf("Similar title already exists: %s", p.Title))
		return false
	}

	// Additional recipe-specific validation
	if len(p.Title) < 5 || len(p.Description) < 10 {
		logger.WriteWarning(fmt.Sprintf("Recipe missing key information: %s", p.Title))
		return false
	}

	// Create the new page with refresh to ensure immediate visibility
	_, err := client.Index().
		Index(IndexName).
		Id(p.ID).
		Refresh("true").
		BodyJson(p).
		Do(ctx)

	if err != nil {
		logger.WriteWarning(fmt.Sprintf("Failed to create the page: %v", err))
		return false
	}

	logger.WriteInfo(fmt.Sprintf("Successfully created new recipe - Title: %s, URL: %s", p.Title, p.URL))
	return true
}

// Helper function to extract source site from URL
func extractSourceSite(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}

// Helper function to extract domain from URL
func extractDomainFromURL(urlPath string) string {
	// Map of known domains based on URL path patterns
	domainPatterns := map[string]string{
		"/cooking/recipe-ideas/": "https://www.delish.com",
		"/recipes/":              "https://www.foodnetwork.com", // Default for common pattern
		"/recipe/":               "https://www.allrecipes.com",  // Default for common pattern
	}

	for pattern, domain := range domainPatterns {
		if strings.Contains(urlPath, pattern) {
			return domain
		}
	}

	return ""
}

// UpdatePage updates an existing page only if the URL is valid and unique
func UpdatePage(id string, params map[string]interface{}) bool {
	ctx := context.Background()

	// If URL is being updated, validate it
	if url, ok := params["url"].(string); ok {
		// Convert relative URL to absolute if needed
		if strings.HasPrefix(url, "/") {
			domain := extractDomainFromURL(url)
			if domain == "" {
				domain = "https://www.example.com" // Fallback
			}
			url = domain + url
			params["url"] = url
		}

		if !isValidRecipeURL(url) {
			logger.WriteWarning(fmt.Sprintf("Invalid URL format: %s", url))
			return false
		}

		// Set source site from URL
		params["source_site"] = extractSourceSite(url)

		// Get the current document to compare URLs
		currentDoc, err := client.Get().
			Index(IndexName).
			Id(id).
			Do(ctx)

		if err != nil {
			logger.WriteWarning(fmt.Sprintf("Failed to get current document: %v", err))
			return false
		}

		var currentPage structs.Page
		if err := json.Unmarshal(currentDoc.Source, &currentPage); err != nil {
			logger.WriteWarning(fmt.Sprintf("Failed to unmarshal current page: %v", err))
			return false
		}

		// Only check for existing URL if it's different from the current URL
		if currentPage.URL != url {
			logger.WriteInfo(fmt.Sprintf("Checking if new URL exists: %s", url))
			if existingURL(client, IndexName, url) {
				logger.WriteWarning(fmt.Sprintf("URL already exists in another document, skipping update: %s", url))
				return false
			}
		}
	}

	// If title is being updated, check for duplicates
	if title, ok := params["title"].(string); ok {
		// Get the current document to compare titles
		currentDoc, err := client.Get().
			Index(IndexName).
			Id(id).
			Do(ctx)

		if err != nil {
			logger.WriteWarning(fmt.Sprintf("Failed to get current document: %v", err))
			return false
		}

		var currentPage structs.Page
		if err := json.Unmarshal(currentDoc.Source, &currentPage); err != nil {
			logger.WriteWarning(fmt.Sprintf("Failed to unmarshal current page: %v", err))
			return false
		}

		// Only check for existing title if it's different from the current title
		if currentPage.Title != title {
			logger.WriteInfo(fmt.Sprintf("Checking if new title exists: %s", title))
			if existingTitle(client, IndexName, title) {
				logger.WriteWarning(fmt.Sprintf("Title already exists in another document, skipping update: %s", title))
				return false
			}
		}

		// If name isn't being updated, sync it with title
		if _, ok := params["name"]; !ok {
			params["name"] = title
		}
	}

	// Update crawl_date
	params["crawl_date"] = time.Now()

	// Perform the update with refresh to ensure immediate visibility
	_, err := client.Update().
		Index(IndexName).
		Id(id).
		Doc(params).
		Refresh("true").
		RetryOnConflict(3).
		Do(ctx)

	if err != nil {
		logger.WriteWarning(fmt.Sprintf("Failed to update the page: %v", err))
		return false
	}

	logger.WriteInfo(fmt.Sprintf("Successfully updated page with ID: %s", id))
	return true
}
