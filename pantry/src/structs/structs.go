package structs

import "time"

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
}
