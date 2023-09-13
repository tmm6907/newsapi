package newsapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Constants containing field options
var (
	BASEURL           string = "https://newsapi.org/v2/"
	ENDPOINTS                = [...]string{"everything", "top-headlines", "top-headlines/sources", "sources"}
	SEARCH_IN_OPTIONS        = [...]string{"title", "description", "content"}
	CATEGORY_OPTIONS         = [...]string{
		"business",
		"entertainment",
		"general",
		"health",
		"science",
		"sports",
		"technology",
	}
	LANGUAGE_OPTIONS = [...]string{
		"ar",
		"de",
		"en",
		"es",
		"fr",
		"he",
		"it",
		"nl",
		"no",
		"pt",
		"ru",
		"sv",
		"ud",
		"zh",
	}
	COUNTRY_OPTIONS = [...]string{
		"ae",
		"ar",
		"at",
		"au",
		"be",
		"bg",
		"br",
		"ca",
		"ch",
		"cn",
		"co",
		"cu",
		"cz",
		"de",
		"eg",
		"fr",
		"gb",
		"gr",
		"hk",
		"hu",
		"id",
		"ie",
		"il",
		"in",
		"it",
		"jp",
		"kr",
		"lt",
		"lv",
		"ma",
		"mx",
		"my",
		"ng",
		"nl",
		"no",
		"nz",
		"ph",
		"pl",
		"pt",
		"ro",
		"rs",
		"ru",
		"sa",
		"se",
		"sg",
		"si",
		"sk",
		"th",
		"tr",
		"tw",
		"ua",
		"us",
		"ve",
		"za",
	}
	SORT_OPTIONS = [...]string{"relevancy", "popularity", "publishedAt"}
)

// Request config for setting query parameters.
type Config struct {
	Query           string
	SearchIn        []string
	Sources         []string
	Domains         []string
	ExcludedDomains []string
	From            string
	To              string
	Language        string
	Country         string
	SortBy          string
	PageSize        int
	Page            int
	Category        string
}

// Validate paramaters and return a URL-encoded string of the parameters.
func (c *Config) clean() (string, error) {
	params := url.Values{}

	if c.Query != "" {
		if len(c.Query) > 500 {
			return "", errors.New("error: query parameter 'q' exceeded max length of 500 characters")
		}
		params.Add("q", c.Query)
	}
	if len(c.SearchIn) > 0 {
		searchInString := ""
		for _, searchField := range c.SearchIn {
			for _, option := range SEARCH_IN_OPTIONS {
				if searchField == option {
					if len(searchInString) == 0 {
						searchInString += searchField
					} else {
						searchInString += fmt.Sprint(",", searchField)
					}
				}
			}
		}
		params.Add("searchIn", searchInString)
		if params.Get("searchIn") == "" {
			return "", errors.New(
				"error: invalid configuration, unrecognized value in query parameter: 'searchIn'",
			)
		}
	}
	if len(c.Sources) > 0 {
		optionString := ""
		for _, option := range c.Sources {
			if len(optionString) == 0 {
				optionString += option
			} else {
				optionString += fmt.Sprint(",", option)
			}
		}
		params.Add("sources", optionString)
	}
	if len(c.Domains) > 0 {
		optionString := ""
		for _, option := range c.Domains {
			if len(optionString) == 0 {
				optionString += option
			} else {
				optionString += fmt.Sprint(",", option)
			}
		}
		params.Add("domains", optionString)
	}
	if len(c.ExcludedDomains) > 0 {
		optionString := ""
		for _, option := range c.ExcludedDomains {
			if len(optionString) == 0 {
				optionString += option
			} else {
				optionString += fmt.Sprint(",", option)
			}
		}
		params.Add("excludeDomains", optionString)
	}
	if c.From != "" {
		params.Set("from", c.From)
	}
	if c.To != "" {
		params.Set("to", c.To)
	}
	if c.Language != "" {
		for _, language := range LANGUAGE_OPTIONS {
			if c.Language == language {
				params.Set("language", c.Language)
			}
		}
		if params.Get("language") == "" {
			return "", errors.New(
				"error: invalid configuration, unrecognized value in query parameter: 'language'",
			)
		}
	}
	if c.Country != "" {
		for _, country := range COUNTRY_OPTIONS {
			if c.Country == country {
				params.Set("country", c.Country)
			}
		}
		if params.Get("country") == "" {
			return "", errors.New(
				"error: invalid configuration, unrecognized value in query parameter: 'country'",
			)
		}
	}
	if c.SortBy != "" {
		for _, option := range SORT_OPTIONS {
			if c.SortBy == option {
				params.Set("sortBy", c.SortBy)
			}
		}
		if params.Get("sortBy") == "" {
			return "", errors.New(
				"error: invalid configuration, unrecognized value in query parameter: 'sortBy'",
			)
		}
	}
	if c.PageSize != 0 {
		if c.PageSize > 100 {
			return "", errors.New("error: query parameter 'pageSize' exceeded max size of 100")
		}
		params.Add("pageSize", fmt.Sprintf("%d", c.PageSize))
	}
	if c.Page != 0 {
		params.Add("page", fmt.Sprintf("%d", c.Page))
	}
	if c.Category != "" {
		for _, option := range CATEGORY_OPTIONS {
			if c.Category == option {
				params.Set("category", c.Category)
				break
			}
		}
		if params.Get("category") == "" {
			return "", errors.New(
				"error: invalid configuration, unrecognized value in query parameter: 'category'",
			)
		}
	}

	return params.Encode(), nil
}

// HTTP client for interfacing with NewsAPI.
type Client struct {
	apiKey string
	client *http.Client
}

// Intitializes a new HTTP client to interface with NewsAPI.
func NewsAPIClient(apiKey string) *Client {
	return &Client{
		client: &http.Client{},
		apiKey: apiKey,
	}
}

// Adds X-Api-Key and content-type headers to requests.
func (c *Client) prepareHeaders(r *http.Request) {
	r.Header.Set("X-API-Key", c.apiKey)
	r.Header.Set("Content-Type", "application/json")
}

// Handels GET requests to NewsAPI.
func (c *Client) Get(endpoint string, config *Config) (*Response, error) {
	formatedParams := ""
	if config != nil {
		paramString, err := config.clean()
		if err != nil {
			return nil, err
		}
		formatedParams = paramString
	} else {
		panic("too many configs")
	}

	for _, option := range ENDPOINTS {
		if endpoint == option {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", BASEURL+endpoint, formatedParams), nil)
			if err != nil {
				return nil, err
			}
			c.prepareHeaders(req)
			resp, err := c.client.Do(req)
			if err != nil {
				return nil, err
			}

			return &Response{
				StatusCode: resp.StatusCode,
				Header:     resp.Header,
				Body:       resp.Body,
				RequestURL: resp.Request.URL,
			}, nil
		}
	}
	return nil, fmt.Errorf("unrecognized endpoint: '%s', try again", endpoint)
}

// Provides controlled access to http.Response from request.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
	RequestURL *url.URL
	Message    string
}

// Represents Article response type.
type ArticleResponse struct {
	Status       string    `json:"status"`
	Code         string    `json:"code"`
	Message      string    `json:"message"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

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

// News article source.
type ArticleSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

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
