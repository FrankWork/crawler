package main

import (
	"io"
	"log"

	"strings"

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
	urlCount := make(map[string]int)

	doc.Find("a").Each(func(index int, sel *goquery.Selection) {
		url, exists := sel.Attr("href")
		url = strings.Trim(url, " \t\n")
		if !exists {
			return
		}
		if len(url) == 0 || url == "#" || strings.ContainsAny(url, "{}") {
			return
		}
		if strings.Index(url, "javascript") == 0 || strings.Index(url, "mailto:") == 0 {
			return
		}
		if strings.Index(url, "http") != 0 {
			// TODO join, redirect
			// fmt.Printf("\tindex(http)!=0 %s\n", url)
			return
		}
		if strings.Index(url, "http") == 0 && !(strings.Contains(url, "163.com") || strings.Contains(url, "netease.com")) {
			// 126.com, youdao.com
			// fmt.Printf("\tout of domain %s\n", url)
			return
		}
		urlCount[url]++
	})

	return urlCount
}
