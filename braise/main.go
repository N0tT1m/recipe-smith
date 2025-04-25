package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	elastic "github.com/olivere/elastic/v7"
)

var client *elastic.Client

const (
	IndexName = "recipes"
)

func initElastic() {
	var err error

	// Create a new Elasticsearch client
	client, err = elastic.NewClient(
		elastic.SetURL("http://192.168.1.78:9200"),
		elastic.SetSniff(false),
	)

	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Check if the Elasticsearch server is running
	info, code, err := client.Ping("http://192.168.1.78:9200").Do(context.Background())
	if err != nil {
		log.Fatalf("Error pinging Elasticsearch: %s", err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}

func main() {
	// Initialize Elasticsearch client
	initElastic()

	// Create a new router
	r := mux.NewRouter()

	r.HandleFunc("/api/recipes/search", searchRecipes).Methods("GET")
	r.HandleFunc("/api/recipes/all", getAllRecipes).Methods("GET")
	r.HandleFunc("/api/recipes/category/{category}", getRecipesByCategory).Methods("GET")
	r.HandleFunc("/api/recipes/recent", getRecentRecipes).Methods("GET")
	// This general route must come AFTER more specific routes
	r.HandleFunc("/api/recipes/{id}", getRecipe).Methods("GET")

	// Setup CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Apply middleware
	handler := corsMiddleware(r)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// Get all recipes
func getAllRecipes(w http.ResponseWriter, r *http.Request) {
	// Get page and size parameters for pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size < 1 {
		size = 10 // Default size
	}

	// Calculate from for pagination
	from := (page - 1) * size

	// Use a simple match_all query without any complex sorting or filtering
	searchService := client.Search().
		Index(IndexName).
		Query(elastic.NewMatchAllQuery()).
		From(from).
		Size(size)

	// Try to sort by name if possible, but make it optional
	// The error might be related to trying to sort on a field that doesn't exist
	// or isn't mapped correctly
	searchResult, err := searchService.Do(context.Background())

	if err != nil {
		log.Printf("Error getting all recipes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process search results
	var recipes []map[string]interface{}
	for _, hit := range searchResult.Hits.Hits {
		var recipe map[string]interface{}
		if err := json.Unmarshal(hit.Source, &recipe); err != nil {
			log.Printf("Error unmarshaling recipe: %s", err)
			continue
		}

		// Add the ID to the recipe
		recipe["id"] = hit.Id
		recipes = append(recipes, recipe)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// Get a single recipe by ID
func getRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the recipe from Elasticsearch
	result, err := client.Get().
		Index(IndexName).
		Id(id).
		Do(context.Background())

	if err != nil {
		if elastic.IsNotFound(err) {
			http.Error(w, "Recipe not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting recipe: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Convert raw source to JSON
	var recipe map[string]interface{}
	if err := json.Unmarshal(result.Source, &recipe); err != nil {
		log.Printf("Error unmarshaling recipe: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the ID to the recipe
	recipe["id"] = result.Id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Search recipes with query
func searchRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Get page and size parameters for pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size < 1 {
		size = 10 // Default size
	}

	// Calculate from for pagination
	from := (page - 1) * size

	// Create a multi-match query
	multiMatchQuery := elastic.NewMultiMatchQuery(query,
		"title^3", // Boost title field
		"name^3",
		"description^2",
		"ingredients^2",
		"instructions",
		"body",
	).Type("best_fields").TieBreaker(0.3)

	// Create search service
	searchResult, err := client.Search().
		Index(IndexName).
		Query(multiMatchQuery).
		From(from).
		Size(size).
		Sort("_score", false). // Sort by relevance
		Do(context.Background())

	if err != nil {
		log.Printf("Error searching recipes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process search results
	var recipes []map[string]interface{}
	for _, hit := range searchResult.Hits.Hits {
		var recipe map[string]interface{}
		if err := json.Unmarshal(hit.Source, &recipe); err != nil {
			log.Printf("Error unmarshaling recipe: %s", err)
			continue
		}

		// Add the ID to the recipe
		recipe["id"] = hit.Id
		recipes = append(recipes, recipe)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// Get recipes by category
func getRecipesByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := strings.ToLower(vars["category"])

	// Get page and size parameters for pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size < 1 {
		size = 10 // Default size
	}

	// Calculate from for pagination
	from := (page - 1) * size

	// Create a query for the category
	// This could be modified based on how you store category information
	matchQuery := elastic.NewBoolQuery().
		Should(
			elastic.NewMatchPhraseQuery("source_site", category),
			elastic.NewMatchPhraseQuery("categories", category),
		).
		MinimumShouldMatch("1")

	// Create search service
	searchResult, err := client.Search().
		Index(IndexName).
		Query(matchQuery).
		From(from).
		Size(size).
		Sort("crawl_date", false). // Sort by date (newest first)
		Do(context.Background())

	if err != nil {
		log.Printf("Error getting recipes by category: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process search results
	var recipes []map[string]interface{}
	for _, hit := range searchResult.Hits.Hits {
		var recipe map[string]interface{}
		if err := json.Unmarshal(hit.Source, &recipe); err != nil {
			log.Printf("Error unmarshaling recipe: %s", err)
			continue
		}

		// Add the ID to the recipe
		recipe["id"] = hit.Id
		recipes = append(recipes, recipe)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// Get recent recipes
func getRecentRecipes(w http.ResponseWriter, r *http.Request) {
	// Get page and size parameters for pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size < 1 {
		size = 10 // Default size
	}

	// Calculate from for pagination
	from := (page - 1) * size

	// Create a query to get all documents, sorted by crawl date
	query := elastic.NewMatchAllQuery()

	// Get current time
	now := time.Now()

	// Create a range filter for the last 30 days
	dateFilter := elastic.NewRangeQuery("crawl_date").
		Gte(now.AddDate(0, 0, -30).Format(time.RFC3339)).
		Lte(now.Format(time.RFC3339))

	// Combine with a bool query
	boolQuery := elastic.NewBoolQuery().
		Must(query).
		Filter(dateFilter)

	// Create search service
	searchResult, err := client.Search().
		Index(IndexName).
		Query(boolQuery).
		From(from).
		Size(size).
		Sort("crawl_date", false). // Sort by date (newest first)
		Do(context.Background())

	if err != nil {
		log.Printf("Error getting recent recipes: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process search results
	var recipes []map[string]interface{}
	for _, hit := range searchResult.Hits.Hits {
		var recipe map[string]interface{}
		if err := json.Unmarshal(hit.Source, &recipe); err != nil {
			log.Printf("Error unmarshaling recipe: %s", err)
			continue
		}

		// Add the ID to the recipe
		recipe["id"] = hit.Id
		recipes = append(recipes, recipe)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}
