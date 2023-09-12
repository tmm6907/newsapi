package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/tmm6907/newsapi/newsapi"
)

const PAGE_SIZE = 50

func main() {
	var result newsapi.ClientResponse
	connection := newsapi.NewClient("81545856d8f0485fada5f42b69afa9d8")
	res, err := connection.Get(
		"everything",
		&newsapi.Config{
			Query:    "\"super conductors\"",
			Country:  "us",
			SortBy:   "relevancy",
			PageSize: PAGE_SIZE,
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
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
	}
	fmt.Println(res.StatusCode)
	fmt.Println(result.Status)
	fmt.Println(result.TotalResults, "results found")
	for i := 0; i < PAGE_SIZE && i < result.TotalResults; i++ {
		fmt.Println("Article", i+1)
		fmt.Println("Title: ", result.Articles[i].Title)
		fmt.Println("Author: ", result.Articles[i].Author)
		fmt.Println("Description: ", result.Articles[i].Description)
		fmt.Println("PublishedAt: ", result.Articles[i].PublishedAt)
		fmt.Println()
	}

}
