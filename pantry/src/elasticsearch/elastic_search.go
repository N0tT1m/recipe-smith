package elasticsearch

// Elastic search client
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
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
	maxRetries := 10
	retryDelay := 5 * time.Second

	// Custom retry strategy for docker-compose initialization
	for connected == false && retries < maxRetries {
		// Create a new elastic client
		client, err = elastic.NewClient(
			elastic.SetURL("http://192.168.1.78:9200"),
			elastic.SetSniff(false),
			elastic.SetHealthcheck(true),
			elastic.SetHealthcheckTimeout(20*time.Second),
			elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(100*time.Millisecond, 5*time.Second))),
		)
		if err != nil {
			logger.WriteWarning(fmt.Sprintf("Failed to connect to Elasticsearch (attempt %d/%d): %v", retries+1, maxRetries, err))
			retries++
			time.Sleep(retryDelay)
		} else {
			connected = true
		}
	}

	if !connected {
		log.Fatalf("Failed to connect to Elasticsearch after %d attempts", maxRetries)
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://192.168.1.78:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	logger.WriteInfo(fmt.Sprintf("Connected to Elasticsearch version %s", esversion))

	return client
}

// ExistsIndex checks if the given index exists or not
func ExistsIndex(i string) bool {
	// Check if index exists
	exists, err := client.IndexExists(i).Do(context.TODO())
	if err != nil {
		logger.WriteError(fmt.Sprintf("Error checking if index exists: %v", err))
		return false
	}

	return exists
}

// CreateIndex creates a new index
func CreateIndex(i string) {
	createIndex, err := client.CreateIndex(IndexName).
		Body(IndexMapping).
		Do(context.Background())
	if err != nil {
		logger.WriteError(fmt.Sprintf("Failed to create index: %v", err))
		return
	}

	if !createIndex.Acknowledged {
		logger.WriteWarning("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}

	logger.WriteInfo(fmt.Sprintf("Created index %s successfully", i))
}

// DeleteIndex in the indexName constant
func DeleteIndex() {
	ctx := context.Background()
	deleteIndex, err := client.DeleteIndex(IndexName).Do(ctx)
	if err != nil {
		// Handle error
		logger.WriteError(fmt.Sprintf("Failed to delete index: %v", err))
		return
	}
	if !deleteIndex.Acknowledged {
		logger.WriteWarning("DeleteIndex was not acknowledged. Check that timeout value is correct.")
	}
	logger.WriteInfo(fmt.Sprintf("Index %s deleted", IndexName))
}

// ExistingPage return a boolean and a page if the title is already
// stored in the database
func ExistingPage(title string) (bool, structs.Page) {
	var exists bool
	var p structs.Page

	ctx := context.Background()

	// Normalize the title before searching
	normalizedTitle := normalizeTitle(title)

	logger.WriteInfo(fmt.Sprintf("Searching for existing page with title: %s", normalizedTitle))

	// Use BoolQuery with MatchPhraseQuery for exact matching
	q := elastic.NewBoolQuery().
		Should(
			elastic.NewMatchPhraseQuery("title", normalizedTitle).Boost(2),
			elastic.NewMatchQuery("title", normalizedTitle).Fuzziness("2"),
		).
		MinimumShouldMatch("1")

	result, err := client.Search().
		Index(IndexName).
		Query(q).
		Size(5).
		Do(ctx)

	if err != nil {
		logger.WriteError(fmt.Sprintf("Failed to search for page: %v", err))
		return false, p
	}

	// Set a similarity threshold for considering a page a match
	similarityThreshold := 0.8

	var ttyp structs.Page
	for _, hit := range result.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &ttyp); err != nil {
			logger.WriteError(fmt.Sprintf("Failed to unmarshal page: %v", err))
			continue
		}

		// Calculate title similarity
		similarity := calculateTitleSimilarity(title, ttyp.Title)
		logger.WriteInfo(fmt.Sprintf("Title similarity: %.2f for '%s' vs '%s'", similarity, title, ttyp.Title))

		// If similarity is above threshold, consider it a match
		if similarity >= similarityThreshold {
			exists, p = true, ttyp
			return exists, p
		}
	}

	return exists, p
}

// normalizeTitle removes extra whitespace and converts to lowercase for consistent matching
func normalizeTitle(title string) string {
	// Remove extra whitespace and convert to lowercase
	normalized := strings.TrimSpace(strings.ToLower(title))

	// Remove common recipe prefixes/suffixes
	normalized = strings.ReplaceAll(normalized, "recipe", "")
	normalized = strings.ReplaceAll(normalized, "best", "")
	normalized = strings.ReplaceAll(normalized, "easy", "")
	normalized = strings.ReplaceAll(normalized, "homemade", "")
	normalized = strings.ReplaceAll(normalized, "the best", "")

	// Remove special characters
	normalized = regexp.MustCompile(`[^\w\s]`).ReplaceAllString(normalized, " ")

	// Replace multiple spaces with a single space
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")

	return strings.TrimSpace(normalized)
}

// isValidRecipeURL checks if a URL is a valid recipe URL
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
		"www.delish.com":        {"/cooking/recipe-ideas/", "/recipe/", "/recipes/"},
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

// calculateTitleSimilarity calculates the similarity between two titles
func calculateTitleSimilarity(title1, title2 string) float64 {
	// Normalize both titles
	norm1 := normalizeTitle(title1)
	norm2 := normalizeTitle(title2)

	// Simple Levenshtein distance calculation
	distance := levenshteinDistance(norm1, norm2)
	maxLength := math.Max(float64(len(norm1)), float64(len(norm2)))

	if maxLength == 0 {
		return 1.0
	}

	return 1.0 - float64(distance)/maxLength
}

// Levenshtein distance implementation
func levenshteinDistance(s1, s2 string) int {
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
func existingURL(urlToCheck string) bool {
	ctx := context.Background()

	// Create a bool query to check for exact URL match
	q := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("url.keyword", urlToCheck))

	// Execute the search
	result, err := client.Search().
		Index(IndexName).
		Query(q).
		Size(1).
		Do(ctx)

	if err != nil {
		logger.WriteWarning(fmt.Sprintf("Failed to check for existing URL: %v", err))
		return false
	}

	// Check if any results were found
	hits := result.TotalHits()
	return hits > 0
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

// CreatePage creates a new page in Elasticsearch
func CreatePage(p structs.Page) bool {
	ctx := context.Background()

	// Validate required fields
	if p.Title == "" || p.Name == "" || p.Ingredients == "" || p.Instructions == "" {
		logger.WriteWarning(fmt.Sprintf("Cannot create page, missing required fields: %s", p.URL))
		return false
	}

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

	// Check if URL already exists
	if existingURL(p.URL) {
		logger.WriteWarning(fmt.Sprintf("URL already exists in database: %s", p.URL))
		return false
	}

	// Basic validation of extracted data
	if len(p.Ingredients) < 10 || len(p.Instructions) < 20 {
		logger.WriteWarning(fmt.Sprintf("Recipe data seems incomplete: %s", p.Title))
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

// UpdatePage updates an existing page in Elasticsearch
func UpdatePage(id string, params map[string]interface{}) bool {
	ctx := context.Background()

	// Validate required fields
	if _, ok := params["title"].(string); !ok || params["title"].(string) == "" {
		logger.WriteWarning("Cannot update page, missing title")
		return false
	}

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
			if existingURL(url) {
				logger.WriteWarning(fmt.Sprintf("URL already exists in another document, skipping update: %s", url))
				return false
			}
		}
	}

	// Validate ingredients and instructions
	if ingredients, ok := params["ingredients"].(string); ok && len(ingredients) < 10 {
		logger.WriteWarning("Ingredients data seems incomplete, skipping update")
		return false
	}

	if instructions, ok := params["instructions"].(string); ok && len(instructions) < 20 {
		logger.WriteWarning("Instructions data seems incomplete, skipping update")
		return false
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

// SearchRecipes performs a search for recipes with the given query
func SearchRecipes(query string, from, size int) ([]structs.Page, int64, error) {
	ctx := context.Background()

	// Create a multi-match query for broader search
	multiMatchQuery := elastic.NewMultiMatchQuery(query,
		"title^3",         // Boost title matches
		"name^2",          // Boost name matches
		"ingredients^1.5", // Boost ingredient matches
		"description",
		"instructions",
	).Type("best_fields").Fuzziness("AUTO")

	// Create a search query
	searchQuery := elastic.NewBoolQuery().
		Should(multiMatchQuery).
		MinimumShouldMatch("1")

	// Execute the search
	result, err := client.Search().
		Index(IndexName).
		Query(searchQuery).
		From(from).
		Size(size).
		Sort("_score", false). // Sort by relevance
		Do(ctx)

	if err != nil {
		return nil, 0, err
	}

	// Parse the search results
	var recipes []structs.Page
	for _, hit := range result.Hits.Hits {
		var recipe structs.Page
		err := json.Unmarshal(hit.Source, &recipe)
		if err != nil {
			logger.WriteWarning(fmt.Sprintf("Failed to unmarshal recipe: %v", err))
			continue
		}
		recipes = append(recipes, recipe)
	}

	return recipes, result.TotalHits(), nil
}

// GetRecentRecipes gets the most recently crawled recipes
func GetRecentRecipes(limit int) ([]structs.Page, error) {
	ctx := context.Background()

	// Query for all recipes, sorted by crawl date
	searchQuery := elastic.NewMatchAllQuery()

	// Execute the search
	result, err := client.Search().
		Index(IndexName).
		Query(searchQuery).
		Sort("crawl_date", false). // Sort by crawl date, newest first
		Size(limit).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	// Parse the search results
	var recipes []structs.Page
	for _, hit := range result.Hits.Hits {
		var recipe structs.Page
		err := json.Unmarshal(hit.Source, &recipe)
		if err != nil {
			logger.WriteWarning(fmt.Sprintf("Failed to unmarshal recipe: %v", err))
			continue
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// FindSimilarRecipes finds recipes similar to the given recipe ID
func FindSimilarRecipes(recipeID string, limit int) ([]structs.Page, error) {
	ctx := context.Background()

	// First get the recipe
	getResult, err := client.Get().
		Index(IndexName).
		Id(recipeID).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	var recipe structs.Page
	if err := json.Unmarshal(getResult.Source, &recipe); err != nil {
		return nil, err
	}

	// Create a query based on recipe ingredients and title
	titleQuery := elastic.NewMatchQuery("title", recipe.Title).Boost(1.5)
	ingredientsQuery := elastic.NewMatchQuery("ingredients", recipe.Ingredients).Boost(2)

	searchQuery := elastic.NewBoolQuery().
		Should(titleQuery, ingredientsQuery).
		MustNot(elastic.NewIdsQuery().Ids(recipeID)). // Exclude the recipe itself
		MinimumShouldMatch("1")

	// Execute the search
	result, err := client.Search().
		Index(IndexName).
		Query(searchQuery).
		Size(limit).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	// Parse the search results
	var recipes []structs.Page
	for _, hit := range result.Hits.Hits {
		var similarRecipe structs.Page
		err := json.Unmarshal(hit.Source, &similarRecipe)
		if err != nil {
			logger.WriteWarning(fmt.Sprintf("Failed to unmarshal recipe: %v", err))
			continue
		}
		recipes = append(recipes, similarRecipe)
	}

	return recipes, nil
}
