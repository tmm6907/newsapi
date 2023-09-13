package newsapi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// HTTP client for interfacing with NewsAPI.
type Client struct {
	apiKey string
	client *http.Client
}

// Provides controlled access to an HTTP response from request.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
	RequestURL *url.URL
	Message    string
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
