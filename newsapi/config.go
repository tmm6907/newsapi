package newsapi

import (
	"errors"
	"fmt"
	"net/url"
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
