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

// Scraper for each website
type Scraper struct {
	url  string
	doc  *goquery.Document
	site string
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
	}
}

// Body returns a string with the body of the page
func (s *Scraper) Body() string {
	body := s.doc.Find("body").Text()
	// Remove leading/ending white spaces
	body = strings.TrimSpace(body)

	return body
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
	var link string

	// Map of popular recipe site patterns
	recipePatterns := map[string][]string{
		"www.delish.com":        {"/cooking/recipe-ideas", "/everyday-cooking/quick-and-easy/", "/recipe/"},
		"www.allrecipes.com":    {"/recipe/", "/recipes/"},
		"www.foodnetwork.com":   {"/recipes/", "/recipe/"},
		"www.epicurious.com":    {"/recipes/", "/recipe/"},
		"www.simplyrecipes.com": {"/recipes/"},
		// Add more sites as needed
	}

	s.doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if !exists {
			return
		}

		// Skip anchors and javascript
		if strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript") {
			return
		}

		// Check if URL belongs to a known recipe site
		for domain, patterns := range recipePatterns {
			if strings.Contains(s.url, domain) {
				// Check against patterns for this domain
				for _, pattern := range patterns {
					if strings.Contains(href, pattern) {
						link = s.buildLinks(href)
						if link != "" {
							links = append(links, link)
						}
						break
					}
				}
				break
			}
		}
	})

	return removeDuplicateLinks(links)
}

// removeDuplicateLinks removes duplicate links from a slice
func removeDuplicateLinks(links []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, link := range links {
		if _, ok := seen[link]; !ok {
			seen[link] = true
			result = append(result, link)
		}
	}

	return result
}

// MetaDataInformation returns the title and description from the page
func (s *Scraper) MetaDataInformation() (string, string) {
	var t string
	var d string

	t = s.doc.Find("title").Contents().Text()

	s.doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("name", "") == "description" || item.AttrOr("property", "") == "og:description" {
			d = item.AttrOr("content", "")
		}
	})

	return t, d
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
			s.doc.Find(".recipe-lead-media img, .lead-media img, .primary-media img").Each(func(i int, item *goquery.Selection) {
				if i == 0 {
					image = item.AttrOr("src", "")
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

	// Try schema.org recipe name
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='name'], [itemscope] [itemprop='name']").Each(func(i int, item *goquery.Selection) {
		if i == 0 {
			name = strings.TrimSpace(item.Text())
		}
	})

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

// GetRecipeTime extracts prep, cook, and total time
func (s *Scraper) GetRecipeTime() (string, string, string) {
	var prepTime, cookTime, totalTime string

	// Try schema.org timing information
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
	if (prepTime == "" || cookTime == "" || totalTime == "") && strings.Contains(s.site, "delish.com") {
		s.doc.Find(".recipe-info-item, .prep-info").Each(func(i int, item *goquery.Selection) {
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

	// Generic fallback for recipe time metadata
	if prepTime == "" || cookTime == "" || totalTime == "" {
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

	return prepTime, cookTime, totalTime
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

// GetRecipeIngredients extracts recipe ingredients
func (s *Scraper) GetRecipeIngredients() string {
	var ingredients []string

	// Try schema.org ingredients
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeIngredient'], [itemtype='http://schema.org/Recipe'] [itemprop='ingredients'], [itemscope] [itemprop='recipeIngredient'], [itemscope] [itemprop='ingredients']").Each(func(i int, item *goquery.Selection) {
		ingredient := strings.TrimSpace(item.Text())
		if ingredient != "" {
			ingredients = append(ingredients, ingredient)
		}
	})

	// Try LD+JSON schema for ingredients
	if len(ingredients) == 0 {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			if len(ingredients) > 0 {
				return
			}

			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "recipeIngredient") {
				// Try to parse the JSON
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(scriptContent), &jsonMap); err == nil {
					// Check if this is a Recipe
					if jsonType, ok := jsonMap["@type"].(string); ok && (jsonType == "Recipe" || strings.Contains(jsonType, "Recipe")) {
						// Try to get ingredient information
						if ingList, ok := jsonMap["recipeIngredient"].([]interface{}); ok {
							for _, ing := range ingList {
								if ingStr, ok := ing.(string); ok {
									ingredients = append(ingredients, ingStr)
								}
							}
						}
					}
				} else {
					// Fallback to regex
					re := regexp.MustCompile(`"recipeIngredient"\s*:\s*\[(.*?)\]`)
					matches := re.FindStringSubmatch(scriptContent)
					if len(matches) >= 2 {
						ingredientListStr := matches[1]
						ingredientRe := regexp.MustCompile(`"([^"]+)"`)
						ingredientMatches := ingredientRe.FindAllStringSubmatch(ingredientListStr, -1)
						for _, match := range ingredientMatches {
							if len(match) >= 2 {
								ingredients = append(ingredients, match[1])
							}
						}
					}
				}
			}
		})
	}

	// Site-specific ingredient extraction for Delish.com
	if len(ingredients) == 0 && strings.Contains(s.site, "delish.com") {
		s.doc.Find(".ingredient-item, .ingredients-item, .ingredient-list li").Each(func(i int, item *goquery.Selection) {
			ingredient := strings.TrimSpace(item.Text())
			if ingredient != "" {
				ingredients = append(ingredients, ingredient)
			}
		})

		// Look for ingredients in the body text format
		if len(ingredients) == 0 {
			ingredientsSection := false
			s.doc.Find("p, li").Each(func(i int, item *goquery.Selection) {
				text := strings.TrimSpace(item.Text())

				// Check if this is the ingredients heading
				if strings.ToLower(text) == "ingredients" {
					ingredientsSection = true
					return
				}

				// If we've reached directions/instructions, stop collecting ingredients
				if strings.Contains(strings.ToLower(text), "direction") || strings.Contains(strings.ToLower(text), "instruction") {
					ingredientsSection = false
					return
				}

				// Collect ingredients if we're in the ingredients section
				if ingredientsSection && text != "" && !strings.Contains(strings.ToLower(text), "for serving") {
					ingredients = append(ingredients, text)
				}
			})
		}
	}

	// Fallback to common ingredient list patterns
	if len(ingredients) == 0 {
		s.doc.Find(".ingredients-item-name, .ingredient, .ingredient-list li").Each(func(i int, item *goquery.Selection) {
			ingredient := strings.TrimSpace(item.Text())
			if ingredient != "" {
				ingredients = append(ingredients, ingredient)
			}
		})
	}

	// Site-specific ingredient extraction as last resort
	if len(ingredients) == 0 {
		switch {
		case strings.Contains(s.site, "allrecipes.com"):
			s.doc.Find("[itemprop='ingredients'], .ingredients-item").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		case strings.Contains(s.site, "foodnetwork.com"):
			s.doc.Find(".o-Ingredients__a-Ingredient, .recipe-ingredients li").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		case strings.Contains(s.site, "epicurious.com"):
			s.doc.Find(".ingredient, .ingredients-list li").Each(func(i int, item *goquery.Selection) {
				ingredient := strings.TrimSpace(item.Text())
				if ingredient != "" {
					ingredients = append(ingredients, ingredient)
				}
			})
		}
	}

	// Clean up ingredients - remove excessive whitespace
	for i, ingredient := range ingredients {
		// Remove extra whitespace
		ingredient = regexp.MustCompile(`\s+`).ReplaceAllString(ingredient, " ")
		ingredients[i] = strings.TrimSpace(ingredient)
	}

	// Join with semicolons as per the DB structure
	return strings.Join(ingredients, ";")
}

// GetRecipeInstructions extracts recipe instructions
func (s *Scraper) GetRecipeInstructions() string {
	var instructions []string

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

	// Special handling for Delish.com
	if len(instructions) == 0 && strings.Contains(s.site, "delish.com") {
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

		// Try to find the directions section in paragraphs
		if len(instructions) == 0 {
			directionsSection := false
			s.doc.Find("p, h2, h3, h4").Each(func(i int, item *goquery.Selection) {
				text := strings.TrimSpace(item.Text())

				// Check for directions/instructions heading
				if strings.Contains(strings.ToLower(text), "direction") || strings.Contains(strings.ToLower(text), "instruction") {
					directionsSection = true
					return
				}

				// Collect instructions if we're in the directions section
				if directionsSection && text != "" && len(text) > 15 {
					// Look for step indicators
					if strings.HasPrefix(text, "Step") || regexp.MustCompile(`^\d+\)`).MatchString(text) {
						instructions = append(instructions, text)
					}
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

	// Site-specific instruction extraction as last resort
	if len(instructions) == 0 {
		switch {
		case strings.Contains(s.site, "allrecipes.com"):
			s.doc.Find(".step, .instructions-section-item").Each(func(i int, item *goquery.Selection) {
				instruction := strings.TrimSpace(item.Text())
				if instruction != "" {
					instructions = append(instructions, instruction)
				}
			})
		case strings.Contains(s.site, "foodnetwork.com"):
			s.doc.Find(".o-Method__m-Step, .recipe-directions-list li").Each(func(i int, item *goquery.Selection) {
				instruction := strings.TrimSpace(item.Text())
				if instruction != "" {
					instructions = append(instructions, instruction)
				}
			})
		case strings.Contains(s.site, "epicurious.com"):
			s.doc.Find(".preparation-step, .preparation-steps li").Each(func(i int, item *goquery.Selection) {
				instruction := strings.TrimSpace(item.Text())
				if instruction != "" {
					instructions = append(instructions, instruction)
				}
			})
		}
	}

	// Clean up instructions
	for i, instruction := range instructions {
		// Remove "Step X:" prefixes
		instruction = regexp.MustCompile(`^Step\s*\d+\s*:?\s*`).ReplaceAllString(instruction, "")

		// Remove excessive whitespace
		instruction = regexp.MustCompile(`\s+`).ReplaceAllString(instruction, " ")

		instructions[i] = strings.TrimSpace(instruction)
	}

	// Join with semicolons as per the DB structure
	return strings.Join(instructions, ";")
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

	// Log what we found for debugging
	log.Printf("Extracted recipe data from %s:", s.url)
	log.Printf("  Title: %s", recipeData["title"])
	log.Printf("  Name: %s", recipeData["name"])
	log.Printf("  Times: Prep=%s, Cook=%s, Total=%s",
		recipeData["prep_time"], recipeData["cook_time"], recipeData["total_time"])
	log.Printf("  Servings: %s", recipeData["servings"])
	log.Printf("  Ingredients: Found %d", strings.Count(recipeData["ingredients"], ";")+1)
	log.Printf("  Instructions: Found %d", strings.Count(recipeData["instructions"], ";")+1)

	return recipeData
}
