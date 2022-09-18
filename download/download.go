package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const startYear = 2000
const endYear = 2021
const baseUrl = "https://www.jreast.co.jp/passenger/"
const saveBaseDir = "htmls"

func main() {
	// 存在しなければカレントディレクトリに保存フォルダを作成する
	checkSaveDir(saveBaseDir)

	var url string
	for i := startYear; i <= endYear; i++ {
		if i == endYear {
			url = fmt.Sprintf("%s%s.html", baseUrl, "index")
		} else {
			url = fmt.Sprintf("%s%d.html", baseUrl, i)
		}
		fmt.Println(url)
		download(url, fmt.Sprintf("%d.html", i))
	}
}

func download(url, saveFilename string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath.Join(saveBaseDir, saveFilename))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func checkSaveDir(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				panic(err)
			}
		}
	}
}
