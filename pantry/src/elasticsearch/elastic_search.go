package elasticsearch

// Elastic search client
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"search-engine-indexer/src/logger"
	"search-engine-indexer/src/structs"
	"strings"
	"time"

	elastic "github.com/olivere/elastic/v7"
)

const (
	IndexName    = "recipes"
	IndexMapping = `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0,
			"analysis": {
	      "analyzer": {
	        "clean_html": {
						"type": "standard",
	          "char_filter": ["html_strip"]
	        }
	      }
	    }
		},
		"mappings":{
			"properties":{
				"title": {
					"type": "text"
				},
				"description": {
					"type": "text"
				},
				"body": {
					"type": "text"
				},
				"url": {
					"type": "text"
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

// isValidDelishURL checks if the URL matches the expected Delish.com pattern
func isValidDelishURL(rawURL string) bool {
	// Must start with https://www.delish.com
	if !strings.HasPrefix(rawURL, "https://www.delish.com/") {
		return false
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Must be exactly www.delish.com
	if parsedURL.Host != "www.delish.com" {
		return false
	}

	// Split the path into segments
	segments := strings.Split(parsedURL.Path, "/")

	// Remove empty segments
	var cleanSegments []string
	for _, s := range segments {
		if s != "" {
			cleanSegments = append(cleanSegments, s)
		}
	}

	// Check for duplicate segments
	seen := make(map[string]bool)
	for _, segment := range cleanSegments {
		if seen[segment] {
			return false // Found duplicate segment
		}
		seen[segment] = true
	}

	// Check if the path follows the expected pattern:
	// /cooking/recipe-ideas/{recipe-id}/{recipe-name}
	if len(cleanSegments) != 4 ||
		cleanSegments[0] != "cooking" ||
		cleanSegments[1] != "recipe-ideas" {
		return false
	}

	// Additional check for any form of path duplication
	if strings.Count(rawURL, "/cooking/recipe-ideas/") > 1 {
		return false
	}

	return true
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
		Size(1).
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

// CreatePage adds a new page to the database only if the URL doesn't already exist
func CreatePage(p structs.Page) bool {
	ctx := context.Background()

	// Convert relative URL to absolute if needed
	if strings.HasPrefix(p.URL, "/") {
		p.URL = "https://www.delish.com" + p.URL
	}

	// Validate URL format
	if !isValidDelishURL(p.URL) {
		logger.WriteWarning(fmt.Sprintf("Invalid URL format or contains duplicates: %s", p.URL))
		return false
	}

	// Log the URL we're checking
	logger.WriteInfo(fmt.Sprintf("Checking for existing URL: %s", p.URL))

	// Check if URL already exists - do this check first before anything else
	if existingURL(client, IndexName, p.URL) {
		logger.WriteWarning(fmt.Sprintf("URL already exists in database, skipping creation: %s", p.URL))
		return false
	}

	// Double-check with a direct search
	searchResult, err := client.Search().
		Index(IndexName).
		Query(elastic.NewTermQuery("url.keyword", p.URL)).
		Size(1).
		Do(ctx)

	if err != nil {
		logger.WriteWarning(fmt.Sprintf("Error during URL search: %v", err))
		return false
	}

	if searchResult.TotalHits() > 0 {
		logger.WriteWarning(fmt.Sprintf("URL found in second check, skipping creation: %s", p.URL))
		return false
	}

	// Normalize the title after URL checks
	p.Title = strings.TrimSpace(strings.ToLower(p.Title))

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

	logger.WriteInfo(fmt.Sprintf("Successfully created new page - Title: %s, URL: %s", p.Title, p.URL))
	return true
}

// UpdatePage updates an existing page only if the URL is valid and unique
func UpdatePage(id string, params map[string]interface{}) bool {
	ctx := context.Background()

	// If URL is being updated, validate it
	if url, ok := params["url"].(string); ok {
		// Convert relative URL to absolute if needed
		if strings.HasPrefix(url, "/") {
			url = "https://www.delish.com" + url
			params["url"] = url
		}

		if !isValidDelishURL(url) {
			logger.WriteWarning(fmt.Sprintf("Invalid URL format or contains duplicates: %s", url))
			return false
		}

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
