package main

import (
	"fmt"
	"github.com/teris-io/shortid"
	"net/url"
	"os"
	"search-engine-indexer/src/elasticsearch"
	"search-engine-indexer/src/logger"
	"search-engine-indexer/src/scraper"
	"search-engine-indexer/src/structs"
	"sync"
	"time"
)

var queue = make(chan string)

func removeDuplicates(strs []string) []string {
	// Create a map to track first occurrence indices
	seen := make(map[string]int)

	// First pass: record the first occurrence of each string
	for i, s := range strs {
		if _, exists := seen[s]; !exists {
			seen[s] = i
		}
	}

	// Create result slice
	result := make([]string, 0, len(seen))

	// Keep only the first occurrence of each string
	for i, s := range strs {
		if seen[s] == i {
			result = append(result, s)
		}
	}

	return result
}

func crawlURL(url string) {
	// Extract links, title and description
	s := scraper.NewScraper(url)
	if s == nil {
		return
	}

	links := s.Links()
	recipeData := s.GetRecipeData()

	// Get basic metadata
	title := recipeData["title"]
	description := recipeData["description"]
	body := s.Body()

	fmtdString := fmt.Sprintf("Title: %s", title)
	logger.WriteInfo(fmtdString)

	// Check if the page exists
	existsLink, page := elasticsearch.ExistingPage(title)

	if !existsLink {
		// Create the new page in the database
		id, _ := shortid.Generate()
		newPage := structs.Page{
			ID:           id,
			Title:        title,
			Description:  description,
			Body:         body,
			URL:          url,
			Image:        recipeData["image"],
			Name:         recipeData["name"],
			PrepTime:     recipeData["prep_time"],
			CookTime:     recipeData["cook_time"],
			TotalTime:    recipeData["total_time"],
			Calories:     recipeData["calories"],
			Servings:     recipeData["servings"],
			Ingredients:  recipeData["ingredients"],
			Instructions: recipeData["instructions"],
			SourceSite:   extractSourceSite(url),
			CrawlDate:    time.Now(),
		}

		success := elasticsearch.CreatePage(newPage)
		if !success {
			return
		}

		fmt.Println("Page", url, "created")
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
			"source_site":  extractSourceSite(url),
		}

		success := elasticsearch.UpdatePage(page.ID, params)

		if !success {
			return
		}

		fmt.Println("Page", title, "with ID", page.ID, "updated")
	}

	// Queue new links for crawling
	for _, link := range links {
		go func(l string) {
			queue <- l
		}(link)
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

func worker(wg *sync.WaitGroup, id int) {
	for link := range queue {
		crawlURL(link)
	}

	wg.Done()
}

func checkIndexPresence() {
	elasticsearch.NewElasticSearchClient()
	exists := elasticsearch.ExistsIndex(elasticsearch.IndexName)
	if !exists {
		elasticsearch.CreateIndex(elasticsearch.IndexName)
	}
}

// Allocate workers and start crawling with the first URL
func startCrawling(start string) {
	checkIndexPresence()

	var wg sync.WaitGroup
	numberOfWorkers := 10

	// Send first url to channel
	go func(s string) {
		queue <- s
	}(start)

	// Create worker pool with numberOfWorkers workers
	wg.Add(numberOfWorkers)
	for i := 1; i <= numberOfWorkers; i++ {
		go worker(&wg, i)
	}

	wg.Wait()
}

func deleteIndex() {
	elasticsearch.NewElasticSearchClient()
	elasticsearch.DeleteIndex()
}

// Update main.go to crawl multiple recipe sites
func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not option provided, please specify one of the options below:")
		fmt.Println()
		fmt.Println("1. If you want to crawl recipe sites:")
		fmt.Println("\tgo run *.go index")
		fmt.Println()
		fmt.Println("2. If you want to crawl a specific recipe site:")
		fmt.Println("\tgo run *.go index URL")
		fmt.Println()
		fmt.Println("3. If you want to delete the pages index from elastic search:")
		fmt.Println("\tgo run *.go delete")
		return
	}

	switch args[1] {
	case "index":
		// List of popular recipe sites to start crawling
		startURLs := []string{
			"https://www.delish.com/cooking/recipe-ideas/",
			"https://www.allrecipes.com/recipes/",
			"https://www.foodnetwork.com/recipes",
			"https://www.epicurious.com/recipes",
			"https://www.simplyrecipes.com/recipes/",
		}

		// If URL is provided, use only that one
		if len(args) >= 3 {
			startURLs = []string{args[2]}
		}

		// Start crawling each URL
		for _, url := range startURLs {
			fmt.Printf("Starting crawl from: %s\n", url)
			startCrawling(url)
		}
	case "delete":
		deleteIndex()
	}
}
