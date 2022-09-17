package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const endYear = 2021
const startYear = 2000

type record struct {
	rank  int
	count int
}

var stationList []string

var stationMap map[string]map[int]record

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

	/*
		if year <= 2019 {
			body := transform.NewReader(bufio.NewReader(f), japanese.ShiftJIS.NewDecoder())
			doc, err = goquery.NewDocumentFromReader(body)
			if err != nil {
				log.Fatal(err)
			}
		}*/
	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {

		rank := s.Find("td:nth-child(1)").Text()
		station := s.Find("td.stationName").Text()
		count := s.Find("td:nth-child(5)").Text()

		if year <= 2012 {
			station = s.Find("td:nth-child(2)").Text()
		}

		if year <= 2011 {
			count = s.Find("td:nth-child(3)").Text()
		}

		if year <= 2019 {
			station, _ = sjis_to_utf8(station)
		}

		if rank != "" {
			// fmt.Printf("%s %s %s\n", rank, station, count)
			//r := record{rank, station, count}
			if year == endYear {
				stationList = append(stationList, station)
				stationMap[station] = make(map[int]record)
			}

			count = strings.Replace(count, ",", "", -1)
			ci, _ := strconv.Atoi(count)
			ranki, _ := strconv.Atoi(rank)
			cs := stationMap[station]

			// csがnilの時は最新の年度TOP100に存在しない駅名のためSKIP
			if cs != nil {
				cs[year] = record{count: ci, rank: ranki}
				stationMap[station] = cs
			}

		}

	})
}

func main() {
	stationMap = make(map[string]map[int]record)

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
			v := cs[i].count

			fmt.Print(v)
		}
		fmt.Println("")
	}

	//fmt.Printf("%v\n", stationMap)

}

func sjis_to_utf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}
