package scraper

import (
	"encoding/json"
	"fmt"
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
<<<<<<< HEAD
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
		},
		Selectors: SiteSelectors{
			RecipeTitle:       "h1.entry-title, h1.recipe-title",
			RecipeDescription: ".recipe-description, .entry-summary",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li",
			RecipeInstructions: ".recipe-instructions li, .instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
			RecipeLinks:       []string{"a[href*='/recipe/']", "a[href*='/recipes/']", "a[href*='/blog/']"},
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
			RecipeTitle:       "h1.recipe-header-title, h1.recipe-title",
			RecipeDescription: ".recipe-summary, .recipe-description",
			RecipeIngredients: ".recipe-ingredients li, .ingredients li",
			RecipeInstructions: ".recipe-instructions li, .instructions li",
			RecipeTime:        ".recipe-time, .prep-time, .cook-time",
			RecipeServings:    ".recipe-servings, .servings",
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
	// Remove www. prefix if present
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
=======
	site string
>>>>>>> 6d26aa9a75e420916d7aa3fa0179960e9d47d47a
}

// NewScraper builds a new scraper for the website
func NewScraper(u string) *Scraper {
	if !strings.HasPrefix(u, "http") {
		return nil
	}

	// Create a client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Set user agent to mimic a browser
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching URL %s: %v", u, err)
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("Failed to fetch %s, status code: %d", u, response.StatusCode)
		return nil
	}

<<<<<<< HEAD
	site := getSiteConfig(u)

	return &Scraper{
		url:  u,
		doc:  d,
		site: site,
=======
	d, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Printf("Error parsing document: %v", err)
		return nil
	}

	// Parse the URL to get the site name
	parsedURL, err := url.Parse(u)
	if err != nil {
		log.Printf("Error parsing URL: %v", err)
		return nil
	}

	log.Printf("Successfully created scraper for %s", u)
	return &Scraper{
		url:  u,
		doc:  d,
		site: parsedURL.Host,
>>>>>>> 6d26aa9a75e420916d7aa3fa0179960e9d47d47a
	}
}

// Body returns a string with the body of the page
func (s *Scraper) Body() string {
	if s.site != nil && isRecipeURL(s.url, s.site) {
		return s.getStructuredRecipeData()
	}

	body := s.doc.Find("body").Text()
	// Remove leading/ending white spaces
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
		// Fallback to body text if no structured data found
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
		parsedURL, err := url.Parse(s.url)
		if err != nil {
			return ""
		}

		baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
		link = baseURL + href
	} else if !strings.HasPrefix(href, "http") {
		parsedURL, err := url.Parse(s.url)
		if err != nil {
			return ""
		}

		// Extract the base directory
		urlPath := parsedURL.Path
		lastSlash := strings.LastIndex(urlPath, "/")
		if lastSlash != -1 {
			urlPath = urlPath[:lastSlash+1]
		} else {
			urlPath = "/"
		}

		baseURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, urlPath)
		link = baseURL + href
	} else {
		link = href
	}

	link = strings.TrimRight(link, "/")
	link = strings.TrimRight(link, ":")

	return link
}

// Links returns a slice of links from the page
func (s *Scraper) Links() []string {
	links := make([]string, 0)
<<<<<<< HEAD

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
=======
	var link string

	// Map of popular recipe site patterns - expanded
	recipePatterns := map[string][]string{
		"www.delish.com":        {"/cooking/recipe-ideas", "/everyday-cooking/quick-and-easy/", "/recipe/", "/recipes/"},
		"www.allrecipes.com":    {"/recipe/", "/recipes/", "/article/", "/gallery/"},
		"www.foodnetwork.com":   {"/recipes/", "/recipe/", "/cuisine/", "/profiles/"},
		"www.epicurious.com":    {"/recipes/", "/recipe/", "/ingredient/", "/cuisine/"},
		"www.simplyrecipes.com": {"/recipes/", "/cuisine/", "/meal-type/", "/dietary-considerations/"},
		"thepioneerwoman.com":   {"/food-cooking/recipes/", "/food-cooking/meals/"},
		"www.bonappetit.com":    {"/recipe/", "/recipes/", "/story/", "/gallery/"},
		"www.marthastewart.com": {"/recipe/", "/recipes/", "/food/", "/dining/"},
		"www.seriouseats.com":   {"/recipes/", "/cuisine/", "/ingredients/", "/techniques/"},
		"cooking.nytimes.com":   {"/recipes/", "/collections/", "/guides/"},
	}

	// First collect all links
	allLinks := make([]string, 0)

	s.doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if !exists {
			return
		}

		// Skip anchors and javascript
		if strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript") {
			return
		}

		link = s.buildLinks(href)
		if link != "" {
			allLinks = append(allLinks, link)
>>>>>>> 6d26aa9a75e420916d7aa3fa0179960e9d47d47a
		}
	}
	return false
}

<<<<<<< HEAD
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
=======
	// Then filter by patterns for current site
	currentHost := s.site
	for domain, patterns := range recipePatterns {
		if strings.Contains(currentHost, domain) {
			// Check against patterns for this domain
			for _, link := range allLinks {
				for _, pattern := range patterns {
					if strings.Contains(link, pattern) {
						links = append(links, link)
						break
					}
				}
			}
			break
		}
	}

	// If no links matched patterns, include all links from same domain
	// This helps avoid getting stuck on sites with different URL patterns
	if len(links) == 0 {
		for _, link := range allLinks {
			parsedURL, err := url.Parse(link)
			if err != nil {
				continue
			}

			// Only include links from the same domain
			if strings.Contains(parsedURL.Host, currentHost) {
				links = append(links, link)
			}
		}
	}

	return removeDuplicateLinks(links)
}

// removeDuplicateLinks removes duplicate links from a slice
func removeDuplicateLinks(links []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, link := range links {
		if _, ok := seen[link]; !ok {
>>>>>>> 6d26aa9a75e420916d7aa3fa0179960e9d47d47a
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

// GetRecipeIngredients extracts recipe ingredients
func (s *Scraper) GetRecipeIngredients() string {
	var ingredients []string

	switch {
	case strings.Contains(s.site, "allrecipes.com"):
		// Check for ingredients in structured lists - most common in Allrecipes
		s.doc.Find(".ingredients-item-name, .ingredients-item, .checklist__item, .mntl-structured-ingredients__list-item").Each(func(i int, item *goquery.Selection) {
			ingredient := strings.TrimSpace(item.Text())
			if ingredient != "" {
				ingredients = append(ingredients, ingredient)
			}
		})

		// Try the newer structured ingredients format (as of 2025)
		if len(ingredients) == 0 {
			s.doc.Find(".mm-recipes-structured-ingredients__list-item p").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		}
	}

	// If still no ingredients, try schema.org approach
	if len(ingredients) == 0 {
		// Try schema.org ingredients
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeIngredient'], [itemtype='http://schema.org/Recipe'] [itemprop='ingredients'], [itemscope] [itemprop='recipeIngredient'], [itemscope] [itemprop='ingredients']").Each(func(i int, item *goquery.Selection) {
			ingredient := strings.TrimSpace(item.Text())
			if ingredient != "" {
				ingredients = append(ingredients, ingredient)
			}
		})
	}

	// Try LD+JSON schema for ingredients
	if len(ingredients) == 0 {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			if len(ingredients) > 0 {
				return
			}

			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") &&
				(strings.Contains(scriptContent, "recipeIngredient") ||
					strings.Contains(scriptContent, "ingredients")) {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check if this is a Recipe
					if jsonType, ok := jsonMap["@type"].(string); ok && (jsonType == "Recipe" || strings.Contains(jsonType, "Recipe")) {
						// Try to get ingredient information - first try recipeIngredient
						if ingList, ok := jsonMap["recipeIngredient"].([]interface{}); ok {
							for _, ing := range ingList {
								if ingStr, ok := ing.(string); ok {
									ingredients = append(ingredients, ingStr)
								}
							}
						} else if ingList, ok := jsonMap["ingredients"].([]interface{}); ok {
							// Then try ingredients
							for _, ing := range ingList {
								if ingStr, ok := ing.(string); ok {
									ingredients = append(ingredients, ingStr)
								}
							}
						}
					}

					// Handle nested recipe data
					if graphData, ok := jsonMap["@graph"].([]interface{}); ok && len(ingredients) == 0 {
						for _, itemData := range graphData {
							if itemMap, ok := itemData.(map[string]interface{}); ok {
								if itemType, ok := itemMap["@type"].(string); ok &&
									(itemType == "Recipe" || strings.Contains(itemType, "Recipe")) {
									if ingList, ok := itemMap["recipeIngredient"].([]interface{}); ok {
										for _, ing := range ingList {
											if ingStr, ok := ing.(string); ok {
												ingredients = append(ingredients, ingStr)
											}
										}
									} else if ingList, ok := itemMap["ingredients"].([]interface{}); ok {
										for _, ing := range ingList {
											if ingStr, ok := ing.(string); ok {
												ingredients = append(ingredients, ingStr)
											}
										}
									}
								}
							}
						}
					}
				} else {
					// Fallback to regex
					reIngArr := regexp.MustCompile(`"recipeIngredient"\s*:\s*\[(.*?)\]`)
					matches := reIngArr.FindStringSubmatch(scriptContent)
					if len(matches) >= 2 {
						ingredientListStr := matches[1]
						ingRe := regexp.MustCompile(`"([^"]+)"`)
						ingMatches := ingRe.FindAllStringSubmatch(ingredientListStr, -1)
						for _, match := range ingMatches {
							if len(match) >= 2 {
								ingredients = append(ingredients, match[1])
							}
						}
					}

					// Try old schema ingredients property
					if len(ingredients) == 0 {
						reIngOld := regexp.MustCompile(`"ingredients"\s*:\s*\[(.*?)\]`)
						matches := reIngOld.FindStringSubmatch(scriptContent)
						if len(matches) >= 2 {
							ingredientListStr := matches[1]
							ingRe := regexp.MustCompile(`"([^"]+)"`)
							ingMatches := ingRe.FindAllStringSubmatch(ingredientListStr, -1)
							for _, match := range ingMatches {
								if len(match) >= 2 {
									ingredients = append(ingredients, match[1])
								}
							}
						}
					}
				}
			}
		})
	}

	// Site-specific ingredient extraction for common recipe sites
	if len(ingredients) == 0 {
		switch {
		case strings.Contains(s.site, "delish.com"):
			s.doc.Find(".ingredient-item, .ingredients-item, .ingredient-list li").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		case strings.Contains(s.site, "allrecipes.com"):
			s.doc.Find(".ingredients-item-name, .ingredients-item, .checklist__item").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		case strings.Contains(s.site, "foodnetwork.com"):
			s.doc.Find(".o-Ingredients__a-Ingredient, .recipe-ingredients li, .ingredient-list li").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		case strings.Contains(s.site, "epicurious.com"):
			s.doc.Find(".ingredient, .ingredients-list li, .ingredient-item").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		case strings.Contains(s.site, "simplyrecipes.com"):
			s.doc.Find(".ingredient, .ingredients-list li, .ingredient-list li").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		}
	}

	// Generic fallback to common ingredient list patterns
	if len(ingredients) == 0 {
		s.doc.Find(".ingredients-item-name, .ingredient, .ingredient-list li, .ingredients-list li, .recipe-ingredients li").Each(func(i int, item *goquery.Selection) {
			ingredient := strings.TrimSpace(item.Text())
			if ingredient != "" {
				ingredients = append(ingredients, ingredient)
			}
		})
	}

	// More aggressive fallback - find any list in a section that appears to be ingredients
	if len(ingredients) == 0 {
		// Look for section headings
		s.doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, item *goquery.Selection) {
			headingText := strings.ToLower(strings.TrimSpace(item.Text()))
			if strings.Contains(headingText, "ingredient") {
				// Look for list items in the next sibling elements
				item.NextAll().EachWithBreak(func(j int, sibling *goquery.Selection) bool {
					// Stop if we hit another heading
					if sibling.Is("h1, h2, h3, h4, h5, h6") {
						return false
					}

					// Look for list items
					sibling.Find("li").Each(func(k int, listItem *goquery.Selection) {
						ingredient := strings.TrimSpace(listItem.Text())
						if ingredient != "" && len(ingredient) > 2 && !strings.Contains(strings.ToLower(ingredient), "instruction") {
							ingredients = append(ingredients, ingredient)
						}
					})

					return true
				})
			}
		})
	}

	// Clean up ingredients - remove excessive whitespace, duplicates
	uniqueIngredients := make(map[string]bool)
	var cleanedIngredients []string

	for _, ingredient := range ingredients {
		// Remove extra whitespace
		ingredient = regexp.MustCompile(`\s+`).ReplaceAllString(ingredient, " ")
		ingredient = strings.TrimSpace(ingredient)

		// Skip empty or very short ingredients
		if ingredient == "" || len(ingredient) < 2 {
			continue
		}

		// Check for duplicates
		if _, exists := uniqueIngredients[strings.ToLower(ingredient)]; !exists {
			uniqueIngredients[strings.ToLower(ingredient)] = true
			cleanedIngredients = append(cleanedIngredients, ingredient)
		}
	}

	// Join with semicolons as per the DB structure
	return strings.Join(cleanedIngredients, ";")
}

// GetRecipeInstructions extracts recipe instructions
func (s *Scraper) GetRecipeInstructions() string {
	var instructions []string

	switch {
	case strings.Contains(s.site, "allrecipes.com"):
		// Check for instructions in Allrecipes specific structure
		s.doc.Find(".step, .instructions-section-item, .recipe-directions__list--item, .step-item, .mntl-sc-block-group--LI").Each(func(i int, item *goquery.Selection) {
			// For newer Allrecipes structure, get the text from paragraph inside step
			instructionText := ""
			paragraphs := item.Find("p")
			if paragraphs.Length() > 0 {
				paragraphs.Each(func(j int, p *goquery.Selection) {
					text := strings.TrimSpace(p.Text())
					if text != "" {
						if instructionText != "" {
							instructionText += " "
						}
						instructionText += text
					}
				})
			} else {
				instructionText = strings.TrimSpace(item.Text())
			}

			if instructionText != "" {
				instructions = append(instructions, instructionText)
			}
		})
	}

	// If still no instructions, try schema.org approach
	if len(instructions) == 0 {
		// Try schema.org instructions
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeInstructions'], [itemscope] [itemprop='recipeInstructions']").Each(func(i int, item *goquery.Selection) {
			// Check if this is a container with multiple steps
			steps := item.Find("li, .step")
			if steps.Length() > 0 {
				steps.Each(func(j int, step *goquery.Selection) {
					instruction := strings.TrimSpace(step.Text())
					if instruction != "" {
						instructions = append(instructions, instruction)
					}
				})
			} else {
				// Single instruction text
				instruction := strings.TrimSpace(item.Text())
				if instruction != "" {
					instructions = append(instructions, instruction)
				}
			}
		})
	}

	// Try LD+JSON schema for instructions - complex version handling multiple formats
	if len(instructions) == 0 {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			if len(instructions) > 0 {
				return
			}

			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "recipeInstructions") {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check if this is a Recipe
					if jsonType, ok := jsonMap["@type"].(string); ok && (jsonType == "Recipe" || strings.Contains(jsonType, "Recipe")) {
						// Try to get instruction information - handle different formats
						if instList, ok := jsonMap["recipeInstructions"].([]interface{}); ok {
							for _, inst := range instList {
								// Format 1: Array of strings
								if instStr, ok := inst.(string); ok {
									instructions = append(instructions, instStr)
								} else if instObj, ok := inst.(map[string]interface{}); ok {
									// Format 2: Array of HowToStep objects
									if text, ok := instObj["text"].(string); ok {
										instructions = append(instructions, text)
									}
								}
							}
						} else if instStr, ok := jsonMap["recipeInstructions"].(string); ok {
							// Format 3: Single string with all instructions
							// Split by periods or newlines
							splits := regexp.MustCompile(`[.;]\s+`).Split(instStr, -1)
							for _, s := range splits {
								s = strings.TrimSpace(s)
								if s != "" {
									instructions = append(instructions, s)
								}
							}
						}
					}

					// Handle nested recipe data in @graph
					if graphData, ok := jsonMap["@graph"].([]interface{}); ok && len(instructions) == 0 {
						for _, itemData := range graphData {
							if itemMap, ok := itemData.(map[string]interface{}); ok {
								if itemType, ok := itemMap["@type"].(string); ok &&
									(itemType == "Recipe" || strings.Contains(itemType, "Recipe")) {
									if instList, ok := itemMap["recipeInstructions"].([]interface{}); ok {
										for _, inst := range instList {
											if instStr, ok := inst.(string); ok {
												instructions = append(instructions, instStr)
											} else if instObj, ok := inst.(map[string]interface{}); ok {
												if text, ok := instObj["text"].(string); ok {
													instructions = append(instructions, text)
												}
											}
										}
									}
								}
							}
						}
					}
				} else {
					// Multiple regex patterns to handle different JSON-LD formats for instructions

					// Try to match HowToStep format first (most common)
					stepFormat := regexp.MustCompile(`"recipeInstructions"\s*:\s*\[\s*\{\s*"@type"\s*:\s*"HowToStep",\s*"text"\s*:\s*"([^"]+)"`)
					stepMatches := stepFormat.FindAllStringSubmatch(scriptContent, -1)
					for _, match := range stepMatches {
						if len(match) >= 2 {
							instructions = append(instructions, match[1])
						}
					}

					// If no matches, try array of strings format
					if len(instructions) == 0 {
						stringFormat := regexp.MustCompile(`"recipeInstructions"\s*:\s*\[\s*"([^"]+)"`)
						stringMatches := stringFormat.FindAllStringSubmatch(scriptContent, -1)
						for _, match := range stringMatches {
							if len(match) >= 2 {
								instructions = append(instructions, match[1])
							}
						}
					}

					// If still no matches, try single string format
					if len(instructions) == 0 {
						singleFormat := regexp.MustCompile(`"recipeInstructions"\s*:\s*"([^"]+)"`)
						singleMatches := singleFormat.FindStringSubmatch(scriptContent)
						if len(singleMatches) >= 2 {
							// Split by periods or newlines
							splits := regexp.MustCompile(`[.;]\s+`).Split(singleMatches[1], -1)
							for _, s := range splits {
								s = strings.TrimSpace(s)
								if s != "" {
									instructions = append(instructions, s)
								}
							}
						}
					}
				}
			}
		})
	}

	// Site-specific extraction
	if len(instructions) == 0 {
		switch {
		case strings.Contains(s.site, "delish.com"):
			// Check for step-by-step instructions in ordered lists
			s.doc.Find("ol li").Each(func(i int, item *goquery.Selection) {
				instruction := strings.TrimSpace(item.Text())
				if instruction != "" && len(instruction) > 10 {
					instructions = append(instructions, instruction)
				}
			})

			// Look for direction elements
			if len(instructions) == 0 {
				s.doc.Find(".direction-item, .preparation-steps li, .preparation-step").Each(func(i int, item *goquery.Selection) {
					instruction := strings.TrimSpace(item.Text())
					if instruction != "" {
						instructions = append(instructions, instruction)
					}
				})
			}
		case strings.Contains(s.site, "allrecipes.com"):
			s.doc.Find(".step, .instructions-section-item, .recipe-directions__list--item, .step-item").Each(func(i int, item *goquery.Selection) {
				instruction := strings.TrimSpace(item.Text())
				if instruction != "" {
					instructions = append(instructions, instruction)
				}
			})
		case strings.Contains(s.site, "foodnetwork.com"):
			s.doc.Find(".o-Method__m-Step, .recipe-directions-list li, .direction-lists li").Each(func(i int, item *goquery.Selection) {
				instruction := strings.TrimSpace(item.Text())
				if instruction != "" {
					instructions = append(instructions, instruction)
				}
			})
		case strings.Contains(s.site, "epicurious.com"):
			s.doc.Find(".preparation-step, .preparation-steps li, .instruction-step").Each(func(i int, item *goquery.Selection) {
				instruction := strings.TrimSpace(item.Text())
				if instruction != "" {
					instructions = append(instructions, instruction)
				}
			})
		}
	}

	// Fallback to common instruction list patterns
	if len(instructions) == 0 {
		s.doc.Find(".instructions-item-name, .recipe-directions__list--item, .prep-steps li, .recipe-method-step, .recipe-instructions li").Each(func(i int, item *goquery.Selection) {
			instruction := strings.TrimSpace(item.Text())
			if instruction != "" {
				instructions = append(instructions, instruction)
			}
		})
	}

	// More aggressive fallback - find any list in a section that appears to be instructions
	if len(instructions) == 0 {
		// Look for section headings that indicate instructions
		instructionKeywords := []string{"instruction", "direction", "method", "preparation", "steps", "how to"}

		s.doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, item *goquery.Selection) {
			headingText := strings.ToLower(strings.TrimSpace(item.Text()))

			// Check if heading contains instruction keywords
			isInstructionHeading := false
			for _, keyword := range instructionKeywords {
				if strings.Contains(headingText, keyword) {
					isInstructionHeading = true
					break
				}
			}

			if isInstructionHeading {
				// Look for list items in the next sibling elements
				item.NextAll().EachWithBreak(func(j int, sibling *goquery.Selection) bool {
					// Stop if we hit another heading
					if sibling.Is("h1, h2, h3, h4, h5, h6") {
						return false
					}

					// Look for ordered list items
					sibling.Find("ol li").Each(func(k int, listItem *goquery.Selection) {
						instruction := strings.TrimSpace(listItem.Text())
						if instruction != "" && len(instruction) > 10 {
							instructions = append(instructions, instruction)
						}
					})

					// If no ordered list, look for any list items
					if len(instructions) == 0 {
						sibling.Find("li").Each(func(k int, listItem *goquery.Selection) {
							instruction := strings.TrimSpace(listItem.Text())
							if instruction != "" && len(instruction) > 10 {
								instructions = append(instructions, instruction)
							}
						})
					}

					// If still no instructions, look for paragraphs
					if len(instructions) == 0 {
						sibling.Find("p").Each(func(k int, para *goquery.Selection) {
							instruction := strings.TrimSpace(para.Text())
							if instruction != "" && len(instruction) > 20 {
								instructions = append(instructions, instruction)
							}
						})
					}

					return true
				})
			}
		})
	}

	// Last resort: Look for numbered paragraphs or paragraphs with step indicators
	if len(instructions) == 0 {
		var numberedInstructions []string

		s.doc.Find("p").Each(func(i int, item *goquery.Selection) {
			text := strings.TrimSpace(item.Text())

			// Check for numbered steps (e.g., "1.", "Step 1:", etc.)
			stepRegex := regexp.MustCompile(`^(?:Step\s*)?(\d+)[.:)]`)
			if stepRegex.MatchString(text) && len(text) > 15 {
				numberedInstructions = append(numberedInstructions, text)
			}
		})

		// If we found some numbered paragraphs, use them
		if len(numberedInstructions) > 0 {
			instructions = numberedInstructions
		}
	}

	// Clean up instructions
	var cleanedInstructions []string
	for _, instruction := range instructions {
		// Remove "Step X:" prefixes
		instruction = regexp.MustCompile(`^Step\s*\d+\s*:?\s*`).ReplaceAllString(instruction, "")

		// Remove excessive whitespace
		instruction = regexp.MustCompile(`\s+`).ReplaceAllString(instruction, " ")
		instruction = strings.TrimSpace(instruction)

		// Skip empty instructions or very short ones (likely not actually instructions)
		if instruction == "" || len(instruction) < 5 {
			continue
		}

		cleanedInstructions = append(cleanedInstructions, instruction)
	}

	// Join with semicolons as per the DB structure
	return strings.Join(cleanedInstructions, ";")
}

// GetRecipeData extracts all recipe data in one function
func (s *Scraper) GetRecipeData() map[string]string {
	recipeData := make(map[string]string)

	// Get basic metadata
	title, description := s.MetaDataInformation()
	recipeData["title"] = title
	recipeData["description"] = description

	// Get recipe-specific data
	recipeData["image"] = s.GetRecipeImage()
	recipeData["name"] = s.GetRecipeName()

	// If name wasn't found, use the title
	if recipeData["name"] == "" {
		recipeData["name"] = title
	}

	// Get prep time and total time, then calculate cook time
	prepTime, cookTime, totalTime := s.GetRecipeTime()
	recipeData["prep_time"] = prepTime
	recipeData["cook_time"] = cookTime
	recipeData["total_time"] = totalTime

	calories, servings := s.GetRecipeNutrition()
	recipeData["calories"] = calories
	recipeData["servings"] = servings

	recipeData["ingredients"] = s.GetRecipeIngredients()
	recipeData["instructions"] = s.GetRecipeInstructions()

	// Get categories
	categories := s.GetRecipeCategories()
	if len(categories) > 0 {
		recipeData["categories"] = strings.Join(categories, ";")
	}

	// Extract source site from URL
	parsedURL, err := url.Parse(s.url)
	if err == nil {
		recipeData["source_site"] = parsedURL.Host
	}

	recipeData["url"] = s.url

	// Enhanced validation and recovery for essential fields
	// If name is still empty, try finding it from h1 tags
	if recipeData["name"] == "" {
		s.doc.Find("h1").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				recipeData["name"] = strings.TrimSpace(item.Text())
			}
		})
	}

	// Last resort for name - use URL path
	if recipeData["name"] == "" && parsedURL != nil {
		path := parsedURL.Path
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			// Convert kebab-case to Title Case
			lastPart = strings.ReplaceAll(lastPart, "-", " ")
			lastPart = strings.ReplaceAll(lastPart, "_", " ")
			// Title case it
			words := strings.Fields(lastPart)
			for i, word := range words {
				if len(word) > 0 {
					r := []rune(word)
					r[0] = []rune(strings.ToUpper(string(r[0])))[0]
					words[i] = string(r)
				}
			}
			recipeData["name"] = strings.Join(words, " ")
		}
	}

	// Extract site-specific structured data
	extractSiteSpecificData(s, recipeData)

	// Try to extract ingredients and instructions from JSON if not already found
	if recipeData["ingredients"] == "" || recipeData["instructions"] == "" {
		tryExtractFromJSON(s, recipeData)
	}

	// Check if ingredients or instructions are still missing
	if recipeData["ingredients"] == "" {
		log.Printf("Warning: No ingredients found for %s, attempting recovery", s.url)
		// Last resort - check if the body contains the word "ingredient"
		bodyText := s.Body()
		if strings.Contains(bodyText, "ingredient") || strings.Contains(bodyText, "Ingredient") {
			log.Printf("Found ingredient text in body but couldn't extract structured data")
			recipeData["ingredients"] = "Ingredients mentioned in page"
		}
	}

	if recipeData["instructions"] == "" {
		log.Printf("Warning: No instructions found for %s, attempting recovery", s.url)
		// Last resort - check if the body contains instruction-related words
		bodyText := s.Body()
		instructionKeywords := []string{"instruction", "direction", "step", "method", "preparation"}
		for _, keyword := range instructionKeywords {
			if strings.Contains(bodyText, keyword) || strings.Contains(bodyText, strings.Title(keyword)) {
				log.Printf("Found %s text in body but couldn't extract structured data", keyword)
				recipeData["instructions"] = fmt.Sprintf("%s mentioned in page", strings.Title(keyword))
				break
			}
		}
	}

	// Log what we found for debugging
	log.Printf("Extracted recipe data from %s:", s.url)
	log.Printf("  Title: %s", recipeData["title"])
	log.Printf("  Name: %s", recipeData["name"])
	log.Printf("  Times: Prep=%s, Cook=%s, Total=%s",
		recipeData["prep_time"], recipeData["cook_time"], recipeData["total_time"])
	log.Printf("  Servings: %s", recipeData["servings"])

	// Check if ingredients were found
	if recipeData["ingredients"] != "" {
		ingredientCount := strings.Count(recipeData["ingredients"], ";") + 1
		log.Printf("  Ingredients: Found %d", ingredientCount)
	} else {
		log.Printf("  Ingredients: None found")
	}

	// Check if instructions were found
	if recipeData["instructions"] != "" {
		instructionCount := strings.Count(recipeData["instructions"], ";") + 1
		log.Printf("  Instructions: Found %d", instructionCount)
	} else {
		log.Printf("  Instructions: None found")
	}

	return recipeData
}

// GetRecipeTime extracts prep and total time, then calculates cook time
func (s *Scraper) GetRecipeTime() (string, string, string) {
	var prepTime, cookTime, totalTime string

	switch {
	case strings.Contains(s.site, "allrecipes.com"):
		// Try the Allrecipes specific format for times
		s.doc.Find(".recipe-meta-item").Each(func(i int, item *goquery.Selection) {
			headerText := item.Find(".recipe-meta-item-header").Text()
			valueText := item.Find(".recipe-meta-item-body").Text()

			headerText = strings.ToLower(strings.TrimSpace(headerText))
			valueText = strings.TrimSpace(valueText)

			if strings.Contains(headerText, "prep") || strings.Contains(headerText, "prep:") {
				prepTime = valueText
			} else if strings.Contains(headerText, "cook") || strings.Contains(headerText, "cook:") {
				cookTime = valueText
			} else if strings.Contains(headerText, "total") || strings.Contains(headerText, "total:") {
				totalTime = valueText
			}
		})

		// Try the newer structure (mm-recipes-details format)
		if prepTime == "" || cookTime == "" || totalTime == "" {
			s.doc.Find(".mm-recipes-details__item").Each(func(i int, item *goquery.Selection) {
				label := item.Find(".mm-recipes-details__label").Text()
				value := item.Find(".mm-recipes-details__value").Text()

				label = strings.ToLower(strings.TrimSpace(label))
				value = strings.TrimSpace(value)

				if strings.Contains(label, "prep time") {
					prepTime = value
				} else if strings.Contains(label, "cook time") {
					cookTime = value
				} else if strings.Contains(label, "total time") {
					totalTime = value
				}
			})
		}
	}

	// If any times are still missing, try schema.org approach
	if prepTime == "" {
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='prepTime'], [itemscope] [itemprop='prepTime']").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				prepTime = item.AttrOr("content", "")
				// Convert ISO duration to simple format if needed
				if strings.HasPrefix(prepTime, "PT") {
					prepTime = convertISODuration(prepTime)
				} else if prepTime == "" {
					prepTime = strings.TrimSpace(item.Text())
				}
			}
		})
	}

	if cookTime == "" {
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='cookTime'], [itemscope] [itemprop='cookTime']").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				cookTime = item.AttrOr("content", "")
				// Convert ISO duration to simple format if needed
				if strings.HasPrefix(cookTime, "PT") {
					cookTime = convertISODuration(cookTime)
				} else if cookTime == "" {
					cookTime = strings.TrimSpace(item.Text())
				}
			}
		})
	}

	if totalTime == "" {
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='totalTime'], [itemscope] [itemprop='totalTime']").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				totalTime = item.AttrOr("content", "")
				// Convert ISO duration to simple format if needed
				if strings.HasPrefix(totalTime, "PT") {
					totalTime = convertISODuration(totalTime)
				} else if totalTime == "" {
					totalTime = strings.TrimSpace(item.Text())
				}
			}
		})
	}

	// Try LD+JSON schema for times
	if prepTime == "" || cookTime == "" || totalTime == "" {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check if this is a Recipe
					if jsonType, ok := jsonMap["@type"].(string); ok && (jsonType == "Recipe" || strings.Contains(jsonType, "Recipe")) {
						if prepTime == "" {
							if pt, ok := jsonMap["prepTime"].(string); ok {
								prepTime = convertISODuration(pt)
							}
						}

						if cookTime == "" {
							if ct, ok := jsonMap["cookTime"].(string); ok {
								cookTime = convertISODuration(ct)
							}
						}

						if totalTime == "" {
							if tt, ok := jsonMap["totalTime"].(string); ok {
								totalTime = convertISODuration(tt)
							}
						}
					}
				} else {
					// Fallback to regex
					if prepTime == "" {
						re := regexp.MustCompile(`"prepTime"\s*:\s*"([^"]+)"`)
						matches := re.FindStringSubmatch(scriptContent)
						if len(matches) >= 2 {
							prepTime = convertISODuration(matches[1])
						}
					}

					if cookTime == "" {
						re := regexp.MustCompile(`"cookTime"\s*:\s*"([^"]+)"`)
						matches := re.FindStringSubmatch(scriptContent)
						if len(matches) >= 2 {
							cookTime = convertISODuration(matches[1])
						}
					}

					if totalTime == "" {
						re := regexp.MustCompile(`"totalTime"\s*:\s*"([^"]+)"`)
						matches := re.FindStringSubmatch(scriptContent)
						if len(matches) >= 2 {
							totalTime = convertISODuration(matches[1])
						}
					}
				}
			}
		})
	}

	// Site-specific time extraction for Delish.com
	if (prepTime == "" || totalTime == "") && strings.Contains(s.site, "delish.com") {
		s.doc.Find(".recipe-info-item, .prep-info").Each(func(i int, item *goquery.Selection) {
			text := strings.ToLower(item.Text())

			if strings.Contains(text, "prep") && prepTime == "" {
				prepTime = extractTime(text)
			} else if strings.Contains(text, "total") && totalTime == "" {
				totalTime = extractTime(text)
			}
		})
	}

	// Generic fallback for recipe time metadata
	if prepTime == "" || totalTime == "" || cookTime == "" {
		// Extract times based on common recipe site patterns
		s.doc.Find(".recipe-meta-item, .recipe-meta, .recipe-details, .recipe-info").Each(func(i int, item *goquery.Selection) {
			text := strings.ToLower(item.Text())
			if strings.Contains(text, "prep") && prepTime == "" {
				prepTime = extractTime(text)
			} else if strings.Contains(text, "cook") && cookTime == "" {
				cookTime = extractTime(text)
			} else if strings.Contains(text, "total") && totalTime == "" {
				totalTime = extractTime(text)
			}
		})
	}

	// Calculate missing times if possible
	if prepTime != "" && totalTime != "" && cookTime == "" {
		cookTime = calculateCookTime(prepTime, totalTime)
	} else if prepTime == "" && cookTime != "" && totalTime != "" {
		// Try to calculate prep time from cook and total time
		cookMinutes := extractMinutes(cookTime)
		totalMinutes := extractMinutes(totalTime)
		if cookMinutes > 0 && totalMinutes > 0 && totalMinutes >= cookMinutes {
			prepMinutes := totalMinutes - cookMinutes
			if prepMinutes >= 60 {
				hours := prepMinutes / 60
				minutes := prepMinutes % 60
				if minutes > 0 {
					prepTime = fmt.Sprintf("%d hr %d min", hours, minutes)
				} else {
					prepTime = fmt.Sprintf("%d hr", hours)
				}
			} else {
				prepTime = fmt.Sprintf("%d min", prepMinutes)
			}
		}
	} else if totalTime == "" && prepTime != "" && cookTime != "" {
		// Calculate total time from prep and cook time
		prepMinutes := extractMinutes(prepTime)
		cookMinutes := extractMinutes(cookTime)
		if prepMinutes > 0 && cookMinutes > 0 {
			totalMinutes := prepMinutes + cookMinutes
			if totalMinutes >= 60 {
				hours := totalMinutes / 60
				minutes := totalMinutes % 60
				if minutes > 0 {
					totalTime = fmt.Sprintf("%d hr %d min", hours, minutes)
				} else {
					totalTime = fmt.Sprintf("%d hr", hours)
				}
			} else {
				totalTime = fmt.Sprintf("%d min", totalMinutes)
			}
		}
	}

	return prepTime, cookTime, totalTime
}

// GetRecipeYield extracts servings/yield specifically for Allrecipes
func (s *Scraper) GetRecipeYield() string {
	var servings string

	switch {
	case strings.Contains(s.site, "allrecipes.com"):
		// Try recipe-meta-item format
		s.doc.Find(".recipe-meta-item").Each(func(i int, item *goquery.Selection) {
			headerText := item.Find(".recipe-meta-item-header").Text()
			valueText := item.Find(".recipe-meta-item-body").Text()

			headerText = strings.ToLower(strings.TrimSpace(headerText))
			valueText = strings.TrimSpace(valueText)

			if strings.Contains(headerText, "servings") ||
				strings.Contains(headerText, "yield") ||
				strings.Contains(headerText, "makes") {
				servings = valueText
			}
		})

		// Try the newer structure (mm-recipes-details format)
		if servings == "" {
			s.doc.Find(".mm-recipes-details__item").Each(func(i int, item *goquery.Selection) {
				label := item.Find(".mm-recipes-details__label").Text()
				value := item.Find(".mm-recipes-details__value").Text()

				label = strings.ToLower(strings.TrimSpace(label))
				value = strings.TrimSpace(value)

				if strings.Contains(label, "servings") {
					servings = value
				} else if strings.Contains(label, "yield") {
					servings = value
				}
			})
		}
	}

	// If still not found, try schema.org approach
	if servings == "" {
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeYield'], [itemscope] [itemprop='recipeYield']").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				servings = strings.TrimSpace(item.Text())
				// Look for content attribute if text is empty
				if servings == "" {
					servings = item.AttrOr("content", "")
				}
			}
		})
	}

	// Try LD+JSON schema for yield
	if servings == "" {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check if this is a Recipe
					if jsonType, ok := jsonMap["@type"].(string); ok && (jsonType == "Recipe" || strings.Contains(jsonType, "Recipe")) {
						// Try to get servings information
						if yield, ok := jsonMap["recipeYield"].(string); ok {
							servings = yield
						} else if yield, ok := jsonMap["recipeYield"].([]interface{}); ok && len(yield) > 0 {
							if yieldStr, ok := yield[0].(string); ok {
								servings = yieldStr
							}
						} else if yield, ok := jsonMap["yield"].(string); ok {
							servings = yield
						}
					}
				} else {
					// Fallback to regex
					re := regexp.MustCompile(`"recipeYield"\s*:\s*"([^"]+)"`)
					matches := re.FindStringSubmatch(scriptContent)
					if len(matches) >= 2 {
						servings = matches[1]
					}
				}
			}
		})
	}

	return servings
}

// GetRecipeImage extracts the recipe image URL
func (s *Scraper) GetRecipeImage() string {
	var image string

	// First try for og:image meta tag (most common)
	s.doc.Find("meta[property='og:image']").Each(func(i int, item *goquery.Selection) {
		if i == 0 { // Just get the first one
			image = item.AttrOr("content", "")
		}
	})

	// If not found, try for schema.org recipe image
	if image == "" {
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='image'], [itemscope] [itemprop='image']").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				image = item.AttrOr("src", "")
				if image == "" {
					image = item.AttrOr("content", "")
				}
			}
		})
	}

	// Check for JSON-LD schema with images
	if image == "" {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			if image != "" {
				return
			}

			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "image") {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check for image property
					if imgVal, ok := jsonMap["image"]; ok {
						// Image could be a string or an object
						switch img := imgVal.(type) {
						case string:
							image = img
						case map[string]interface{}:
							if url, ok := img["url"].(string); ok {
								image = url
							}
						case []interface{}:
							if len(img) > 0 {
								// Get the first image
								if imgStr, ok := img[0].(string); ok {
									image = imgStr
								} else if imgObj, ok := img[0].(map[string]interface{}); ok {
									if url, ok := imgObj["url"].(string); ok {
										image = url
									}
								}
							}
						}
					}
				} else {
					// Try regex as fallback
					re := regexp.MustCompile(`"image"\s*:\s*"([^"]+)"`)
					matches := re.FindStringSubmatch(scriptContent)
					if len(matches) >= 2 {
						image = matches[1]
					}
				}
			}
		})
	}

	// Site-specific image extraction as fallback
	if image == "" {
		switch {
		case strings.Contains(s.site, "delish.com"):
			s.doc.Find(".content-lede-image img, .recipe-image img, [data-journey-content] img").Each(func(i int, item *goquery.Selection) {
				if i == 0 {
					image = item.AttrOr("src", "")
					if image == "" {
						image = item.AttrOr("data-src", "")
					}
				}
			})
		case strings.Contains(s.site, "allrecipes.com"):
			s.doc.Find(".recipe-lead-media img, .lead-media img, .primary-media img, .primary-image__image").Each(func(i int, item *goquery.Selection) {
				if i == 0 {
					image = item.AttrOr("src", "")
					if image == "" {
						// Some newer Allrecipes images use data-src
						image = item.AttrOr("data-src", "")
					}
				}
			})
		case strings.Contains(s.site, "foodnetwork.com"):
			s.doc.Find(".m-MediaBlock__a-Image img, .recipe-lead-photo img").Each(func(i int, item *goquery.Selection) {
				if i == 0 {
					image = item.AttrOr("src", "")
				}
			})
		}
	}

	// Generic fallback: look for large images in the recipe
	if image == "" {
		s.doc.Find("img").Each(func(i int, item *goquery.Selection) {
			if image != "" {
				return
			}

			// Check for large images
			width, _ := item.Attr("width")
			height, _ := item.Attr("height")

			// Convert to int if possible
			widthInt := 0
			heightInt := 0
			fmt.Sscanf(width, "%d", &widthInt)
			fmt.Sscanf(height, "%d", &heightInt)

			if widthInt >= 400 && heightInt >= 300 {
				image = item.AttrOr("src", "")
			}
		})
	}

	return image
}

// GetRecipeName extracts the recipe name
func (s *Scraper) GetRecipeName() string {
	var name string

	switch {
	case strings.Contains(s.site, "allrecipes.com"):
		// Look for the recipe heading in both old and new Allrecipes formats
		s.doc.Find(".recipe-heading h1, .article-heading, .headline, .recipe-title").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				name = strings.TrimSpace(item.Text())
			}
		})
	}

	// If still no name, try schema.org recipe name
	if name == "" {
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='name'], [itemscope] [itemprop='name']").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				name = strings.TrimSpace(item.Text())
			}
		})
	}

	// Try LD+JSON schema
	if name == "" {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			if name != "" {
				return
			}

			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "name") {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check if this is a Recipe
					if jsonType, ok := jsonMap["@type"].(string); ok && (jsonType == "Recipe" || strings.Contains(jsonType, "Recipe")) {
						if nameVal, ok := jsonMap["name"].(string); ok {
							name = nameVal
						}
					}
				} else {
					// Fallback to regex
					re := regexp.MustCompile(`"name"\s*:\s*"([^"]+)"`)
					matches := re.FindStringSubmatch(scriptContent)
					if len(matches) >= 2 {
						name = matches[1]
					}
				}
			}
		})
	}

	// Site-specific extraction
	if name == "" {
		switch {
		case strings.Contains(s.site, "delish.com"):
			s.doc.Find("h1.content-header-title, h1.recipe-title, h1[data-journey-content]").Each(func(i int, item *goquery.Selection) {
				if i == 0 {
					name = strings.TrimSpace(item.Text())
				}
			})
		}
	}

	// Fallback to h1 if schema not found
	if name == "" {
		s.doc.Find("h1").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				name = strings.TrimSpace(item.Text())
			}
		})
	}

	// If still not found, try for common recipe title classes
	if name == "" {
		s.doc.Find(".recipe-title, .recipe-name, .recipeName, .headline").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				name = strings.TrimSpace(item.Text())
			}
		})
	}

	return name
}

// GetRecipeNutrition extracts calories and servings
func (s *Scraper) GetRecipeNutrition() (string, string) {
	var calories, servings string

	// Try schema.org nutrition information
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='nutrition'] [itemprop='calories'], [itemscope] [itemprop='nutrition'] [itemprop='calories']").Each(func(i int, item *goquery.Selection) {
		if i == 0 {
			calories = strings.TrimSpace(item.Text())
		}
	})

	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeYield'], [itemscope] [itemprop='recipeYield']").Each(func(i int, item *goquery.Selection) {
		if i == 0 {
			servings = strings.TrimSpace(item.Text())
			// Look for content attribute if text is empty
			if servings == "" {
				servings = item.AttrOr("content", "")
			}
		}
	})

	// Try LD+JSON schema for nutrition
	if calories == "" || servings == "" {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check if this is a Recipe
					if jsonType, ok := jsonMap["@type"].(string); ok && (jsonType == "Recipe" || strings.Contains(jsonType, "Recipe")) {
						// Try to get nutrition information
						if calories == "" {
							if nutrition, ok := jsonMap["nutrition"].(map[string]interface{}); ok {
								if cal, ok := nutrition["calories"].(string); ok {
									calories = cal
								}
							}
						}

						// Try to get servings information
						if servings == "" {
							if yield, ok := jsonMap["recipeYield"].(string); ok {
								servings = yield
							} else if yield, ok := jsonMap["recipeYield"].([]interface{}); ok && len(yield) > 0 {
								if yieldStr, ok := yield[0].(string); ok {
									servings = yieldStr
								}
							}
						}
					}
				} else {
					// Fallback to regex
					if calories == "" {
						re := regexp.MustCompile(`"calories"\s*:\s*"([^"]+)"`)
						matches := re.FindStringSubmatch(scriptContent)
						if len(matches) >= 2 {
							calories = matches[1]
						}
					}

					if servings == "" {
						re := regexp.MustCompile(`"recipeYield"\s*:\s*"([^"]+)"`)
						matches := re.FindStringSubmatch(scriptContent)
						if len(matches) >= 2 {
							servings = matches[1]
						}
					}
				}
			}
		})
	}

	// Site-specific extraction for Delish.com
	if (calories == "" || servings == "") && strings.Contains(s.site, "delish.com") {
		s.doc.Find(".recipe-info-item, .nutrition-info").Each(func(i int, item *goquery.Selection) {
			text := strings.ToLower(item.Text())

			if strings.Contains(text, "cal") && calories == "" {
				re := regexp.MustCompile(`(\d+)\s*cal`)
				matches := re.FindStringSubmatch(text)
				if len(matches) >= 2 {
					calories = matches[1]
				}
			} else if (strings.Contains(text, "serv") || strings.Contains(text, "yield") || strings.Contains(text, "makes")) && servings == "" {
				re := regexp.MustCompile(`(\d+[-\s]?\d*)`)
				matches := re.FindStringSubmatch(text)
				if len(matches) >= 2 {
					servings = matches[1]
				}
			}
		})
	}

	// Generic nutrition extraction as fallback
	if calories == "" {
		s.doc.Find(".nutrition-info, .recipe-nutrition, .nutrition-data").Each(func(i int, item *goquery.Selection) {
			text := strings.ToLower(item.Text())
			if strings.Contains(text, "calorie") || strings.Contains(text, "cal") {
				re := regexp.MustCompile(`(\d+)\s*cal`)
				matches := re.FindStringSubmatch(text)
				if len(matches) >= 2 {
					calories = matches[1]
				}
			}
		})
	}

	if servings == "" {
		s.doc.Find(".recipe-meta-item, .servings, .yield").Each(func(i int, item *goquery.Selection) {
			text := strings.ToLower(item.Text())
			if strings.Contains(text, "serv") || strings.Contains(text, "yield") || strings.Contains(text, "makes") {
				re := regexp.MustCompile(`(\d+[-\s]?\d*)`)
				matches := re.FindStringSubmatch(text)
				if len(matches) >= 2 {
					servings = matches[1]
				}
			}
		})
	}

	// Clean up calories to just have the number
	if calories != "" {
		reCalNum := regexp.MustCompile(`(\d+)`)
		calMatches := reCalNum.FindStringSubmatch(calories)
		if len(calMatches) >= 2 {
			calories = calMatches[1]
		}
	}

	return calories, servings
}

// GetRecipeCategories extracts recipe categories or tags
func (s *Scraper) GetRecipeCategories() []string {
	var categories []string

	// Try schema.org categories
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeCategory'], [itemscope] [itemprop='recipeCategory']").Each(func(i int, item *goquery.Selection) {
		category := strings.TrimSpace(item.Text())
		if category != "" {
			categories = append(categories, category)
		}
	})

	// Try LD+JSON schema for categories
	if len(categories) == 0 {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			if len(categories) > 0 {
				return
			}

			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "recipeCategory") {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check if this is a Recipe
					if jsonType, ok := jsonMap["@type"].(string); ok && (jsonType == "Recipe" || strings.Contains(jsonType, "Recipe")) {
						// Try to get category information
						if catList, ok := jsonMap["recipeCategory"].([]interface{}); ok {
							for _, cat := range catList {
								if catStr, ok := cat.(string); ok {
									categories = append(categories, catStr)
								}
							}
						} else if catStr, ok := jsonMap["recipeCategory"].(string); ok {
							categories = append(categories, catStr)
						}
					}
				} else {
					// Fallback to regex
					re := regexp.MustCompile(`"recipeCategory"\s*:\s*\[(.*?)\]`)
					matches := re.FindStringSubmatch(scriptContent)
					if len(matches) >= 2 {
						categoryListStr := matches[1]
						categoryRe := regexp.MustCompile(`"([^"]+)"`)
						categoryMatches := categoryRe.FindAllStringSubmatch(categoryListStr, -1)
						for _, match := range categoryMatches {
							if len(match) >= 2 {
								categories = append(categories, match[1])
							}
						}
					} else {
						re := regexp.MustCompile(`"recipeCategory"\s*:\s*"([^"]+)"`)
						matches := re.FindStringSubmatch(scriptContent)
						if len(matches) >= 2 {
							categories = append(categories, matches[1])
						}
					}
				}
			}
		})
	}

	// Site specific extraction for Delish.com
	if len(categories) == 0 && strings.Contains(s.site, "delish.com") {
		s.doc.Find(".tag-link, .recipe-tags a, .content-header-label").Each(func(i int, item *goquery.Selection) {
			category := strings.TrimSpace(item.Text())
			if category != "" {
				categories = append(categories, category)
			}
		})
	}

	// Fallback to common category tags
	if len(categories) == 0 {
		s.doc.Find(".recipe-categories a, .recipe-tags a, .category-tag").Each(func(i int, item *goquery.Selection) {
			category := strings.TrimSpace(item.Text())
			if category != "" {
				categories = append(categories, category)
			}
		})
	}

	return categories
}

// extractSiteSpecificData handles site-specific data extraction
func extractSiteSpecificData(s *Scraper, recipeData map[string]string) {
	// Different sites have different structures
	switch {
	case strings.Contains(s.site, "delish.com"):
		// Delish has unique selectors for certain elements
		s.doc.Find(".recipe-info-item").Each(func(i int, item *goquery.Selection) {
			text := item.Text()
			if strings.Contains(strings.ToLower(text), "yield") || strings.Contains(strings.ToLower(text), "serves") {
				servingReg := regexp.MustCompile(`\d+`)
				matches := servingReg.FindString(text)
				if matches != "" {
					recipeData["servings"] = matches
				}
			}
		})

	case strings.Contains(s.site, "allrecipes.com"):
		// AllRecipes structure
		if recipeData["servings"] == "" {
			s.doc.Find(".recipe-meta-item").Each(func(i int, item *goquery.Selection) {
				headerText := item.Find(".recipe-meta-item-header").Text()
				if strings.Contains(strings.ToLower(headerText), "servings") {
					servingText := item.Find(".recipe-meta-item-body").Text()
					recipeData["servings"] = strings.TrimSpace(servingText)
				}
			})
		}

	case strings.Contains(s.site, "foodnetwork.com"):
		// Food Network structure
		if recipeData["servings"] == "" {
			s.doc.Find(".o-RecipeInfo__m-Yield").Each(func(i int, item *goquery.Selection) {
				recipeData["servings"] = strings.TrimSpace(item.Text())
			})
		}
	}
}

// tryExtractFromJSON attempts to extract recipe data from JSON-LD scripts
func tryExtractFromJSON(s *Scraper, recipeData map[string]string) {
	s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
		scriptContent := item.Text()

		if !strings.Contains(scriptContent, "Recipe") {
			return
		}

		// Try to parse the JSON
		var jsonData interface{}
		if err := json.Unmarshal([]byte(scriptContent), &jsonData); err != nil {
			log.Printf("Error parsing JSON-LD: %v", err)
			return
		}

		// Helper function to recursively extract recipe data
		var extractRecipeData func(data interface{})
		extractRecipeData = func(data interface{}) {
			// Check if this is a map (JSON object)
			if dataMap, ok := data.(map[string]interface{}); ok {
				// Check if this is a Recipe
				if typeVal, ok := dataMap["@type"].(string); ok {
					if typeVal == "Recipe" || strings.Contains(typeVal, "Recipe") {
						// Extract recipe properties

						// Name
						if recipeData["name"] == "" {
							if name, ok := dataMap["name"].(string); ok && name != "" {
								recipeData["name"] = name
							}
						}

						// Image
						if recipeData["image"] == "" {
							if image, ok := dataMap["image"].(string); ok && image != "" {
								recipeData["image"] = image
							} else if imageObj, ok := dataMap["image"].(map[string]interface{}); ok {
								if url, ok := imageObj["url"].(string); ok && url != "" {
									recipeData["image"] = url
								}
							} else if imageArr, ok := dataMap["image"].([]interface{}); ok && len(imageArr) > 0 {
								if imgStr, ok := imageArr[0].(string); ok && imgStr != "" {
									recipeData["image"] = imgStr
								} else if imgObj, ok := imageArr[0].(map[string]interface{}); ok {
									if url, ok := imgObj["url"].(string); ok && url != "" {
										recipeData["image"] = url
									}
								}
							}
						}

						// Description
						if recipeData["description"] == "" {
							if desc, ok := dataMap["description"].(string); ok && desc != "" {
								recipeData["description"] = desc
							}
						}

						// Ingredients
						if recipeData["ingredients"] == "" {
							// Check various possible field names
							for _, field := range []string{"recipeIngredient", "ingredients"} {
								if ings, ok := dataMap[field].([]interface{}); ok && len(ings) > 0 {
									var ingredients []string
									for _, ing := range ings {
										if ingStr, ok := ing.(string); ok && ingStr != "" {
											ingredients = append(ingredients, ingStr)
										}
									}
									if len(ingredients) > 0 {
										recipeData["ingredients"] = strings.Join(ingredients, ";")
										break
									}
								}
							}
						}

						// Instructions
						if recipeData["instructions"] == "" {
							if insts, ok := dataMap["recipeInstructions"].([]interface{}); ok && len(insts) > 0 {
								var instructions []string
								for _, inst := range insts {
									if instStr, ok := inst.(string); ok && instStr != "" {
										instructions = append(instructions, instStr)
									} else if instObj, ok := inst.(map[string]interface{}); ok {
										if text, ok := instObj["text"].(string); ok && text != "" {
											instructions = append(instructions, text)
										}
									}
								}
								if len(instructions) > 0 {
									recipeData["instructions"] = strings.Join(instructions, ";")
								}
							} else if instStr, ok := dataMap["recipeInstructions"].(string); ok && instStr != "" {
								// Split into steps if it's a single string
								steps := strings.Split(instStr, ". ")
								recipeData["instructions"] = strings.Join(steps, ";")
							}
						}

						// Prep Time
						if recipeData["prep_time"] == "" {
							if pt, ok := dataMap["prepTime"].(string); ok && pt != "" {
								recipeData["prep_time"] = convertISODuration(pt)
							}
						}

						// Cook Time
						if recipeData["cook_time"] == "" {
							if ct, ok := dataMap["cookTime"].(string); ok && ct != "" {
								recipeData["cook_time"] = convertISODuration(ct)
							}
						}

						// Total Time
						if recipeData["total_time"] == "" {
							if tt, ok := dataMap["totalTime"].(string); ok && tt != "" {
								recipeData["total_time"] = convertISODuration(tt)
							}
						}

						// Servings
						if recipeData["servings"] == "" {
							if yield, ok := dataMap["recipeYield"].(string); ok && yield != "" {
								recipeData["servings"] = yield
							} else if yieldArr, ok := dataMap["recipeYield"].([]interface{}); ok && len(yieldArr) > 0 {
								if yStr, ok := yieldArr[0].(string); ok && yStr != "" {
									recipeData["servings"] = yStr
								}
							}
						}

						// Calories
						if recipeData["calories"] == "" {
							if nutrition, ok := dataMap["nutrition"].(map[string]interface{}); ok {
								if cal, ok := nutrition["calories"].(string); ok && cal != "" {
									recipeData["calories"] = cal
								}
							}
						}
					}
				}

				// Process nested arrays (like @graph)
				for key, value := range dataMap {
					if key == "@graph" {
						if graphArr, ok := value.([]interface{}); ok {
							for _, item := range graphArr {
								extractRecipeData(item)
							}
						}
					}
				}
			} else if dataArr, ok := data.([]interface{}); ok {
				// If it's an array, process each item
				for _, item := range dataArr {
					extractRecipeData(item)
				}
			}
		}

		// Start the extraction process
		extractRecipeData(jsonData)
	})
}

// Helper function to convert ISO duration to human-readable time
func convertISODuration(isoDuration string) string {
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(isoDuration)

	if len(matches) >= 4 {
		hours := matches[1]
		minutes := matches[2]
		seconds := matches[3]

		if hours != "" && minutes != "" {
			return hours + " hr " + minutes + " min"
		} else if hours != "" {
			return hours + " hr"
		} else if minutes != "" {
			return minutes + " min"
		} else if seconds != "" {
			return seconds + " sec"
		}
	}

	return isoDuration
}

// Helper function to extract time from text
func extractTime(text string) string {
	// First try to find pattern like "5 min" or "1 hour 30 min"
	re := regexp.MustCompile(`(\d+)\s*(min|hour|hr|minute|h|m)(?:\s+(\d+)\s*(min|minute|m))?`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 3 {
		if len(matches) >= 5 && matches[3] != "" {
			// Format with hours and minutes
			return fmt.Sprintf("%s %s %s %s", matches[1], matches[2], matches[3], matches[4])
		}
		// Simple format
		return matches[1] + " " + matches[2]
	}

	// Try simpler extraction - just grab first number followed by time unit
	re = regexp.MustCompile(`(\d+)\s*(minutes|minute|mins|min|hours|hour|hrs|hr)`)
	matches = re.FindStringSubmatch(text)
	if len(matches) >= 3 {
		return matches[1] + " " + matches[2]
	}

	return ""
}

// Helper function to calculate cook time from prep time and total time
func calculateCookTime(prepTime, totalTime string) string {
	// Extract minutes from prep time
	prepMinutes := extractMinutes(prepTime)

	// Extract minutes from total time
	totalMinutes := extractMinutes(totalTime)

	// If either conversion failed, return empty string
	if prepMinutes < 0 || totalMinutes < 0 {
		return ""
	}

	// Ensure total time is greater than or equal to prep time
	if totalMinutes < prepMinutes {
		return ""
	}

	// Calculate cook time in minutes
	cookMinutes := totalMinutes - prepMinutes

	// Convert back to string format
	if cookMinutes >= 60 {
		hours := cookMinutes / 60
		minutes := cookMinutes % 60
		if minutes > 0 {
			return fmt.Sprintf("%d hr %d min", hours, minutes)
		}
		return fmt.Sprintf("%d hr", hours)
	}

	return fmt.Sprintf("%d min", cookMinutes)
}

// Helper function to extract minutes from time string
func extractMinutes(timeStr string) int {
	// Extract hours
	reHour := regexp.MustCompile(`(\d+)\s*(hour|hr|h)`)
	hourMatches := reHour.FindStringSubmatch(timeStr)
	hours := 0
	if len(hourMatches) >= 3 {
		fmt.Sscanf(hourMatches[1], "%d", &hours)
	}

	// Extract minutes
	reMin := regexp.MustCompile(`(\d+)\s*(min|minute|m)`)
	minMatches := reMin.FindStringSubmatch(timeStr)
	minutes := 0
	if len(minMatches) >= 3 {
		fmt.Sscanf(minMatches[1], "%d", &minutes)
	}

	// If no matches found, return error code
	if len(hourMatches) < 3 && len(minMatches) < 3 {
		return -1
	}

	// Convert to total minutes
	return hours*60 + minutes
}
