package main

import (
	"io"
	"log"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/iconv.v1"
)

func newDoc(body io.ReadCloser, encoding string) *goquery.Document {
	if encoding != "utf-8" {
		cd, err := iconv.Open("utf-8", encoding) // convert encoding to utf-8
		if err != nil {
			log.Fatal("iconv.Open failed!")
		}
		defer cd.Close()

		bufSize := 0 // default if zero
		reader := iconv.NewReader(cd, body, bufSize)
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Fatal(err)
		}
		return doc
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
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
