package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const endYear = 2021
const startYear = 2020

type record struct {
	rank    string
	station string
	count   string
}

var stationList []string

var stationMap map[string]map[int]int

func scrape(path string, year int) {

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	stationMap = make(map[string]map[int]int)
	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {

		rank := s.Find("td:nth-child(1)").Text()
		station := s.Find("td.stationName").Text()
		count := s.Find("td:nth-child(5)").Text()

		if rank != "" {
			//	fmt.Printf("%s %s %s\n", rank, station, count)
			//r := record{rank, station, count}
			if year == endYear {
				stationList = append(stationList, station)
				stationMap[station] = map[int]int{}
			}

			ci, _ := strconv.Atoi(count)
			cs := stationMap[station]

			if cs == nil {
				cs = make(map[int]int)
			}

			cs[year] = ci
			stationMap[station] = cs
		}

	})
}

func main() {
	for i := endYear; i >= startYear; i-- {
		path := fmt.Sprintf("%s%d.html", "./htmls/", i)
		//fmt.Println(path)
		scrape(path, i)
	}
	genCSV()
}

func genCSV() {
	// 1行目を出力する
	for _, v := range stationList {
		fmt.Printf(",")
		fmt.Print(v)
	}
	fmt.Println("")

	for i := startYear; i <= endYear; i++ {
		fmt.Print(i)
		for _, v := range stationList {
			fmt.Printf(",")

			cs := stationMap[v]
			v := cs[i]

			fmt.Print(v)
		}
		fmt.Println("")
	}

}
