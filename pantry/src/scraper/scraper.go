package scraper

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// RecipeSite defines the structure for recipe site configurations
type RecipeSite struct {
	Domain      string
	URLPatterns []string
	Selectors   SiteSelectors
}

// SiteSelectors holds CSS selectors for different recipe sites
type SiteSelectors struct {
	RecipeTitle       string
	RecipeDescription string
	RecipeIngredients string
	RecipeInstructions string
	RecipeTime        string
	RecipeServings    string
	RecipeLinks       []string
}

// Scraper for each website
type Scraper struct {
	url  string
	doc  *goquery.Document
	site *RecipeSite
}

var recipeSites = []RecipeSite{
	{
		Domain: "pinchofyum.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-content p:first-of-type",
			RecipeIngredients: ".recipe-ingredients li, .wp-block-recipe-card-ingredients li",
			RecipeInstructions: ".recipe-instructions li, .wp-block-recipe-card-instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "minimalistbaker.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-summary",
			RecipeIngredients: ".recipe-ingredients li, ul.ingredients li",
			RecipeInstructions: ".recipe-instructions li, ol.instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .total-time",
			RecipeServings:    ".recipe-servings, .yield",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "cookieandkate.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-summary",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li",
			RecipeInstructions: ".recipe-instructions li, .instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "loveandlemons.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-summary",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li",
			RecipeInstructions: ".recipe-instructions li, .instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "smittenkitchen.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
			"/blog/",
			"/20", // Year-based URLs like /2025/06/recipe-name/
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title, h1, .post-title",
			RecipeDescription: ".recipe-description, .entry-summary, .entry-content p:first-of-type",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li, .entry-content p:contains('cup'), .entry-content p:contains('tablespoon'), .entry-content p:contains('teaspoon')",
			RecipeInstructions: ".recipe-instructions li, .instructions li, .entry-content p:contains('mix'), .entry-content p:contains('combine'), .entry-content p:contains('heat')",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time, .entry-content p:contains('minute'), .entry-content p:contains('hour')",
			RecipeServings:    ".recipe-servings, .servings, .entry-content p:contains('serve'), .entry-content p:contains('yield')",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']", "a[href*='/blog/']", "a[href*='/20']"},
		},
	},
	{
		Domain: "seriouseats.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.heading__title, h1.recipe-title",
			RecipeDescription: ".recipe-about, .recipe-description",
			RecipeIngredients: ".recipe-ingredients li, .structured-ingredients__list-item",
			RecipeInstructions: ".recipe-procedures li, .recipe-instructions li",
			RecipeTime:        ".recipe-time, .total-time, .active-time",
			RecipeServings:    ".recipe-yield, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "halfbakedharvest.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-summary",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li",
			RecipeInstructions: ".recipe-instructions li, .instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "101cookbooks.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-summary",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li",
			RecipeInstructions: ".recipe-instructions li, .instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "food52.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1, .recipe-title, .recipe-header-title, [data-testid='recipe-title']",
			RecipeDescription: ".recipe-summary, .recipe-description, .recipe-intro, .recipe-about, .intro",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li, .ingredient-list li, [data-testid='ingredient'], .recipe-ingredient",
			RecipeInstructions: ".recipe-instructions li, .instructions li, .direction-list li, [data-testid='instruction'], .recipe-instruction, .recipe-method li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time, .total-time, [data-testid='recipe-time']",
			RecipeServings:    ".recipe-servings, .servings, .yield, [data-testid='servings']",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "budgetbytes.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-summary",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li",
			RecipeInstructions: ".recipe-instructions li, .instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
	{
		Domain: "thewoksoflife.com",
		URLPatterns: []string{
			"/recipe/",
			"/recipes/",
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-summary",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li",
			RecipeInstructions: ".recipe-instructions li, .instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']"},
		},
	},
}

// getSiteConfig returns the recipe site configuration for a given URL
func getSiteConfig(u string) *RecipeSite {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil
	}

	hostname := strings.ToLower(parsedURL.Hostname())
	hostname = strings.TrimPrefix(hostname, "www.")

	for _, site := range recipeSites {
		if strings.Contains(hostname, site.Domain) {
			return &site
		}
	}
	return nil
}

// isRecipeURL checks if the URL matches recipe patterns for the site
func isRecipeURL(u string, site *RecipeSite) bool {
	if site == nil {
		return false
	}

	for _, pattern := range site.URLPatterns {
		if strings.Contains(u, pattern) {
			return true
		}
	}
	return false
}

// NewScraper builds a new scraper for the website with retries and better error handling
func NewScraper(u string) *Scraper {
	if !strings.HasPrefix(u, "http") {
		return nil
	}

	// Add rate limiting and retries
	client := &http.Client{
		Timeout: 45 * time.Second,
	}

	// Try up to 3 times with exponential backoff
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			log.Printf("Failed to create request for %s (attempt %d): %v", u, attempt, err)
			if attempt == 3 {
				return nil
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Set more realistic headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		response, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to fetch %s (attempt %d): %v", u, attempt, err)
			if attempt == 3 {
				return nil
			}
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}
		defer response.Body.Close()

		// Handle various status codes more gracefully
		if response.StatusCode == 404 {
			log.Printf("Page not found (404) for %s", u)
			return nil
		} else if response.StatusCode == 403 {
			log.Printf("Access forbidden (403) for %s", u)
			return nil
		} else if response.StatusCode >= 500 {
			log.Printf("Server error (%d) for %s (attempt %d)", response.StatusCode, u, attempt)
			if attempt == 3 {
				return nil
			}
			time.Sleep(time.Duration(attempt) * 3 * time.Second)
			continue
		} else if response.StatusCode != 200 {
			log.Printf("Non-200 status code %d for %s", response.StatusCode, u)
			return nil
		}

		// Successfully fetched page

		// Read response body and handle compression
		var reader io.Reader = response.Body
		
		// Check if response is gzipped and decompress if needed
		if response.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(response.Body)
			if err != nil {
				fmt.Printf("Failed to create gzip reader for %s (attempt %d): %v\n", u, attempt, err)
				if attempt == 3 {
					return nil
				}
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			defer gzipReader.Close()
			reader = gzipReader
		}
		
		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			fmt.Printf("Failed to read response body for %s (attempt %d): %v\n", u, attempt, err)
			if attempt == 3 {
				return nil
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Response body read successfully

		d, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
		if err != nil {
			log.Printf("Failed to parse HTML for %s (attempt %d): %v", u, attempt, err)
			if attempt == 3 {
				return nil
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		site := getSiteConfig(u)

		// HTML parsed successfully

		return &Scraper{
			url:  u,
			doc:  d,
			site: site,
		}
	}

	return nil
}

// Body returns a string with the body of the page
func (s *Scraper) Body() string {
	if s.site != nil && isRecipeURL(s.url, s.site) {
		return s.getStructuredRecipeData()
	}

	body := s.doc.Find("body").Text()
	body = strings.TrimSpace(body)

	return body
}

// getStructuredRecipeData extracts structured recipe data using site-specific selectors
func (s *Scraper) getStructuredRecipeData() string {
	// First try to extract JSON-LD structured data
	if jsonData := s.extractJSONLD(); jsonData != "" {
		return jsonData
	}

	var recipeData strings.Builder

	if s.site == nil {
		return s.doc.Find("body").Text()
	}

	selectors := s.site.Selectors

	// Extract recipe title
	if title := s.doc.Find(selectors.RecipeTitle).First().Text(); title != "" {
		recipeData.WriteString("TITLE: " + strings.TrimSpace(title) + "\n\n")
	}

	// Extract recipe description
	if description := s.doc.Find(selectors.RecipeDescription).First().Text(); description != "" {
		recipeData.WriteString("DESCRIPTION: " + strings.TrimSpace(description) + "\n\n")
	}

	// Extract ingredients
	recipeData.WriteString("INGREDIENTS:\n")
	s.doc.Find(selectors.RecipeIngredients).Each(func(i int, sel *goquery.Selection) {
		ingredient := strings.TrimSpace(sel.Text())
		if ingredient != "" {
			recipeData.WriteString("- " + ingredient + "\n")
		}
	})
	recipeData.WriteString("\n")

	// Extract instructions
	recipeData.WriteString("INSTRUCTIONS:\n")
	s.doc.Find(selectors.RecipeInstructions).Each(func(i int, sel *goquery.Selection) {
		instruction := strings.TrimSpace(sel.Text())
		if instruction != "" {
			recipeData.WriteString(fmt.Sprintf("%d. %s\n", i+1, instruction))
		}
	})
	recipeData.WriteString("\n")

	// Extract timing information
	if time := s.doc.Find(selectors.RecipeTime).First().Text(); time != "" {
		recipeData.WriteString("TIME: " + strings.TrimSpace(time) + "\n")
	}

	// Extract servings
	if servings := s.doc.Find(selectors.RecipeServings).First().Text(); servings != "" {
		recipeData.WriteString("SERVINGS: " + strings.TrimSpace(servings) + "\n")
	}

	result := recipeData.String()
	if result == "" {
		return s.doc.Find("body").Text()
	}

	return result
}

// extractJSONLD extracts recipe data from JSON-LD structured data
func (s *Scraper) extractJSONLD() string {
	var recipeData strings.Builder

	s.doc.Find("script[type='application/ld+json']").Each(func(i int, sel *goquery.Selection) {
		jsonText := sel.Text()
		var data interface{}
		
		if err := json.Unmarshal([]byte(jsonText), &data); err != nil {
			return
		}

		if recipe := s.findRecipeInJSON(data); recipe != nil {
			if title, ok := recipe["name"].(string); ok && title != "" {
				recipeData.WriteString("TITLE: " + title + "\n\n")
			}

			if description, ok := recipe["description"].(string); ok && description != "" {
				recipeData.WriteString("DESCRIPTION: " + description + "\n\n")
			}

			// Extract ingredients
			if ingredients, ok := recipe["recipeIngredient"].([]interface{}); ok {
				recipeData.WriteString("INGREDIENTS:\n")
				for _, ing := range ingredients {
					if ingredient, ok := ing.(string); ok {
						recipeData.WriteString("- " + ingredient + "\n")
					}
				}
				recipeData.WriteString("\n")
			}

			// Extract instructions
			if instructions, ok := recipe["recipeInstructions"].([]interface{}); ok {
				recipeData.WriteString("INSTRUCTIONS:\n")
				for idx, inst := range instructions {
					var instruction string
					if instObj, ok := inst.(map[string]interface{}); ok {
						if text, ok := instObj["text"].(string); ok {
							instruction = text
						}
					} else if instStr, ok := inst.(string); ok {
						instruction = instStr
					}
					if instruction != "" {
						recipeData.WriteString(fmt.Sprintf("%d. %s\n", idx+1, instruction))
					}
				}
				recipeData.WriteString("\n")
			}

			// Extract timing
			if prepTime, ok := recipe["prepTime"].(string); ok && prepTime != "" {
				recipeData.WriteString("PREP TIME: " + prepTime + "\n")
			}
			if cookTime, ok := recipe["cookTime"].(string); ok && cookTime != "" {
				recipeData.WriteString("COOK TIME: " + cookTime + "\n")
			}
			if totalTime, ok := recipe["totalTime"].(string); ok && totalTime != "" {
				recipeData.WriteString("TOTAL TIME: " + totalTime + "\n")
			}

			// Extract servings/yield
			if yield, ok := recipe["recipeYield"].(string); ok && yield != "" {
				recipeData.WriteString("SERVINGS: " + yield + "\n")
			} else if yieldArr, ok := recipe["recipeYield"].([]interface{}); ok && len(yieldArr) > 0 {
				if yieldStr, ok := yieldArr[0].(string); ok {
					recipeData.WriteString("SERVINGS: " + yieldStr + "\n")
				}
			}
		}
	})

	return recipeData.String()
}

// findRecipeInJSON recursively searches for Recipe schema in JSON-LD data
func (s *Scraper) findRecipeInJSON(data interface{}) map[string]interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		if typeField, ok := v["@type"].(string); ok && strings.ToLower(typeField) == "recipe" {
			return v
		}
		// Search in nested objects
		for _, value := range v {
			if recipe := s.findRecipeInJSON(value); recipe != nil {
				return recipe
			}
		}
	case []interface{}:
		// Search in array elements
		for _, item := range v {
			if recipe := s.findRecipeInJSON(item); recipe != nil {
				return recipe
			}
		}
	}
	return nil
}

func (s *Scraper) buildLinks(href string) string {
	var link string

	if strings.HasPrefix(href, "/") {
		link = strings.Join([]string{s.url, href}, "")
	} else {
		link = href
	}

	link = strings.TrimRight(link, "/")
	link = strings.TrimRight(link, ":")

	return link
}

// Links returns an array with all the links from the website
func (s *Scraper) Links() []string {
	links := make([]string, 0)

	if s.site != nil {
		// Use site-specific link selectors
		for _, selector := range s.site.Selectors.RecipeLinks {
			s.doc.Find(selector).Each(func(index int, item *goquery.Selection) {
				href, exists := item.Attr("href")
				if !exists {
					return
				}

				if s.isValidLink(href) {
					link := s.buildLinks(href)
					if link != "" && s.isRecipeLink(link) {
						links = append(links, link)
					}
				}
			})
		}
	} else {
		// Fallback to generic recipe link patterns
		s.doc.Find("body a").Each(func(index int, item *goquery.Selection) {
			href, exists := item.Attr("href")
			if !exists {
				return
			}

			if s.isValidLink(href) && s.isGenericRecipeLink(href) {
				link := s.buildLinks(href)
				if link != "" {
					links = append(links, link)
				}
			}
		})
	}

	return s.removeDuplicateLinks(links)
}

// isValidLink checks if a link is valid for crawling
func (s *Scraper) isValidLink(href string) bool {
	return !strings.HasPrefix(href, "#") && 
		   !strings.HasPrefix(href, "javascript") && 
		   !strings.HasPrefix(href, "mailto:") &&
		   !strings.HasPrefix(href, "tel:") &&
		   href != "" &&
		   !strings.Contains(href, "twitter.com") &&
		   !strings.Contains(href, "facebook.com") &&
		   !strings.Contains(href, "instagram.com") &&
		   !strings.Contains(href, "pinterest.com")
}

// isRecipeLink checks if the link is likely a recipe based on site configuration
func (s *Scraper) isRecipeLink(link string) bool {
	if s.site == nil {
		return s.isGenericRecipeLink(link)
	}

	for _, pattern := range s.site.URLPatterns {
		if strings.Contains(link, pattern) {
			return true
		}
	}
	return false
}

// isGenericRecipeLink checks for common recipe URL patterns
func (s *Scraper) isGenericRecipeLink(href string) bool {
	recipePatterns := []string{
		"/recipe/", "/recipes/", "/cooking/", "/food/", "/dish/", "/meal/",
		"recipe-ideas", "quick-and-easy", "chicken", "tacos", "pasta", "soup",
		"dessert", "breakfast", "lunch", "dinner", "appetizer", "snack",
		"vegetarian", "vegan", "healthy", "easy",
	}

	href = strings.ToLower(href)
	for _, pattern := range recipePatterns {
		if strings.Contains(href, pattern) {
			return true
		}
	}
	return false
}

// removeDuplicateLinks removes duplicate URLs from the links slice
func (s *Scraper) removeDuplicateLinks(links []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, link := range links {
		if !seen[link] {
			seen[link] = true
			result = append(result, link)
		}
	}

	return result
}

// MetaDataInformation returns the title and description from the page
func (s *Scraper) MetaDataInformation() (string, string) {
	var title, description string

	// Try to get recipe-specific title if this is a recipe page
	if s.site != nil && isRecipeURL(s.url, s.site) {
		if recipeTitle := s.doc.Find(s.site.Selectors.RecipeTitle).First().Text(); recipeTitle != "" {
			title = strings.TrimSpace(recipeTitle)
		}
		if recipeDesc := s.doc.Find(s.site.Selectors.RecipeDescription).First().Text(); recipeDesc != "" {
			description = strings.TrimSpace(recipeDesc)
		}
	}

	// Fallback to standard title if no recipe title found
	if title == "" {
		title = s.doc.Find("title").Contents().Text()
		title = strings.TrimSpace(title)
	}

	// Fallback to meta description if no recipe description found
	if description == "" {
		s.doc.Find("meta").Each(func(index int, item *goquery.Selection) {
			if item.AttrOr("name", "") == "description" || 
			   item.AttrOr("property", "") == "og:description" ||
			   item.AttrOr("property", "") == "og:title" {
				if content := item.AttrOr("content", ""); content != "" {
					if description == "" || item.AttrOr("name", "") == "description" {
						description = content
					}
				}
			}
		})
		description = strings.TrimSpace(description)
	}

	// Clean up title and description
	title = s.cleanText(title)
	description = s.cleanText(description)

	return title, description
}

// cleanText removes extra whitespace and common unwanted characters
func (s *Scraper) cleanText(text string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	
	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)
	
	return text
}

// GetRecipeData extracts comprehensive recipe data from the page
func (s *Scraper) GetRecipeData() map[string]string {
	data := make(map[string]string)

	// Get title and description
	title, description := s.MetaDataInformation()
	data["title"] = title
	data["description"] = description

	// First try JSON-LD structured data
	jsonRecipe := s.extractJSONLDRecipe()
	if jsonRecipe != nil {
		if name, ok := jsonRecipe["name"].(string); ok {
			data["name"] = name
		}
		if desc, ok := jsonRecipe["description"].(string); ok && data["description"] == "" {
			data["description"] = desc
		}
		if image, ok := jsonRecipe["image"].(string); ok {
			data["image"] = image
		} else if imageArr, ok := jsonRecipe["image"].([]interface{}); ok && len(imageArr) > 0 {
			if imageStr, ok := imageArr[0].(string); ok {
				data["image"] = imageStr
			}
		}
		if prepTime, ok := jsonRecipe["prepTime"].(string); ok {
			data["prep_time"] = prepTime
		}
		if cookTime, ok := jsonRecipe["cookTime"].(string); ok {
			data["cook_time"] = cookTime
		}
		if totalTime, ok := jsonRecipe["totalTime"].(string); ok {
			data["total_time"] = totalTime
		}
		if calories, ok := jsonRecipe["calories"].(string); ok {
			data["calories"] = calories
		}
		if yield, ok := jsonRecipe["recipeYield"].(string); ok {
			data["servings"] = yield
		} else if yieldArr, ok := jsonRecipe["recipeYield"].([]interface{}); ok && len(yieldArr) > 0 {
			if yieldStr, ok := yieldArr[0].(string); ok {
				data["servings"] = yieldStr
			}
		}

		// Extract ingredients
		if ingredients, ok := jsonRecipe["recipeIngredient"].([]interface{}); ok {
			var ingredientsList []string
			for _, ing := range ingredients {
				if ingredient, ok := ing.(string); ok {
					ingredientsList = append(ingredientsList, strings.TrimSpace(ingredient))
				}
			}
			data["ingredients"] = strings.Join(ingredientsList, ";")
		}

		// Extract instructions
		if instructions, ok := jsonRecipe["recipeInstructions"].([]interface{}); ok {
			var instructionsList []string
			for _, inst := range instructions {
				var instruction string
				if instObj, ok := inst.(map[string]interface{}); ok {
					if text, ok := instObj["text"].(string); ok {
						instruction = text
					}
				} else if instStr, ok := inst.(string); ok {
					instruction = instStr
				}
				if instruction != "" {
					instructionsList = append(instructionsList, strings.TrimSpace(instruction))
				}
			}
			data["instructions"] = strings.Join(instructionsList, ";")
		}
	}

	// If JSON-LD didn't provide data, try site-specific selectors
	if s.site != nil && isRecipeURL(s.url, s.site) {
		selectors := s.site.Selectors

		if data["name"] == "" {
			if name := s.doc.Find(selectors.RecipeTitle).First().Text(); name != "" {
				data["name"] = strings.TrimSpace(name)
			}
		}

		if data["description"] == "" {
			if desc := s.doc.Find(selectors.RecipeDescription).First().Text(); desc != "" {
				data["description"] = strings.TrimSpace(desc)
			}
		}

		if data["ingredients"] == "" {
			var ingredients []string
			s.doc.Find(selectors.RecipeIngredients).Each(func(i int, sel *goquery.Selection) {
				ingredient := strings.TrimSpace(sel.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
			if len(ingredients) > 0 {
				data["ingredients"] = strings.Join(ingredients, ";")
			}
		}

		if data["instructions"] == "" {
			var instructions []string
			s.doc.Find(selectors.RecipeInstructions).Each(func(i int, sel *goquery.Selection) {
				instruction := strings.TrimSpace(sel.Text())
				if instruction != "" {
					instructions = append(instructions, instruction)
				}
			})
			if len(instructions) > 0 {
				data["instructions"] = strings.Join(instructions, ";")
			}
		}

		if data["prep_time"] == "" || data["cook_time"] == "" || data["total_time"] == "" {
			if timeText := s.doc.Find(selectors.RecipeTime).First().Text(); timeText != "" {
				data["total_time"] = strings.TrimSpace(timeText)
			}
		}

		if data["servings"] == "" {
			if servings := s.doc.Find(selectors.RecipeServings).First().Text(); servings != "" {
				data["servings"] = strings.TrimSpace(servings)
			}
		}
	}

	// Try to extract image from meta tags if not found
	if data["image"] == "" {
		s.doc.Find("meta").Each(func(index int, item *goquery.Selection) {
			if item.AttrOr("property", "") == "og:image" || item.AttrOr("name", "") == "twitter:image" {
				if content := item.AttrOr("content", ""); content != "" {
					data["image"] = content
				}
			}
		})
	}

	// Use title as name if name is still empty
	if data["name"] == "" {
		data["name"] = data["title"]
	}

	// If we still don't have ingredients or instructions, try narrative extraction
	// This is particularly useful for blog-style recipe sites like smittenkitchen.com
	if (data["ingredients"] == "" || data["instructions"] == "") {
		narrativeData := s.extractNarrativeRecipe()
		
		if data["ingredients"] == "" && narrativeData["ingredients"] != "" {
			data["ingredients"] = narrativeData["ingredients"]
		}
		
		if data["instructions"] == "" && narrativeData["instructions"] != "" {
			data["instructions"] = narrativeData["instructions"]
		}
	}

	return data
}

// extractJSONLDRecipe extracts the first recipe from JSON-LD data
func (s *Scraper) extractJSONLDRecipe() map[string]interface{} {
	var recipe map[string]interface{}

	s.doc.Find("script[type='application/ld+json']").Each(func(i int, sel *goquery.Selection) {
		if recipe != nil {
			return // Already found a recipe
		}

		jsonText := sel.Text()
		
		// Clean up the JSON text to handle common formatting issues
		jsonText = strings.TrimSpace(jsonText)
		// Try to fix newlines in string literals (common issue with food52.com)
		jsonText = strings.ReplaceAll(jsonText, "\n", "\\n")
		jsonText = strings.ReplaceAll(jsonText, "\r", "\\r")
		jsonText = strings.ReplaceAll(jsonText, "\t", "\\t")
		
		var data interface{}
		
		if err := json.Unmarshal([]byte(jsonText), &data); err != nil {
			// Try once more with a simpler approach - remove all control characters
			jsonTextSimple := regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(sel.Text(), " ")
			if err2 := json.Unmarshal([]byte(jsonTextSimple), &data); err2 != nil {
				return
			}
		}

		recipe = s.findRecipeInJSON(data)
	})

	return recipe
}

// extractNarrativeRecipe attempts to extract recipe data from narrative text for sites like smittenkitchen.com
func (s *Scraper) extractNarrativeRecipe() map[string]string {
	data := make(map[string]string)
	bodyText := s.doc.Find(".entry-content, .post-content, .content").Text()
	
	// Try to extract ingredients from patterns
	lines := strings.Split(bodyText, "\n")
	var ingredients []string
	var instructions []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Look for ingredient patterns (measurements + ingredients)
		if strings.Contains(line, "cup") || strings.Contains(line, "tablespoon") || 
		   strings.Contains(line, "teaspoon") || strings.Contains(line, "pound") ||
		   strings.Contains(line, "ounce") || strings.Contains(line, "gram") {
			// This looks like an ingredient
			if len(line) < 200 { // Reasonable ingredient length
				ingredients = append(ingredients, line)
			}
		}
		
		// Look for instruction patterns
		if strings.Contains(line, "mix") || strings.Contains(line, "combine") ||
		   strings.Contains(line, "heat") || strings.Contains(line, "bake") ||
		   strings.Contains(line, "cook") || strings.Contains(line, "add") ||
		   strings.Contains(line, "stir") || strings.Contains(line, "pour") {
			// This looks like an instruction
			if len(line) > 20 && len(line) < 500 { // Reasonable instruction length
				instructions = append(instructions, line)
			}
		}
	}
	
	if len(ingredients) > 0 {
		data["ingredients"] = strings.Join(ingredients, ";")
	}
	
	if len(instructions) > 0 {
		data["instructions"] = strings.Join(instructions, ";")
	}
	
	return data
}