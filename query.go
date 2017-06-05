package main

import (
	"bytes"
	"io"
	"log"
	"regexp"

	"strings"

	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var regx = regexp.MustCompile(`(.+?);? ?(charset=(.+))?$`)

func newDocFromByte(body []byte, url string) *goquery.Document {
	reader := bytes.NewReader(body)
	return newDocFromReader(reader, url)
}

func newDocFromReader(reader io.Reader, rawurl string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader(reader) failed! : %s\n", rawurl)
		return nil
	}
	urlPointer, err := url.Parse(rawurl)
	if err != nil {
		log.Printf("url.Parse failed in newDocFromReader! : %s\n", rawurl)
		return nil
	}
	doc.Url = urlPointer

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
	// fmt.Println(doc.Url.String())

	urlCount := make(map[string]int)
	// log.Println("=======================")
	// log.Println(doc.Text())
	doc.Find("a").Each(func(index int, sel *goquery.Selection) {
		rawurl, exists := sel.Attr("href")
		rawurl = strings.Trim(rawurl, " \t\n")
		// log.Println(rawurl)
		if !exists {
			return
		}
		if len(rawurl) == 0 || rawurl == "#" || strings.ContainsAny(rawurl, "{}") {
			return
		}
		if strings.Index(rawurl, "javascript") == 0 || strings.Index(rawurl, "mailto:") == 0 {
			return
		}
		if strings.Index(rawurl, "tel:") == 0 {
			return
		}

		// resolve relative url path
		urlPointer, err := url.Parse(rawurl)
		if err != nil {
			log.Printf("url.Parse failed: %s\n", rawurl)
			return
		}
		urlPointer = doc.Url.ResolveReference(urlPointer)
		rawurl = urlPointer.String()

		domainCount := 0
		for idx := range domains {
			if strings.Contains(rawurl, domains[idx]) {
				domainCount++
			}
		}
		if domainCount == 0 {
			// log.Printf("out of domain: %s\n", rawurl)
			return
		}

		urlCount[rawurl]++
	})

	return urlCount
}
