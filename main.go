package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func scrape(path string) {

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
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
	for i := 2021; i >= 2000; i-- {
		path := fmt.Sprintf("%s%d.html", "./htmls/", i)
		fmt.Println(path)
		scrape(path)
	}
}
