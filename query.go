package main

import (
	"bytes"
	"io"
	"log"
	"regexp"

	"strings"

	"github.com/PuerkitoBio/goquery"
)

var regx = regexp.MustCompile(`(.+?);? ?(charset=(.+))?$`)

func newDocFromByte(body []byte, url string) *goquery.Document {
	reader := bytes.NewReader(body)
	return newDocFromReader(reader, url)
}

func newDocFromReader(reader io.Reader, url string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader(reader) failed! : %s\n", url)
		return nil
	}
	return doc
}

func contentAndEncoding(contentType string) (string, string) {
	result := regx.FindAllStringSubmatch(contentType, 3)

	//[[text/html; charset=utf-8 text/html charset=utf-8 utf-8]]
	content := result[0][1]
	encoding := result[0][3]
	return content, encoding
}

func charsetInHTML(doc *goquery.Document) string {
	// <meta http-equiv="Content-Type" content="text/html; charset=gb2312" />
	// <meta charset="gb2312">
	charset := ""
	exists := false

	doc.Find("meta").EachWithBreak(func(index int, sel *goquery.Selection) bool {
		charset, exists = sel.Attr("charset")
		if exists {
			return false // break loop
		}

		charset, exists = sel.Attr("content")
		if exists {
			_, charset = contentAndEncoding(charset)
			if charset != "" {
				return false
			}
		}

		return true
	})

	return charset
}

func getTitle(doc *goquery.Document) string {
	return doc.Find("title").Text()
}

func getLinks(doc *goquery.Document) map[string]int {
	urlCount := make(map[string]int)
	// log.Println("=======================")
	// log.Println(doc.Text())
	doc.Find("a").Each(func(index int, sel *goquery.Selection) {
		url, exists := sel.Attr("href")
		url = strings.Trim(url, " \t\n")
		// log.Println(url)
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
