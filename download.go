package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TODO merge 2 function as 1
func parseCharset(contentType []string, url string) string {
	defaultEncoding := "utf-8"

	if len(contentType) == 1 {
		prefix := "charset="
		index := strings.Index(contentType[0], prefix)
		if index == -1 {
			log.Printf("no charset in contentType, %s\n", url)
			return defaultEncoding
		}
		encoding := contentType[0][index+len(prefix):]
		return strings.ToLower(encoding)
	}
	log.Printf("len(contentType) != 1, %s\n", url)
	return defaultEncoding
}

func parseContentType(contentType []string, url string) string {
	defaultType := "text/html"
	if len(contentType) == 1 {
		prefix := "Content-Type: "
		index := strings.Index(contentType[0], prefix)
		if index == -1 {
			log.Printf("no `Content-Type` in `contentType`, %s\n", url)
			return defaultType
		}
		endIndex := strings.Index(contentType[0], "; charset")
		if endIndex == -1 {
			endIndex = len(contentType[0])
		}
		defaultType = contentType[0][index+len(prefix) : endIndex]
		return strings.ToLower(defaultType)
	}
	log.Printf("len(contentType) != 1, %s\n", url)
	return defaultType
}

func request(rawurl string) *goquery.Document {
	response, err := http.Get(rawurl)
	if err != nil {
		log.Printf("http get failed!, %s\n", rawurl)
	}
	defer response.Body.Close()

	contentType := response.Header["Content-Type"]
	encoding := parseCharset(contentType, rawurl)
	content := parseContentType(contentType, rawurl)
	if content != "text/html" {
		log.Printf("%s %s\n", content, rawurl)
		return nil
	}
	doc := newDoc(response.Body, encoding)
	return doc
}
