package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	base := "https://www.jreast.co.jp/passenger/"

	var url string
	for i := 2000; i <= 2021; i++ {
		if i == 2021 {
			url = fmt.Sprintf("%s%s.html", base, "index")
		} else {
			url = fmt.Sprintf("%s%d.html", base, i)
		}
		fmt.Println(url)
		download(url, fmt.Sprintf("%d.html", i))
	}
}

func download(url, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("./htmls/" + filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
