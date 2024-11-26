package main

import (
	"fmt"
	"github.com/teris-io/shortid"
	"os"
	"search-engine-indexer/src/elasticsearch"
	"search-engine-indexer/src/logger"
	"search-engine-indexer/src/scraper"
	"search-engine-indexer/src/structs"
	"strings"
	"sync"
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

	var joinedUrl string
	var newLinks []string

	links := s.Links()
	fmt.Println(links)
	title, description := s.MetaDataInformation()
	body := s.Body()

	fmtdString := fmt.Sprintf("Title: %s", title)
	logger.WriteInfo(fmtdString)

	fmt.Println(url)

	for _, link := range links {
		splitUrl := strings.Split(link, "/")

		splitUrl = removeDuplicates(splitUrl)

		joinedUrl = strings.Join(splitUrl, "/")

		fmtdString = fmt.Sprintf("Url being used: %s", joinedUrl)
		logger.WriteInfo(fmtdString)

		newLinks = append(newLinks, link)
	}

	// Check if the page exists
	existsLink, page := elasticsearch.ExistingPage(joinedUrl)

	if existsLink {
		// Update the page in database
		params := map[string]interface{}{
			"title":       title,
			"description": description,
			"body":        body,
		}

		success := elasticsearch.UpdatePage(page.ID, params)
		if !success {
			return
		}

		fmt.Println("Page", joinedUrl, "with ID", page.ID, "update")
	} else {
		// Create the new page in the database
		id, _ := shortid.Generate()
		newPage := structs.Page{
			ID:          id,
			Title:       title,
			Description: description,
			Body:        body,
			URL:         joinedUrl,
		}

		success := elasticsearch.CreatePage(newPage)
		if !success {
			return
		}

		fmt.Println("Page", url, "created")

		for _, link := range links {
			go func(l string) {
				queue <- l
			}(link)
		}
	}
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

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not option provided, please specify one of the options below:")
		fmt.Println()
		fmt.Println("1. If you want to crawl the internet:")
		fmt.Println("\tgo run *.go index CRAWLING_START_URL")
		fmt.Println()
		fmt.Println("2. If you want to delete the pages index from elastic search:")
		fmt.Println("\tgo run *.go delete")
		return
	}

	switch args[1] {
	case "index":
		startCrawling(args[2])
	case "delete":
		deleteIndex()
	}
}
