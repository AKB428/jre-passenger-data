package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get("https://www.jreast.co.jp/passenger/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {

		title := s.Find("td.stationName").Text()
		count := s.Find("td:nth-child(5)").Text()
		fmt.Printf("%d: %s %s\n", i, title, count)
	})
}

func main() {
	ExampleScrape()
}
