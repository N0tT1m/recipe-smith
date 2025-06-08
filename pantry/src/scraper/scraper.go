package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

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
}

// NewScraper builds a new scraper for the website
func NewScraper(u string) *Scraper {
	if !strings.HasPrefix(u, "http") {
		return nil
	}

	response, err := http.Get(u)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer response.Body.Close()

	d, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	site := getSiteConfig(u)

	return &Scraper{
		url:  u,
		doc:  d,
		site: site,
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
