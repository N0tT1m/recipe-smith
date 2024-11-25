package structs

// Page struct to store in database
type Page struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// SearchResult struct to handle search queries
type SearchResult struct {
	Pages []Page `json:"pages"`
	Input string `json:"input"`
}
