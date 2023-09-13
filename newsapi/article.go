package newsapi

// News article.
type Article struct {
	Source      ArticleSource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	ImageURL    string        `json:"urlToImage"`
	PublishedAt string        `json:"publishedAt"`
	Content     string        `json:"content"`
}

// Represents Article response type.
type ArticleResponse struct {
	Status       string    `json:"status"`
	Code         string    `json:"code"`
	Message      string    `json:"message"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

// News article source.
type ArticleSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
