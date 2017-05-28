package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/iconv.v1"
)

func get(url string) *goquery.Document {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	contentType := response.Header["Content-Type"]
	encoding := charset(contentType)

	var doc *goquery.Document
	if encoding == "utf-8" {
		doc, err = goquery.NewDocumentFromReader(response.Body)
	} else {
		cd, err := iconv.Open("utf-8", encoding) // convert encoding to utf-8
		if err != nil {
			log.Fatal("iconv.Open failed!")
		}
		defer cd.Close()

		bufSize := 0 // default if zero
		reader := iconv.NewReader(cd, response.Body, bufSize)
		doc, err = goquery.NewDocumentFromReader(reader)
	}

	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func charset(contentType []string) string {
	if len(contentType) == 1 {
		prefix := "charset="
		index := strings.Index(contentType[0], prefix)
		if index == -1 {
			log.Fatal("no charset in contentType")
		}
		encoding := contentType[0][index+len(prefix):]
		return strings.ToLower(encoding)
	}
	log.Fatal(contentType)
	return "utf-8"
}

func getTitle(doc *goquery.Document) string {
	return doc.Find("title").Text()
}

func getLinks(doc *goquery.Document) map[string]int {
	// var indices []int
	// var urlArr []string
	urlCount := make(map[string]int)

	doc.Find("a").Each(func(index int, sel *goquery.Selection) {
		url, _ := sel.Attr("href")
		// urlArr = append(urlArr, url)
		urlCount[url]++
		// indices = append(indices, index)
	})
	// fmt.Println(len(urlArr))
	// fmt.Println(len(urlCount))
	// fmt.Println(urlArr[0])

	return urlCount
}

func goQuery() {
	url := "http://www.163.com"
	doc := get(url)

	title := getTitle(doc)
	fmt.Println(title)

	urlCount := getLinks(doc)
	fmt.Println("\"\":", urlCount[""])
	c := 0
	for key, value := range urlCount {
		if idx := strings.Index(key, "http"); idx != 0 {
			fmt.Println(key, ":", value)
		}
		if value > 1 {
			// fmt.Println(key)
			c++
		}
	}
	fmt.Println("c: ", c)
}
