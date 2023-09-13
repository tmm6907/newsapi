package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tmm6907/newsapi/newsapi"
)

const PAGE_SIZE = 3

func main() {
	var everythingResult newsapi.ArticleResponse
	var headlinesResult newsapi.ArticleResponse
	var sourceResult newsapi.SourceResponse

	apiKey := os.Getenv("NEWS_API_KEY")

	// Initialize new client
	connection := newsapi.NewsAPIClient(apiKey)

	// Set query paramaters of GET request using newsapi.Config
	// GET /everything?q=iphone&excludedDomains=techcrunch.com,thenextweb.com
	everythingRes, err := connection.Get(
		"everything",
		&newsapi.Config{
			Sources:         []string{"bbc-news"},
			ExcludedDomains: []string{"techcrunch.com", "thenextweb.com"},
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer everythingRes.Body.Close()
	body, err := io.ReadAll(everythingRes.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Deserialize JSON everything response to a newsapi.ArticleResponse
	if err := json.Unmarshal(body, &everythingResult); err != nil {
		fmt.Println("Failed to unmarshal", err)
	}

	fmt.Println("Request URL:", everythingRes.RequestURL)
	fmt.Println("Status Code:", everythingRes.StatusCode, strings.ToUpper(everythingResult.Status))
	if everythingRes.StatusCode != 200 {
		fmt.Println("ErrorMessage:", everythingResult.Message)
		fmt.Println()
	}
	fmt.Println(everythingResult.TotalResults, "results found")
	for i := 0; i < PAGE_SIZE && i < everythingResult.TotalResults; i++ {
		fmt.Println("Article", i+1)
		fmt.Println("Title:", everythingResult.Articles[i].Title)
		fmt.Println("Author:", everythingResult.Articles[i].Author)
		fmt.Println("Description:", everythingResult.Articles[i].Description)
		fmt.Println("PublishedAt:", everythingResult.Articles[i].PublishedAt)
		fmt.Println()
	}

	// GET /top-headlines?searchIn=title&pageSize=PAGE_SIZE *errors*
	headlinesRes, err := connection.Get(
		"top-headlines",
		&newsapi.Config{
			SearchIn: []string{"title"},
			PageSize: PAGE_SIZE,
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer headlinesRes.Body.Close()
	headlinesBody, err := io.ReadAll(headlinesRes.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Deserialize JSON top-headlines response to a newsapi.ArticleResponse
	if err := json.Unmarshal(headlinesBody, &headlinesResult); err != nil {
		fmt.Println("Failed to unmarshal", err)
	}

	fmt.Println("Request URL:", headlinesRes.RequestURL)
	fmt.Println("Status:", headlinesRes.StatusCode, strings.ToUpper(headlinesResult.Status))
	if headlinesRes.StatusCode != 200 {
		fmt.Println("ErrorMessage:", headlinesResult.Message)
		fmt.Println()
	}
	fmt.Println(headlinesResult.TotalResults, "results found")
	for i := 0; i < PAGE_SIZE && i < headlinesResult.TotalResults; i++ {
		fmt.Println("Article", i+1)
		fmt.Println("Title:", headlinesResult.Articles[i].Title)
		fmt.Println("Author:", headlinesResult.Articles[i].Author)
		fmt.Println("Description:", headlinesResult.Articles[i].Description)
		fmt.Println("PublishedAt:", headlinesResult.Articles[i].PublishedAt)
		fmt.Println()
	}

	// GET /sources?
	sourcesRes, err := connection.Get(
		"sources",
		&newsapi.Config{},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer sourcesRes.Body.Close()
	sourcesBody, err := io.ReadAll(sourcesRes.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Deserialize JSON sources response to a newsapi.SourceResponse
	if err := json.Unmarshal(sourcesBody, &sourceResult); err != nil {
		fmt.Println("Failed to unmarshal", err)
	}

	fmt.Println("Request URL:", sourcesRes.RequestURL)
	fmt.Println("Status:", sourcesRes.StatusCode, strings.ToUpper(sourceResult.Status))
	if sourcesRes.StatusCode != 200 {
		fmt.Println("Error Message:", sourceResult.Message)
		fmt.Println()
	}
	fmt.Println(len(sourceResult.Sources), "results found")

	for i := 0; i < PAGE_SIZE && i < len(sourceResult.Sources); i++ {
		fmt.Println("Source", i+1)
		fmt.Println("ID:", sourceResult.Sources[i].ID)
		fmt.Println("Name:", sourceResult.Sources[i].Name)
		fmt.Println("Description:", sourceResult.Sources[i].Description)
		fmt.Println("URL:", sourceResult.Sources[i].URL)
		fmt.Println("Category:", sourceResult.Sources[i].Category)
		fmt.Println()
	}

}
