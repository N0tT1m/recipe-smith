package structs

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Page is the main struct for storing recipe data
type Page struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Body         string    `json:"body"`
	URL          string    `json:"url"`
	Image        string    `json:"image"`
	Name         string    `json:"name"`
	PrepTime     string    `json:"prep_time"`
	CookTime     string    `json:"cook_time"`
	TotalTime    string    `json:"total_time"`
	Calories     string    `json:"calories"`
	Servings     string    `json:"servings"`
	Ingredients  string    `json:"ingredients"`
	Instructions string    `json:"instructions"`
	SourceSite   string    `json:"source_site"`
	CrawlDate    time.Time `json:"crawl_date"`
	Categories   string    `json:"categories,omitempty"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SearchResult represents a search result
type SearchResult struct {
	TotalHits int64  `json:"total_hits"`
	Pages     []Page `json:"pages"`
}

// RecipeIngredient represents a single ingredient with its components
type RecipeIngredient struct {
	Original   string `json:"original"`
	Quantity   string `json:"quantity,omitempty"`
	Unit       string `json:"unit,omitempty"`
	Ingredient string `json:"ingredient"`
	Notes      string `json:"notes,omitempty"`
}

// ParsedRecipe represents a recipe with parsed ingredients and instructions
type ParsedRecipe struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	URL          string             `json:"url"`
	Image        string             `json:"image"`
	Name         string             `json:"name"`
	PrepTime     string             `json:"prep_time"`
	CookTime     string             `json:"cook_time"`
	TotalTime    string             `json:"total_time"`
	Calories     string             `json:"calories"`
	Servings     string             `json:"servings"`
	Ingredients  []RecipeIngredient `json:"ingredients"`
	Instructions []string           `json:"instructions"`
	SourceSite   string             `json:"source_site"`
	CrawlDate    time.Time          `json:"crawl_date"`
	Categories   []string           `json:"categories,omitempty"`
}

// RecipeStats represents aggregated statistics about crawled recipes
type RecipeStats struct {
	TotalRecipes         int             `json:"total_recipes"`
	RecipesBySite        map[string]int  `json:"recipes_by_site"`
	AvgIngredientsCount  float64         `json:"avg_ingredients_count"`
	AvgInstructionsCount float64         `json:"avg_instructions_count"`
	MostCommonCategories []CategoryCount `json:"most_common_categories"`
	LastCrawled          time.Time       `json:"last_crawled"`
}

// CategoryCount represents a category and its count
type CategoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// RecipeFilter represents filtering options for recipe search
type RecipeFilter struct {
	Ingredients []string `json:"ingredients,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	PrepTime    int      `json:"prep_time,omitempty"`  // Max prep time in minutes
	CookTime    int      `json:"cook_time,omitempty"`  // Max cook time in minutes
	TotalTime   int      `json:"total_time,omitempty"` // Max total time in minutes
	MaxCalories int      `json:"max_calories,omitempty"`
	MinCalories int      `json:"min_calories,omitempty"`
	Source      string   `json:"source,omitempty"` // Source website
}

// CrawlStatus represents the status of a recipe crawl operation
type CrawlStatus struct {
	StartTime      time.Time `json:"start_time"`
	TotalProcessed int       `json:"total_processed"`
	TotalSuccess   int       `json:"total_success"`
	TotalFailed    int       `json:"total_failed"`
	CurrentURL     string    `json:"current_url,omitempty"`
	IsRunning      bool      `json:"is_running"`
	LastError      string    `json:"last_error,omitempty"`
}

// SimilarRecipe represents a recipe that is similar to another
type SimilarRecipe struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	Image             string  `json:"image"`
	SimilarityScore   float64 `json:"similarity_score"`
	SharedIngredients int     `json:"shared_ingredients"`
	TotalIngredients  int     `json:"total_ingredients"`
	IngredientOverlap float64 `json:"ingredient_overlap"` // Percentage of shared ingredients
}

// RecipeRecommendation represents a recipe recommendation
type RecipeRecommendation struct {
	Recipe          Page     `json:"recipe"`
	ReasonForRec    string   `json:"reason_for_recommendation"`
	MatchingTags    []string `json:"matching_tags,omitempty"`
	PopularityScore float64  `json:"popularity_score"`
}

// ParseIngredients parses a semicolon-separated ingredient string into structured ingredients
func ParseIngredients(ingredients string) []RecipeIngredient {
	if ingredients == "" {
		return []RecipeIngredient{}
	}

	ingredientsList := strings.Split(ingredients, ";")
	result := make([]RecipeIngredient, 0, len(ingredientsList))

	for _, ing := range ingredientsList {
		ing = strings.TrimSpace(ing)
		if ing == "" {
			continue
		}

		// Create a basic ingredient structure
		parsed := RecipeIngredient{
			Original:   ing,
			Ingredient: ing,
		}

		// Try to extract quantity
		if quantityRegex := regexp.MustCompile(`^([\d\s./]+)`); quantityRegex.MatchString(ing) {
			matches := quantityRegex.FindStringSubmatch(ing)
			if len(matches) > 1 {
				parsed.Quantity = strings.TrimSpace(matches[1])
				parsed.Ingredient = strings.TrimSpace(ing[len(matches[1]):])
			}
		}

		// Try to extract unit
		unitPattern := `^\s*([\d\s./]+)?\s*(cup|cups|tbsp|tsp|tablespoon|tablespoons|teaspoon|teaspoons|oz|ounce|ounces|lb|pound|pounds|g|gram|grams|kg|kilogram|kilograms|ml|milliliter|milliliters|l|liter|liters|pinch|dash|can|cans|clove|cloves)\s*`
		if unitRegex := regexp.MustCompile(unitPattern); unitRegex.MatchString(parsed.Ingredient) {
			matches := unitRegex.FindStringSubmatch(parsed.Ingredient)
			if len(matches) > 2 {
				if parsed.Quantity == "" && matches[1] != "" {
					parsed.Quantity = strings.TrimSpace(matches[1])
				}
				parsed.Unit = strings.TrimSpace(matches[2])
				parsed.Ingredient = strings.TrimSpace(parsed.Ingredient[len(matches[0]):])
			}
		}

		// Check for notes in parentheses
		if notesRegex := regexp.MustCompile(`(.*?)\s*\((.*?)\)\s*(.*)`); notesRegex.MatchString(parsed.Ingredient) {
			matches := notesRegex.FindStringSubmatch(parsed.Ingredient)
			if len(matches) > 3 {
				parsed.Notes = strings.TrimSpace(matches[2])
				parsed.Ingredient = strings.TrimSpace(matches[1] + " " + matches[3])
			}
		}

		// Do some final cleanup
		parsed.Ingredient = strings.TrimSpace(parsed.Ingredient)

		// Remove trailing commas, semicolons, etc.
		parsed.Ingredient = regexp.MustCompile(`[,;.]+$`).ReplaceAllString(parsed.Ingredient, "")

		result = append(result, parsed)
	}

	return result
}

// ParseInstructions parses a semicolon-separated instruction string into a slice of strings
func ParseInstructions(instructions string) []string {
	if instructions == "" {
		return []string{}
	}

	instructionsList := strings.Split(instructions, ";")
	result := make([]string, 0, len(instructionsList))

	for _, instruction := range instructionsList {
		instruction = strings.TrimSpace(instruction)
		if instruction == "" {
			continue
		}

		// Remove step numbers (e.g., "Step 1:", "1.", etc.)
		instruction = regexp.MustCompile(`^(?:Step\s*)?(?:\d+\.?\s*:?\s*)`).ReplaceAllString(instruction, "")

		// Remove common JSON-LD artifacts that might have been incorrectly parsed
		instruction = regexp.MustCompile(`@type|HowToStep|text|[{}[\]"`).ReplaceAllString(instruction, "")

		instruction = strings.TrimSpace(instruction)
		if instruction != "" {
			result = append(result, instruction)
		}
	}

	return result
}

// ParseCategories parses a semicolon-separated category string into a slice of strings
func ParseCategories(categories string) []string {
	if categories == "" {
		return []string{}
	}

	categoriesList := strings.Split(categories, ";")
	result := make([]string, 0, len(categoriesList))

	for _, category := range categoriesList {
		category = strings.TrimSpace(category)
		if category == "" {
			continue
		}
		result = append(result, category)
	}

	return result
}

// ToPageStruct converts a ParsedRecipe back to a Page struct for storage
func (pr *ParsedRecipe) ToPageStruct() Page {
	// Convert ingredients back to semicolon-separated string
	ingredientsStrings := make([]string, len(pr.Ingredients))
	for i, ing := range pr.Ingredients {
		ingredientsStrings[i] = ing.Original
	}

	// Convert instructions back to semicolon-separated string
	instructionsString := strings.Join(pr.Instructions, ";")

	// Convert categories back to semicolon-separated string
	categoriesString := strings.Join(pr.Categories, ";")

	return Page{
		ID:           pr.ID,
		Title:        pr.Title,
		Description:  pr.Description,
		Body:         "", // Body is not typically stored in ParsedRecipe
		URL:          pr.URL,
		Image:        pr.Image,
		Name:         pr.Name,
		PrepTime:     pr.PrepTime,
		CookTime:     pr.CookTime,
		TotalTime:    pr.TotalTime,
		Calories:     pr.Calories,
		Servings:     pr.Servings,
		Ingredients:  strings.Join(ingredientsStrings, ";"),
		Instructions: instructionsString,
		SourceSite:   pr.SourceSite,
		CrawlDate:    pr.CrawlDate,
		Categories:   categoriesString,
	}
}

// FromPageStruct converts a Page struct to a ParsedRecipe
func FromPageStruct(page Page) ParsedRecipe {
	return ParsedRecipe{
		ID:           page.ID,
		Title:        page.Title,
		Description:  page.Description,
		URL:          page.URL,
		Image:        page.Image,
		Name:         page.Name,
		PrepTime:     page.PrepTime,
		CookTime:     page.CookTime,
		TotalTime:    page.TotalTime,
		Calories:     page.Calories,
		Servings:     page.Servings,
		Ingredients:  ParseIngredients(page.Ingredients),
		Instructions: ParseInstructions(page.Instructions),
		SourceSite:   page.SourceSite,
		CrawlDate:    page.CrawlDate,
		Categories:   ParseCategories(page.Categories),
	}
}

// EstimateTimeInMinutes converts time strings like "1 hr 30 min" to minutes
func EstimateTimeInMinutes(timeStr string) int {
	if timeStr == "" {
		return 0
	}

	total := 0

	// Extract hours
	hourRegex := regexp.MustCompile(`(\d+)\s*(?:hr|hour|h)`)
	if matches := hourRegex.FindStringSubmatch(timeStr); len(matches) > 1 {
		hours := 0
		fmt.Sscanf(matches[1], "%d", &hours)
		total += hours * 60
	}

	// Extract minutes
	minRegex := regexp.MustCompile(`(\d+)\s*(?:min|minute|m)`)
	if matches := minRegex.FindStringSubmatch(timeStr); len(matches) > 1 {
		minutes := 0
		fmt.Sscanf(matches[1], "%d", &minutes)
		total += minutes
	}

	return total
}
