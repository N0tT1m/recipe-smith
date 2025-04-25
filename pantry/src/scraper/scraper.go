package scraper

import (
	"fmt"
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

	response, err := client.Get(u)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("Failed to fetch %s, status code: %d\n", u, response.StatusCode)
		return nil
	}

	d, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Parse the URL to get the site name
	parsedURL, err := url.Parse(u)
	if err != nil {
		fmt.Println(err)
		return nil
	}

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
		"www.delish.com":        {"/cooking/recipe-ideas", "/everyday-cooking/quick-and-easy/"},
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

	return links
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
		s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='image']").Each(func(i int, item *goquery.Selection) {
			if i == 0 {
				image = item.AttrOr("src", "")
				if image == "" {
					image = item.AttrOr("content", "")
				}
			}
		})
	}

	// Site-specific image extraction as fallback
	if image == "" {
		switch {
		case strings.Contains(s.site, "delish.com"):
			s.doc.Find(".content-lede-image img, .recipe-image img").Each(func(i int, item *goquery.Selection) {
				if i == 0 {
					image = item.AttrOr("src", "")
				}
			})
		case strings.Contains(s.site, "allrecipes.com"):
			s.doc.Find(".recipe-lead-media img, .lead-media img").Each(func(i int, item *goquery.Selection) {
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
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='name']").Each(func(i int, item *goquery.Selection) {
		if i == 0 {
			name = strings.TrimSpace(item.Text())
		}
	})

	// Try LD+JSON schema
	if name == "" {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "name") {
				re := regexp.MustCompile(`"name"\s*:\s*"([^"]+)"`)
				matches := re.FindStringSubmatch(scriptContent)
				if len(matches) >= 2 {
					name = matches[1]
				}
			}
		})
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
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='prepTime']").Each(func(i int, item *goquery.Selection) {
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

	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='cookTime']").Each(func(i int, item *goquery.Selection) {
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

	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='totalTime']").Each(func(i int, item *goquery.Selection) {
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
		})
	}

	// Site-specific time extraction as fallback
	if prepTime == "" || cookTime == "" || totalTime == "" {
		// Extract times based on common recipe site patterns
		s.doc.Find(".recipe-meta-item, .recipe-meta, .recipe-details").Each(func(i int, item *goquery.Selection) {
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
	re := regexp.MustCompile(`(\d+)\s*(min|hour|hr|minute|h|m)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 3 {
		return matches[1] + " " + matches[2]
	}
	return ""
}

// Helper function to convert ISO duration to human-readable time
func convertISODuration(isoDuration string) string {
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?`)
	matches := re.FindStringSubmatch(isoDuration)

	if len(matches) >= 3 {
		hours := matches[1]
		minutes := matches[2]

		if hours != "" && minutes != "" {
			return hours + " hr " + minutes + " min"
		} else if hours != "" {
			return hours + " hr"
		} else if minutes != "" {
			return minutes + " min"
		}
	}

	return isoDuration
}

// GetRecipeNutrition extracts calories and servings
func (s *Scraper) GetRecipeNutrition() (string, string) {
	var calories, servings string

	// Try schema.org nutrition information
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='nutrition'] [itemprop='calories']").Each(func(i int, item *goquery.Selection) {
		if i == 0 {
			calories = strings.TrimSpace(item.Text())
		}
	})

	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeYield']").Each(func(i int, item *goquery.Selection) {
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
		})
	}

	// Site-specific nutrition extraction as fallback
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
				re := regexp.MustCompile(`(\d+)[-\s]?(\d+)?`)
				matches := re.FindStringSubmatch(text)
				if len(matches) >= 2 {
					servings = matches[0]
				}
			}
		})
	}

	return calories, servings
}

// GetRecipeIngredients extracts recipe ingredients
func (s *Scraper) GetRecipeIngredients() string {
	var ingredients []string

	// Try schema.org ingredients
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeIngredient'], [itemtype='http://schema.org/Recipe'] [itemprop='ingredients']").Each(func(i int, item *goquery.Selection) {
		ingredient := strings.TrimSpace(item.Text())
		if ingredient != "" {
			ingredients = append(ingredients, ingredient)
		}
	})

	// Try LD+JSON schema for ingredients
	if len(ingredients) == 0 {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "recipeIngredient") {
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
		})
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
		case strings.Contains(s.site, "delish.com"):
			s.doc.Find(".ingredient-item, .ingredients-item").Each(func(i int, item *goquery.Selection) {
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

	// Join with semicolons as per the DB structure
	return strings.Join(ingredients, ";")
}

// GetRecipeInstructions extracts recipe instructions
func (s *Scraper) GetRecipeInstructions() string {
	var instructions []string

	// Try schema.org instructions
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeInstructions']").Each(func(i int, item *goquery.Selection) {
		instruction := strings.TrimSpace(item.Text())
		if instruction != "" {
			instructions = append(instructions, instruction)
		}
	})

	// Try LD+JSON schema for instructions
	if len(instructions) == 0 {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "recipeInstructions") {
				// Check if instructions are an array
				re := regexp.MustCompile(`"recipeInstructions"\s*:\s*\[(.*?)\]`)
				matches := re.FindStringSubmatch(scriptContent)
				if len(matches) >= 2 {
					instructionListStr := matches[1]
					instructionRe := regexp.MustCompile(`"([^"]+)"`)
					instructionMatches := instructionRe.FindAllStringSubmatch(instructionListStr, -1)
					for _, match := range instructionMatches {
						if len(match) >= 2 {
							instructions = append(instructions, match[1])
						}
					}
				} else {
					// Check if it's a string
					re := regexp.MustCompile(`"recipeInstructions"\s*:\s*"([^"]+)"`)
					matches := re.FindStringSubmatch(scriptContent)
					if len(matches) >= 2 {
						instructions = append(instructions, matches[1])
					}
				}
			}
		})
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
		case strings.Contains(s.site, "delish.com"):
			s.doc.Find(".direction-item, .preparation-steps li").Each(func(i int, item *goquery.Selection) {
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

	// Join with semicolons as per the DB structure
	return strings.Join(instructions, ";")
}

// GetRecipeCategories extracts recipe categories or tags
func (s *Scraper) GetRecipeCategories() []string {
	var categories []string

	// Try schema.org categories
	s.doc.Find("[itemtype='http://schema.org/Recipe'] [itemprop='recipeCategory']").Each(func(i int, item *goquery.Selection) {
		category := strings.TrimSpace(item.Text())
		if category != "" {
			categories = append(categories, category)
		}
	})

	// Try LD+JSON schema for categories
	if len(categories) == 0 {
		s.doc.Find("script[type='application/ld+json']").Each(func(i int, item *goquery.Selection) {
			scriptContent := item.Text()
			if strings.Contains(scriptContent, "Recipe") && strings.Contains(scriptContent, "recipeCategory") {
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

	return recipeData
}
