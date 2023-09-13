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

const PAGE_SIZE = 5

func main() {
	var articleResult newsapi.ArticleResponse
	// var sourceResult newsapi.SourceResponse
	apiKey := os.Getenv("NEWS_API_KEY")
	connection := newsapi.NewClient(apiKey)
	res, err := connection.Get(
		"everything",
		&newsapi.Config{
			Query:           "iphone",
			ExcludedDomains: []string{"techcrunch.com", "thenextweb.com"},
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	// if err := json.Unmarshal(body, &sourceResult); err != nil { // Parse []byte to go struct pointer
	// 	fmt.Println(err)
	// }
	if err := json.Unmarshal(body, &articleResult); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Failed to unmarshal", err)
	}
	fmt.Println(res.StatusCode, strings.ToUpper(articleResult.Status))
	fmt.Println(articleResult.Status)
	fmt.Println(articleResult.TotalResults, "results found")
	for i := 0; i < PAGE_SIZE && i < articleResult.TotalResults; i++ {
		fmt.Println("Article", i+1)
		fmt.Println("Title:", articleResult.Articles[i].Title)
		fmt.Println("Author:", articleResult.Articles[i].Author)
		fmt.Println("Description:", articleResult.Articles[i].Description)
		fmt.Println("PublishedAt:", articleResult.Articles[i].PublishedAt)
		fmt.Println()
	}
	// for i := 0; i < PAGE_SIZE && i < len(sourceResult.Sources); i++ {
	// 	fmt.Println("Source", i+1)
	// 	fmt.Println("ID:", sourceResult.Sources[i].ID)
	// 	fmt.Println("Name:", sourceResult.Sources[i].Name)
	// 	fmt.Println("Description:", sourceResult.Sources[i].Description)
	// 	fmt.Println("URL:", sourceResult.Sources[i].URL)
	// 	fmt.Println("Category:", sourceResult.Sources[i].Category)
	// 	fmt.Println()
	// }

}
