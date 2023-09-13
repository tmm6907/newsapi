package newsapi

// Represents Source response type.
type SourceResponse struct {
	Status  string   `json:"status"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Sources []Source `json:"sources"`
}

// Source obejct definition.
type Source struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	Country     string `json:"country"`
}
