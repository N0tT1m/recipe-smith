package elasticsearch

// Elastic search client
import (
	"context"
	"encoding/json"
	"fmt"
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
		logger.WriteWarning(fmt.Sprintf("Failed to connect to Elasticsearch after %d attempts", maxRetries))
		return nil
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://192.168.1.78:9200")
	if err != nil {
		// Handle error
		logger.WriteWarning(fmt.Sprintf("Failed to get Elasticsearch version: %v", err))
		return nil
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

// Helper function to extract domain from URL path
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

// Helper function to extract source site from URL
func extractSourceSite(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}

// isValidRecipeURL checks if a URL is a valid recipe URL
func isValidRecipeURL(rawURL string) bool {
	// List of allowed recipe domains - expanded
	allowedDomains := []string{
		"www.delish.com",
		"www.allrecipes.com",
		"allrecipes.com", // Some may not have www prefix
		"www.foodnetwork.com",
		"www.epicurious.com",
		"www.simplyrecipes.com",
		"www.bonappetit.com",
		"www.taste.com.au",
		"www.bbcgoodfood.com",
		"www.eatingwell.com",
		"www.seriouseats.com",
		"cooking.nytimes.com",
		"www.tasteofhome.com",
		"www.food.com",
		"www.yummly.com",
		"thepioneerwoman.com",
		"minimalistbaker.com",
		"pinchofyum.com",
		"www.budgetbytes.com",
		"sallysbakingaddiction.com",
		"www.recipetineats.com",
		"www.gimmesomeoven.com",
		"damndelicious.net",
		"www.marthastewart.com",
		"www.myrecipes.com",
		"www.101cookbooks.com",
		"www.skinnytaste.com",
		"cookieandkate.com",
		"smittenkitchen.com",
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Extract host without www. prefix for consistent matching
	host := parsedURL.Host
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	// Check if domain is in allowed list
	domainValid := false
	for _, domain := range allowedDomains {
		// Remove www. from the domain for consistent matching
		allowedHost := domain
		if strings.HasPrefix(allowedHost, "www.") {
			allowedHost = allowedHost[4:]
		}

		if host == allowedHost || parsedURL.Host == domain {
			domainValid = true
			break
		}
	}

	if !domainValid {
		return false
	}

	// Map of domain-specific validation patterns - expanded
	validPathPatterns := map[string][]string{
		"delish.com":                {"/cooking/recipe-ideas/", "/recipe/", "/recipes/"},
		"allrecipes.com":            {"/recipe/", "/recipes/", "/gallery/"},
		"foodnetwork.com":           {"/recipes/", "/recipe/", "/fn-dish/"},
		"epicurious.com":            {"/recipes/", "/recipe/", "/food/views/"},
		"simplyrecipes.com":         {"/recipes/", "/"},
		"bonappetit.com":            {"/recipe/", "/recipes/", "/story/"},
		"taste.com.au":              {"/recipes/", "/recipe/"},
		"bbcgoodfood.com":           {"/recipes/", "/recipe/"},
		"eatingwell.com":            {"/recipe/", "/recipes/"},
		"seriouseats.com":           {"/recipes/", "/"},
		"cooking.nytimes.com":       {"/recipes/", "/"},
		"tasteofhome.com":           {"/recipes/", "/recipe/"},
		"food.com":                  {"/recipe/", "/recipes/"},
		"yummly.com":                {"/recipe/", "/recipes/"},
		"thepioneerwoman.com":       {"/food-cooking/recipes/", "/food-cooking/"},
		"minimalistbaker.com":       {"/recipes/", "/recipe/"},
		"pinchofyum.com":            {"/recipe/", "/"},
		"budgetbytes.com":           {"/recipes/", "/recipe/"},
		"sallysbakingaddiction.com": {"/recipe/", "/"},
		"recipetineats.com":         {"/recipes/", "/recipe/", "/"},
		"gimmesomeoven.com":         {"/"},
		"damndelicious.net":         {"/recipe/", "/"},
		"marthastewart.com":         {"/recipe/", "/recipes/"},
		"myrecipes.com":             {"/recipe/", "/recipes/"},
		"101cookbooks.com":          {"/recipes/", "/recipe/"},
		"skinnytaste.com":           {"/recipe/", "/recipes/"},
		"cookieandkate.com":         {"/recipe/", "/"},
		"smittenkitchen.com":        {"/recipe/", "/"},
	}

	// Extract the host part without www. for consistent matching
	hostForPatterns := parsedURL.Host
	if strings.HasPrefix(hostForPatterns, "www.") {
		hostForPatterns = hostForPatterns[4:]
	}

	// Check domain-specific path patterns
	patterns, exists := validPathPatterns[hostForPatterns]
	if !exists {
		// If we don't have specific patterns, check if path has common recipe indicators
		path := parsedURL.Path
		commonRecipePatterns := []string{"/recipe/", "/recipes/", "-recipe", "recipe-"}
		for _, pattern := range commonRecipePatterns {
			if strings.Contains(path, pattern) {
				return true
			}
		}
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

// isLikelyRecipePage does a best-effort check to see if a URL likely leads to a recipe page
func isLikelyRecipePage(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	path := parsedURL.Path

	// Recipe indicators in the URL path
	recipeIndicators := []string{
		"/recipe/",
		"/recipes/",
		"-recipe",
		"recipe-",
		"-recipes",
		"recipes-",
		"/cook/",
		"how-to-make",
		"how-to-cook",
		"best-ever",
		"easy-",
		"-cake",
		"-soup",
		"-salad",
		"-pie",
		"-bread",
		"-cookie",
		"-meal",
		"-dinner",
		"-breakfast",
		"-lunch",
		"-dessert",
		"-sandwich",
		"-pizza",
		"-pasta",
	}

	for _, indicator := range recipeIndicators {
		if strings.Contains(path, indicator) {
			return true
		}
	}

	// Check if URL has a numeric ID which is common for recipe pages
	recipeIdPattern := regexp.MustCompile(`/\d+(/|$)`)
	if recipeIdPattern.MatchString(path) {
		return true
	}

	// Check for specific patterns we know lead to recipe detail pages
	host := parsedURL.Host
	// Remove www. for consistent matching
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	switch host {
	case "allrecipes.com":
		// AllRecipes has URLs like /recipe/12345/chocolate-cake/
		if strings.Contains(path, "/recipe/") && strings.Count(path, "/") >= 3 {
			return true
		}
	case "foodnetwork.com":
		// FoodNetwork has URLs like /recipes/food-network-kitchens/chocolate-cake-recipe-2109090
		if strings.Contains(path, "/recipes/") && strings.Count(path, "/") >= 3 {
			return true
		}
	case "epicurious.com":
		// Epicurious has URLs like /recipes/food/views/chocolate-cake-107885
		if strings.Contains(path, "/food/views/") {
			return true
		}
	case "simplyrecipes.com":
		// Simply Recipes has URLs like /recipes/chocolate-cake/
		if strings.Count(path, "/") >= 2 && !strings.HasSuffix(path, "/recipes/") {
			return true
		}
	}

	return false
}

// existingURL checks if a URL already exists in the database with error handling
func existingURL(urlToCheck string) (bool, error) {
	ctx := context.Background()

	// Create a bool query to check for exact URL match
	q := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("url.keyword", urlToCheck))

	// Execute the search with error handling
	result, err := client.Search().
		Index(IndexName).
		Query(q).
		Size(1).
		Do(ctx)

	if err != nil {
		return false, err
	}

	// Check if any results were found
	hits := result.TotalHits()
	return hits > 0, nil
}

// CreatePage creates a new page in Elasticsearch with more flexible validation
func CreatePage(p structs.Page) bool {
	ctx := context.Background()

	// More flexible required field validation
	if p.Title == "" && p.Name == "" {
		logger.WriteWarning(fmt.Sprintf("Cannot create page, missing both title and name: %s", p.URL))
		return false
	}

	// If title is missing but name exists, use name as title
	if p.Title == "" {
		p.Title = p.Name
		logger.WriteInfo(fmt.Sprintf("Using name as title for URL: %s", p.URL))
	}

	// If name is missing but title exists, use title as name
	if p.Name == "" {
		p.Name = p.Title
		logger.WriteInfo(fmt.Sprintf("Using title as name for URL: %s", p.URL))
	}

	// Allow empty ingredients if body contains ingredient-related text
	hasIngredientText := p.Ingredients != "" ||
		strings.Contains(strings.ToLower(p.Body), "ingredient") ||
		strings.Contains(p.Body, "Ingredient")

	// Allow empty instructions if body contains instruction-related text
	instructionKeywords := []string{"instructions", "directions", "steps", "method", "preparation"}
	hasInstructionText := p.Instructions != ""

	if !hasInstructionText {
		for _, keyword := range instructionKeywords {
			if strings.Contains(strings.ToLower(p.Body), keyword) ||
				strings.Contains(p.Body, strings.Title(keyword)) {
				hasInstructionText = true
				break
			}
		}
	}

	// Check minimum required data with more flexibility
	if !hasIngredientText || !hasInstructionText {
		// Special case for allrecipes.com - they have a specific structure
		if strings.Contains(p.URL, "allrecipes.com") {
			// For allrecipes, if we have the URL and title/name, let's assume it's valid
			// and create a placeholder for missing data
			if p.Ingredients == "" && hasIngredientText {
				p.Ingredients = "Ingredients mentioned in page but not structured"
				logger.WriteInfo(fmt.Sprintf("Created placeholder ingredients for AllRecipes URL: %s", p.URL))
			}

			if p.Instructions == "" && hasInstructionText {
				p.Instructions = "Instructions mentioned in page but not structured"
				logger.WriteInfo(fmt.Sprintf("Created placeholder instructions for AllRecipes URL: %s", p.URL))
			}
		} else {
			// For other sites, still enforce basic validation
			if !hasIngredientText {
				logger.WriteWarning(fmt.Sprintf("Cannot create page, missing ingredients: %s", p.URL))
				return false
			}

			if !hasInstructionText {
				logger.WriteWarning(fmt.Sprintf("Cannot create page, missing instructions: %s", p.URL))
				return false
			}
		}
	}

	// Convert relative URL to absolute if needed
	if strings.HasPrefix(p.URL, "/") {
		// Extract domain from context or use a default one
		domain := extractDomainFromURL(p.URL)
		if domain == "" {
			domain = "https://www.example.com" // Fallback
		}
		p.URL = domain + p.URL
		logger.WriteInfo(fmt.Sprintf("Converted relative URL to absolute: %s", p.URL))
	}

	// Improved URL validation for recipe sites
	if !isValidRecipeURL(p.URL) {
		// If URL doesn't match standard patterns but appears to be a recipe, allow it
		if isLikelyRecipePage(p.URL) {
			logger.WriteInfo(fmt.Sprintf("URL doesn't match standard patterns but appears to be a recipe: %s", p.URL))
		} else {
			logger.WriteWarning(fmt.Sprintf("Invalid recipe URL format: %s", p.URL))
			return false
		}
	}

	// Set source site from URL
	p.SourceSite = extractSourceSite(p.URL)

	// Set crawl date to current time
	p.CrawlDate = time.Now()

	// Check if URL already exists - with improved error handling
	urlExists, err := existingURL(p.URL)
	if err != nil {
		logger.WriteWarning(fmt.Sprintf("Error checking if URL exists: %v - proceeding anyway", err))
	} else if urlExists {
		logger.WriteWarning(fmt.Sprintf("URL already exists in database: %s", p.URL))
		return false
	}

	// Less stringent validation of extracted data length
	// We still want some basic data, but we'll be more lenient
	if len(p.Ingredients) < 5 && !strings.Contains(p.Ingredients, "placeholder") {
		logger.WriteWarning(fmt.Sprintf("Ingredients seem unusually short: %s", p.Title))
		// But proceed anyway - don't return false
	}

	if len(p.Instructions) < 10 && !strings.Contains(p.Instructions, "placeholder") {
		logger.WriteWarning(fmt.Sprintf("Instructions seem unusually short: %s", p.Title))
		// But proceed anyway - don't return false
	}

	// Create the new page with refresh to ensure immediate visibility
	_, err = client.Index().
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
			urlExists, err := existingURL(url)
			if err != nil {
				logger.WriteWarning(fmt.Sprintf("Error checking if URL exists: %v - proceeding anyway", err))
			} else if urlExists {
				logger.WriteWarning(fmt.Sprintf("URL already exists in another document, skipping update: %s", url))
				return false
			}
		}
	}

	// More flexible validation for ingredients and instructions
	if ingredients, ok := params["ingredients"].(string); ok && ingredients == "" {
		// If ingredients is provided but empty, add a placeholder
		params["ingredients"] = "Ingredients mentioned in page but not structured"
		logger.WriteInfo("Created placeholder ingredients for update")
	}

	if instructions, ok := params["instructions"].(string); ok && instructions == "" {
		// If instructions is provided but empty, add a placeholder
		params["instructions"] = "Instructions mentioned in page but not structured"
		logger.WriteInfo("Created placeholder instructions for update")
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

// calculateTitleSimilarity calculates the similarity between two titles
func calculateTitleSimilarity(title1, title2 string) float64 {
	// Normalize both titles
	norm1 := normalizeTitle(title1)
	norm2 := normalizeTitle(title2)

	// Use Levenshtein distance for similarity
	distance := levenshteinDistance(norm1, norm2)
	maxLength := float64(max(len(norm1), len(norm2)))

	if maxLength == 0 {
		return 1.0
	}

	return 1.0 - float64(distance)/maxLength
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
