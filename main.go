package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const startYear = 2000
const endYear = 2021
const csvDir = "csv"
const countCsvFn = "count.csv"
const rankCsvFn = "rank.csv"

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
			station, _ = sjis2utf8(station)
		}

		if rank != "" {
			if year == endYear {
				stationList = append(stationList, station)
				stationMap[station] = make(map[int]record)
			}

			// 10,000 -> 10000
			count = strings.Replace(count, ",", "", -1)
			ci, _ := strconv.Atoi(count)
			ranki, _ := strconv.Atoi(rank)
			yearMapBySt := stationMap[station]

			// yearMapByStがnilの時は最新の年度TOP100に存在しない駅名のためSKIP
			if yearMapBySt != nil {
				yearMapBySt[year] = record{count: ci, rank: ranki}
				stationMap[station] = yearMapBySt
			}

		}

	})
}

func main() {

	var rankCsvPrint bool
	flag.BoolVar(&rankCsvPrint, "r", false, "rankCsvPrint")
	flag.Parse()

	stationMap = make(map[string]map[int]record)

	for i := endYear; i >= startYear; i-- {
		path := fmt.Sprintf("%s%d.html", "./htmls/", i)
		scrape(path, i)
	}

	genCSV(rankCsvPrint)
}

func genCSV(rank bool) {
	var filePath string
	records := [][]string{}

	checkSaveDir(csvDir)
	if rank {
		filePath = filepath.Join(csvDir, rankCsvFn)
	} else {
		filePath = filepath.Join(csvDir, countCsvFn)
	}

	// ヘッダ行処理
	header := []string{}
	header = append(header, "年度")
	header = append(header, stationList...)
	records = append(records, header)

	for i := startYear; i <= endYear; i++ {
		record := []string{}
		record = append(record, strconv.Itoa(i))
		for _, stationName := range stationList {
			cs := stationMap[stationName]
			var v int
			if rank {
				v = cs[i].rank
			} else {
				v = cs[i].count
			}
			record = append(record, strconv.Itoa(v))
		}
		records = append(records, record)
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(f)
	//https://pkg.go.dev/encoding/csv#Writer.Write
	//w.write + w.Flush するとWindowsDefenderが誤検知するのでw.WriteAllを使う
	err = w.WriteAll(records) // calls Flush internally
	if err != nil {
		log.Fatalln(err)
	}

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func sjis2utf8(str string) (string, error) {
	ret, err := io.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func checkSaveDir(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
