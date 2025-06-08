package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/teris-io/shortid"
	"search-engine-indexer/src/elasticsearch"
	"search-engine-indexer/src/logger"
	"search-engine-indexer/src/scraper"
	"search-engine-indexer/src/structs"
	"sync"
)

// Global variables
var (
	// Queue for URLs to be crawled
	queue = make(chan string, 1000)

	// Set to track URLs that have been added to the queue already
	// to prevent duplicate crawling
	queuedURLs = sync.Map{}

	// Track successfully crawled URLs
	crawledURLs = sync.Map{}

	// Semaphore to limit concurrent requests to a domain
	domainSemaphores = make(map[string]*Semaphore)
	domainLock       sync.Mutex

	// Configuration
	maxCrawlDepth        = 3
	concurrentWorkers    = 10
	crawlDelayPerDomain  = 1 * time.Second
	maxRequestsPerDomain = 5

	// Debug mode for more verbose logging
	debugMode = false
)

// Custom Semaphore implementation for rate limiting
type Semaphore struct {
	c chan struct{}
}

// NewSemaphore creates a new semaphore with the given limit
func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		c: make(chan struct{}, limit),
	}
}

// Acquire acquires a semaphore token
func (s *Semaphore) Acquire() {
	s.c <- struct{}{}
}

// Release releases a semaphore token
func (s *Semaphore) Release() {
	<-s.c
}

// getDomainSemaphore returns a semaphore for the given domain,
// creating one if it doesn't exist
func getDomainSemaphore(domain string) *Semaphore {
	domainLock.Lock()
	defer domainLock.Unlock()

	if sem, ok := domainSemaphores[domain]; ok {
		return sem
	}

	sem := NewSemaphore(maxRequestsPerDomain)
	domainSemaphores[domain] = sem
	return sem
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(strs []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, s := range strs {
		if _, exists := seen[s]; !exists {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}

// isRecipeListingPage determines if a URL is a recipe listing/category page
func isRecipeListingPage(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check path patterns for listing pages
	path := parsedURL.Path

	// Common listing page patterns
	listingPatterns := []string{
		"/recipes/",
		"/cooking/recipe-ideas/",
		"/recipe-ideas/",
		"/recipes-a-z/",
		"/category/",
		"/collections/",
		"/meal-type/",
		"/cuisines/",
		"/cooking-method/",
		"/holidays-events/",
	}

	// Check if the path matches any listing pattern but doesn't have additional segments
	// that would indicate a specific recipe
	for _, pattern := range listingPatterns {
		if strings.Contains(path, pattern) {
			// If the path is exactly the pattern or has only a trailing slash,
			// it's likely a listing page
			if path == pattern || path == pattern+"/" ||
				// Allow for one extra path segment for category pages
				strings.Count(path, "/") <= strings.Count(pattern, "/")+1 {
				return true
			}
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
	}

	for _, indicator := range recipeIndicators {
		if strings.Contains(path, indicator) {
			return true
		}
	}

	// Check for specific patterns we know lead to recipe detail pages
	host := parsedURL.Host
	switch {
	case strings.Contains(host, "delish.com"):
		return strings.Contains(path, "/recipe/") || strings.Contains(path, "/recipes/")
	case strings.Contains(host, "allrecipes.com"):
		return strings.Contains(path, "/recipe/") || strings.Contains(path, "/recipes/") && !isRecipeListingPage(urlStr)
	case strings.Contains(host, "foodnetwork.com"):
		return strings.Contains(path, "/recipes/") && strings.Count(path, "/") >= 3
	case strings.Contains(host, "epicurious.com"):
		return strings.Contains(path, "/recipes/") && !isRecipeListingPage(urlStr)
	case strings.Contains(host, "simplyrecipes.com"):
		return strings.Contains(path, "/recipes/") && strings.Count(path, "/") >= 3
	}

	return false
}

// crawlURL crawls a single URL and extracts recipe data
func crawlURL(urlStr string, depth int) {
	// Check if we've reached maximum crawl depth
	if depth > maxCrawlDepth {
		logger.WriteInfo(fmt.Sprintf("Maximum crawl depth reached for URL: %s", urlStr))
		return
	}

	// Extract domain for rate limiting
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		logger.WriteError(fmt.Sprintf("Failed to parse URL: %s - %v", urlStr, err))
		return
	}
	domain := parsedURL.Host

	// Apply rate limiting per domain
	sem := getDomainSemaphore(domain)
	sem.Acquire()
	defer sem.Release()

	// Add a delay to avoid overloading the server
	time.Sleep(crawlDelayPerDomain)

	// Extract links, title and description
	s := scraper.NewScraper(urlStr)
	if s == nil {
		logger.WriteError(fmt.Sprintf("Failed to create scraper for URL: %s", urlStr))
		return
	}

	// Mark this URL as crawled
	crawledURLs.Store(urlStr, true)

	// Get links for further crawling
	links := s.Links()
	logger.WriteInfo(fmt.Sprintf("Found %d links on page: %s", len(links), urlStr))

	// Check if this is a recipe listing page or an actual recipe page
	if isRecipeListingPage(urlStr) {
		logger.WriteInfo(fmt.Sprintf("Processing as a recipe listing page: %s", urlStr))

		// For listing pages, just extract links and queue them for crawling
		for _, link := range links {
			// Check if we've already queued this URL
			if _, exists := queuedURLs.Load(link); !exists {
				// Mark URL as queued
				queuedURLs.Store(link, true)

				// Add to crawl queue with incremented depth
				go func(l string, d int) {
					queue <- l
				}(link, depth+1)
			}
		}

		return
	}

	// Check if this is a likely recipe page
	isRecipe := isLikelyRecipePage(urlStr)

	// For recipe pages, extract recipe data
	recipeData := s.GetRecipeData()

	// Get basic metadata
	title := recipeData["title"]
	description := recipeData["description"]
	body := s.Body()

	logger.WriteInfo(fmt.Sprintf("Processing potential recipe page: %s", title))

	// Debug mode - log all extracted data
	if debugMode {
		logger.WriteInfo(fmt.Sprintf("Extracted data for URL %s:", urlStr))
		for key, value := range recipeData {
			logger.WriteInfo(fmt.Sprintf("  %s: %s", key, value))
		}
	}

	// Check if the page exists
	existsLink, page := elasticsearch.ExistingPage(title)

	// Get a unique ID for new pages
	id, _ := shortid.Generate()

	// More flexible data checks
	hasName := recipeData["name"] != "" || title != ""
	hasIngredients := recipeData["ingredients"] != "" ||
		strings.Contains(body, "ingredient") ||
		strings.Contains(body, "Ingredient")

	hasInstructions := recipeData["instructions"] != "" ||
		strings.Contains(body, "direction") ||
		strings.Contains(body, "Direction") ||
		strings.Contains(body, "instruction") ||
		strings.Contains(body, "Instruction") ||
		strings.Contains(body, "steps") ||
		strings.Contains(body, "Steps") ||
		strings.Contains(body, "method") ||
		strings.Contains(body, "Method") ||
		strings.Contains(body, "preparation") ||
		strings.Contains(body, "Preparation")

	// Check if it's a likely recipe based on URL pattern AND has some recipe-like content
	hasMinimumData := hasName && (hasIngredients || hasInstructions)

	// If we're unsure, but the URL strongly suggests it's a recipe, try harder to extract data
	if !hasMinimumData && isRecipe {
		logger.WriteInfo(fmt.Sprintf("URL %s appears to be a recipe but missing some data, attempting recovery", urlStr))

		// If name is missing, use title
		if recipeData["name"] == "" {
			recipeData["name"] = title
			hasName = true
		}

		// If we're missing ingredients, do a wider search in the HTML
		if recipeData["ingredients"] == "" && hasIngredients {
			logger.WriteInfo(fmt.Sprintf("Ingredients text found in body for URL %s but not structured, proceeding anyway", urlStr))
		}

		// If we're missing instructions, do a wider search in the HTML
		if recipeData["instructions"] == "" && hasInstructions {
			logger.WriteInfo(fmt.Sprintf("Instructions text found in body for URL %s but not structured, proceeding anyway", urlStr))
		}

		// Re-evaluate minimum data requirement
		hasMinimumData = hasName && (hasIngredients || hasInstructions)
	}

	if !hasMinimumData {
		logger.WriteWarning(fmt.Sprintf("Skipping URL %s - Missing essential recipe data", urlStr))

		// Log what's missing to help debug
		if !hasName {
			logger.WriteWarning(fmt.Sprintf("  Missing name/title in URL %s", urlStr))
		}
		if !hasIngredients {
			logger.WriteWarning(fmt.Sprintf("  Missing ingredients in URL %s", urlStr))
		}
		if !hasInstructions {
			logger.WriteWarning(fmt.Sprintf("  Missing instructions in URL %s", urlStr))
		}

		// Even if we skip storing this page, we still queue its links for crawling
		for _, link := range links {
			if _, exists := queuedURLs.Load(link); !exists {
				queuedURLs.Store(link, true)
				go func(l string, d int) {
					queue <- l
				}(link, depth+1)
			}
		}

		return
	}

	// Use title as name if name is missing
	if recipeData["name"] == "" {
		recipeData["name"] = title
	}

	if !existsLink {
		// Create the new page in the database
		newPage := structs.Page{
			ID:           id,
			Title:        title,
			Description:  description,
			Body:         body,
			URL:          urlStr,
			Image:        recipeData["image"],
			Name:         recipeData["name"],
			PrepTime:     recipeData["prep_time"],
			CookTime:     recipeData["cook_time"],
			TotalTime:    recipeData["total_time"],
			Calories:     recipeData["calories"],
			Servings:     recipeData["servings"],
			Ingredients:  recipeData["ingredients"],
			Instructions: recipeData["instructions"],
			SourceSite:   extractSourceSite(urlStr),
			CrawlDate:    time.Now(),
		}

		success := elasticsearch.CreatePage(newPage)
		if !success {
			logger.WriteError(fmt.Sprintf("Failed to create page for URL: %s", urlStr))

			// Even if storing fails, still queue links for crawling
			for _, link := range links {
				if _, exists := queuedURLs.Load(link); !exists {
					queuedURLs.Store(link, true)
					go func(l string, d int) {
						queue <- l
					}(link, depth+1)
				}
			}

			return
		}

		logger.WriteInfo(fmt.Sprintf("Created new recipe: %s - %s", newPage.ID, urlStr))

		// Save a copy to the filesystem for backup
		saveRecipeToFile(newPage)
	} else {
		// Update the page in database
		params := map[string]interface{}{
			"title":        title,
			"description":  description,
			"body":         body,
			"image":        recipeData["image"],
			"name":         recipeData["name"],
			"prep_time":    recipeData["prep_time"],
			"cook_time":    recipeData["cook_time"],
			"total_time":   recipeData["total_time"],
			"calories":     recipeData["calories"],
			"servings":     recipeData["servings"],
			"ingredients":  recipeData["ingredients"],
			"instructions": recipeData["instructions"],
			"source_site":  extractSourceSite(urlStr),
		}

		success := elasticsearch.UpdatePage(page.ID, params)
		if !success {
			logger.WriteError(fmt.Sprintf("Failed to update page for URL: %s", urlStr))
		} else {
			logger.WriteInfo(fmt.Sprintf("Updated page %s (%s)", page.ID, title))
		}
	}

	// Queue new links for crawling
	for _, link := range links {
		if _, exists := queuedURLs.Load(link); !exists {
			queuedURLs.Store(link, true)
			go func(l string, d int) {
				queue <- l
			}(link, depth+1)
		}
	}
}

// Helper function to extract source site from URL
func extractSourceSite(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}

// saveRecipeToFile saves a backup of the recipe to a JSON file
func saveRecipeToFile(page structs.Page) {
	// Create backups directory if it doesn't exist
	backupDir := "recipe_backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		logger.WriteError(fmt.Sprintf("Failed to create backup directory: %v", err))
		return
	}

	// Create a filename based on the page ID and title
	safeName := strings.ReplaceAll(page.Title, " ", "_")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, ":", "_")
	filename := filepath.Join(backupDir, fmt.Sprintf("%s_%s.json", page.ID, safeName))

	// Marshal the page data to JSON
	jsonData, err := json.MarshalIndent(page, "", "  ")
	if err != nil {
		logger.WriteError(fmt.Sprintf("Failed to marshal page data: %v", err))
		return
	}

	// Write the JSON data to file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		logger.WriteError(fmt.Sprintf("Failed to write backup file: %v", err))
		return
	}

	logger.WriteInfo(fmt.Sprintf("Saved recipe backup to %s", filename))
}

// worker function processes URLs from the queue
func worker(wg *sync.WaitGroup, id int) {
	logger.WriteInfo(fmt.Sprintf("Worker %d started", id))

	for link := range queue {
		// Extract depth from metadata or default to 0
		depth := 0
		crawlURL(link, depth)
	}

	logger.WriteInfo(fmt.Sprintf("Worker %d finished", id))
	wg.Done()
}

// checkIndexPresence ensures the Elasticsearch index exists
func checkIndexPresence() {
	elasticsearch.NewElasticSearchClient()
	exists := elasticsearch.ExistsIndex(elasticsearch.IndexName)
	if !exists {
		logger.WriteInfo("Creating Elasticsearch index...")
		elasticsearch.CreateIndex(elasticsearch.IndexName)
	} else {
		logger.WriteInfo("Elasticsearch index already exists")
	}
}

// startCrawling initializes workers and begins crawling from the starting URL
func startCrawling(startURLs []string) {
	// Ensure Elasticsearch index exists
	checkIndexPresence()

	var wg sync.WaitGroup

	// Set up termination channel with timeout
	timeout := 30 * time.Minute
	done := make(chan bool)

	// Send initial URLs to channel
	for _, startURL := range startURLs {
		go func(url string) {
			queuedURLs.Store(url, true)
			queue <- url
		}(startURL)
	}

	// Create worker pool
	wg.Add(concurrentWorkers)
	for i := 1; i <= concurrentWorkers; i++ {
		go worker(&wg, i)
	}

	// Set up timeout
	go func() {
		time.Sleep(timeout)
		logger.WriteInfo(fmt.Sprintf("Crawler timeout reached after %v", timeout))
		done <- true
	}()

	// Wait for workers to finish or timeout
	go func() {
		wg.Wait()
		done <- true
	}()

	// Wait for done signal
	<-done

	// Print summary
	var crawledCount int
	crawledURLs.Range(func(_, _ interface{}) bool {
		crawledCount++
		return true
	})

	logger.WriteInfo(fmt.Sprintf("Crawling completed. Processed %d URLs.", crawledCount))
}

// deleteIndex removes the Elasticsearch index
func deleteIndex() {
	elasticsearch.NewElasticSearchClient()
	elasticsearch.DeleteIndex()
}

// getPopularRecipeSites returns URLs for popular recipe sites
func getPopularRecipeSites() []string {
	return []string{
		"https://pinchofyum.com/recipes",
		"https://minimalistbaker.com/recipes", 
		"https://cookieandkate.com/recipes",
		"https://loveandlemons.com/recipes",
		"https://smittenkitchen.com/recipes",
		"https://seriouseats.com/recipes",
		"https://halfbakedharvest.com/category/recipes",
		"https://101cookbooks.com/recipes",
		"https://food52.com/recipes",
		"https://budgetbytes.com/category/recipes",
		"https://thewoksoflife.com/recipes",
		"https://www.delish.com/cooking/recipe-ideas/",
		"https://www.allrecipes.com/recipes/",
		"https://www.foodnetwork.com/recipes",
		"https://www.epicurious.com/recipes",
		"https://www.simplyrecipes.com/recipes/",
	}
}

// startRecipeCrawling starts crawling popular recipe sites
func startRecipeCrawling() {
	sites := getPopularRecipeSites()
	fmt.Printf("Starting to crawl %d popular recipe sites...\n", len(sites))
	
	logger.WriteInfo(fmt.Sprintf("Starting crawler with parameters:"))
	logger.WriteInfo(fmt.Sprintf("  Workers: %d", concurrentWorkers))
	logger.WriteInfo(fmt.Sprintf("  Max Depth: %d", maxCrawlDepth))
	logger.WriteInfo(fmt.Sprintf("  Delay: %v", crawlDelayPerDomain))
	logger.WriteInfo(fmt.Sprintf("  Max Requests Per Domain: %d", maxRequestsPerDomain))
	logger.WriteInfo(fmt.Sprintf("  Debug Mode: %t", debugMode))
	logger.WriteInfo(fmt.Sprintf("  Starting URLs: %v", sites))

	startCrawling(sites)

	// Print summary
	var crawledCount int
	crawledURLs.Range(func(_, _ interface{}) bool {
		crawledCount++
		return true
	})
	fmt.Printf("Crawling completed. Processed %d URLs.\n", crawledCount)
}

// main function handles command line arguments and starts the crawler
func main() {
	args := os.Args

	// Setup logging
	logFile, err := os.OpenFile("recipe_crawler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if len(args) < 2 {
		fmt.Println("Not option provided, please specify one of the options below:")
		fmt.Println()
		fmt.Println("1. If you want to crawl popular recipe sites:")
		fmt.Println("\tgo run *.go recipes")
		fmt.Println()
		fmt.Println("2. If you want to crawl a specific recipe site:")
		fmt.Println("\tgo run *.go index URL")
		fmt.Println()
		fmt.Println("3. If you want to crawl with custom parameters:")
		fmt.Println("\tgo run *.go index URL -workers=20 -depth=5 -delay=2 -debug=true")
		fmt.Println()
		fmt.Println("4. If you want to delete the pages index from elastic search:")
		fmt.Println("\tgo run *.go delete")
		fmt.Println()
		fmt.Println("5. If you want to test a specific URL:")
		fmt.Println("\tgo run *.go test-url URL")
		return
	}

	// Parse command-line arguments for custom crawler settings
	for i := 2; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-workers=") {
			fmt.Sscanf(arg[9:], "%d", &concurrentWorkers)
		} else if strings.HasPrefix(arg, "-depth=") {
			fmt.Sscanf(arg[7:], "%d", &maxCrawlDepth)
		} else if strings.HasPrefix(arg, "-delay=") {
			var delay float64
			fmt.Sscanf(arg[7:], "%f", &delay)
			crawlDelayPerDomain = time.Duration(delay * float64(time.Second))
		} else if strings.HasPrefix(arg, "-max-requests=") {
			fmt.Sscanf(arg[14:], "%d", &maxRequestsPerDomain)
		} else if strings.HasPrefix(arg, "-debug=") {
			fmt.Sscanf(arg[7:], "%t", &debugMode)
		}
	}

	switch args[1] {
	case "recipes":
		startRecipeCrawling()

	case "index":
		// Default to popular recipe sites if no URL provided
		startURLs := getPopularRecipeSites()

		// If URL is provided, use only that one
		if len(args) >= 3 && !strings.HasPrefix(args[2], "-") {
			startURLs = []string{args[2]}
		}

		// Log parameters
		logger.WriteInfo(fmt.Sprintf("Starting crawler with parameters:"))
		logger.WriteInfo(fmt.Sprintf("  Workers: %d", concurrentWorkers))
		logger.WriteInfo(fmt.Sprintf("  Max Depth: %d", maxCrawlDepth))
		logger.WriteInfo(fmt.Sprintf("  Delay: %v", crawlDelayPerDomain))
		logger.WriteInfo(fmt.Sprintf("  Max Requests Per Domain: %d", maxRequestsPerDomain))
		logger.WriteInfo(fmt.Sprintf("  Debug Mode: %t", debugMode))
		logger.WriteInfo(fmt.Sprintf("  Starting URLs: %v", startURLs))

		// Start crawling each URL
		fmt.Printf("Starting crawler with %d workers\n", concurrentWorkers)
		fmt.Printf("Max crawl depth: %d\n", maxCrawlDepth)
		startCrawling(startURLs)

		// Print summary
		var crawledCount int
		crawledURLs.Range(func(_, _ interface{}) bool {
			crawledCount++
			return true
		})
		fmt.Printf("Crawling completed. Processed %d URLs.\n", crawledCount)

	case "delete":
		deleteIndex()
		fmt.Println("Index deleted successfully")

	case "test-url":
		if len(args) < 3 {
			fmt.Println("Please provide a URL to test")
			return
		}

		testURL := args[2]
		fmt.Printf("Testing URL: %s\n", testURL)

		// Enable debug mode for testing
		debugMode = true

		// Create scraper
		s := scraper.NewScraper(testURL)
		if s == nil {
			fmt.Println("Failed to create scraper")
			return
		}

		// Extract recipe data
		recipeData := s.GetRecipeData()

		// Print recipe data
		fmt.Println("Recipe Data:")
		fmt.Printf("  Title: %s\n", recipeData["title"])
		fmt.Printf("  Name: %s\n", recipeData["name"])
		fmt.Printf("  Description: %s\n", recipeData["description"])
		fmt.Printf("  Image: %s\n", recipeData["image"])
		fmt.Printf("  Prep Time: %s\n", recipeData["prep_time"])
		fmt.Printf("  Cook Time: %s\n", recipeData["cook_time"])
		fmt.Printf("  Total Time: %s\n", recipeData["total_time"])
		fmt.Printf("  Calories: %s\n", recipeData["calories"])
		fmt.Printf("  Servings: %s\n", recipeData["servings"])

		// Print ingredients
		fmt.Println("  Ingredients:")
		ingredients := strings.Split(recipeData["ingredients"], ";")
		for i, ingredient := range ingredients {
			fmt.Printf("    %d. %s\n", i+1, ingredient)
		}

		// Print instructions
		fmt.Println("  Instructions:")
		instructions := strings.Split(recipeData["instructions"], ";")
		for i, instruction := range instructions {
			fmt.Printf("    %d. %s\n", i+1, instruction)
		}

		// Check if this is a recipe listing or detail page
		fmt.Println("\nURL Analysis:")
		fmt.Printf("  Is recipe listing page: %t\n", isRecipeListingPage(testURL))
		fmt.Printf("  Is likely recipe page: %t\n", isLikelyRecipePage(testURL))

		// Check minimum data requirements
		hasName := recipeData["name"] != "" || recipeData["title"] != ""
		hasIngredients := recipeData["ingredients"] != ""
		hasInstructions := recipeData["instructions"] != ""
		body := s.Body()
		if !hasIngredients {
			hasIngredients = strings.Contains(body, "ingredient") || strings.Contains(body, "Ingredient")
		}
		if !hasInstructions {
			hasInstructions = strings.Contains(body, "direction") ||
				strings.Contains(body, "Direction") ||
				strings.Contains(body, "instruction") ||
				strings.Contains(body, "Instruction") ||
				strings.Contains(body, "steps") ||
				strings.Contains(body, "Steps") ||
				strings.Contains(body, "method") ||
				strings.Contains(body, "Method") ||
				strings.Contains(body, "preparation") ||
				strings.Contains(body, "Preparation")
		}

		fmt.Println("\nData Requirements Check:")
		fmt.Printf("  Has name/title: %t\n", hasName)
		fmt.Printf("  Has ingredients: %t\n", hasIngredients)
		fmt.Printf("  Has instructions: %t\n", hasInstructions)
		fmt.Printf("  Meets minimum requirements: %t\n", hasName && (hasIngredients || hasInstructions))

		// Print links found
		links := s.Links()
		fmt.Printf("\nFound %d links on the page\n", len(links))
		if len(links) > 0 {
			fmt.Println("First 5 links:")
			for i, link := range links {
				if i >= 5 {
					break
				}
				fmt.Printf("  %d. %s\n", i+1, link)
			}
		}

	default:
		fmt.Println("Unknown option:", args[1])
		fmt.Println("Valid options are: recipes, index, delete, test-url")
	}
}