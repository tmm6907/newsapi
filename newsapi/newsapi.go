package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	HOSTNAME          string = "https://newsapi.org/v2/"
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
	LANGUAGES = [...]string{"ar",
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
	COUNTRIES = [...]string{
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

func (c *Config) clean() (string, error) {
	// Validate all fields with limited options and URLencode all fields
	//fmt.Errorf("enter valid for %s field", err)
	params := url.Values{}
	switch {
	case c.Query != "":
		if len(c.Query) > 500 {
			return "", errors.New("error: query parameter 'q' exceeded max length of 500 characters")
		}
		params.Add("q", c.Query)
	case len(c.SearchIn) > 0:
		for _, searchField := range c.SearchIn {
			for _, option := range SEARCH_IN_OPTIONS {
				if searchField == option {
					params.Add("searchIn", searchField)
				}
			}
		}
		if params.Get("searchIn") == "" {
			return "", errors.New(
				"error: invalid configuration, unrecognized value in query parameter: 'searchIn'",
			)
		}
	case len(c.Sources) > 0:
		for _, option := range c.Sources {
			params.Add("sources", option)
		}
	case len(c.Domains) > 0:
		for _, option := range c.Domains {
			params.Add("domains", option)
		}
	case len(c.ExcludedDomains) > 0:
		for _, option := range c.ExcludedDomains {
			params.Add("excludedDomains", option)
		}
	case c.From != "":
		params.Add("from", c.From)
	case c.To != "":
		params.Add("to", c.To)
	case c.Language != "":
		for _, language := range LANGUAGES {
			if c.Language == language {
				params.Set("language", c.Language)
			}
		}
		if params.Get("language") == "" {
			return "", errors.New(
				"error: invalid configuration, unrecognized value in query parameter: 'language'",
			)
		}
	case c.SortBy != "":
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
	case c.PageSize != 0:
		if c.PageSize > 100 {
			return "", errors.New("error: query parameter 'pageSize' exceeded max size of 100")
		}
		params.Add("pageSize", fmt.Sprintf("%d", c.PageSize))
	case c.Page != 0:
		params.Add("page", fmt.Sprintf("%d", c.Page))
	case c.Category != "":
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
	default:
		return "", errors.New("error: no configuration options set, query parameters not satisfied")
	}

	return params.Encode(), nil
}

type NewsAPIClient struct {
	apiKey string
	config Config
	client *http.Client
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
}

func NewClient(apiKey string) *NewsAPIClient {
	return &NewsAPIClient{
		client: &http.Client{},
		apiKey: apiKey,
		config: Config{},
	}
}

func (c *NewsAPIClient) prepareHeaders(r *http.Request) {
	r.Header.Set("X-API-Key", c.apiKey)
	r.Header.Set("Content-Type", "application/json")
}
func (c *NewsAPIClient) Get(endpoint string, config ...Config) (*Response, error) {
	formatedParams := ""
	if len(config) == 1 {
		paramString, err := config[0].clean()
		if err != nil {
			return nil, err
		}
		formatedParams = paramString
		fmt.Println("Formated: ", formatedParams)
	}

	fmt.Println("Count: ", len(config))

	for _, option := range ENDPOINTS {
		if endpoint == option {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", HOSTNAME+endpoint, formatedParams), nil)
			if err != nil {
				return nil, err
			}
			c.prepareHeaders(req)
			resp, err := c.client.Do(req)
			fmt.Println(resp, err)
			if err != nil {
				return &Response{
					StatusCode: resp.StatusCode,
					Header:     resp.Header,
					Body:       resp.Body,
				}, err
			}
			defer resp.Body.Close()
			return &Response{
				StatusCode: resp.StatusCode,
				Header:     resp.Header,
				Body:       resp.Body,
			}, nil
		}
	}
	return nil, fmt.Errorf("unrecognized endpoint: '%s', try again", endpoint)
}

func main() {
	connection := NewClient("81545856d8f0485fada5f42b69afa9d8")
	res, err := connection.Get(
		"everything",
		Config{
			Language: "sv",
			Country:  "se",
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	if res != nil {
		fmt.Println(res.StatusCode)
		fmt.Println(res.Body)

	}
	fmt.Println("none found")

}