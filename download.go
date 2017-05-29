package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func charset(contentType []string) string {
	defaultEncoding := "utf-8"

	if len(contentType) == 1 {
		prefix := "charset="
		index := strings.Index(contentType[0], prefix)
		if index == -1 {
			log.Println("no charset in contentType")
			return defaultEncoding
		}
		encoding := contentType[0][index+len(prefix):]
		return strings.ToLower(encoding)
	}
	log.Fatal("len(contentType) != 1")
	return defaultEncoding
}

func request(rawurl string) *goquery.Document {
	response, err := http.Get(rawurl)
	if err != nil {
		log.Fatal("http get failed!")
	}
	defer response.Body.Close()

	contentType := response.Header["Content-Type"]
	encoding := charset(contentType)

	doc := newDoc(response.Body, encoding)
	maskDupURL(rawurl)
	return doc
}
