package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var regx = regexp.MustCompile(`(.+); ?charset=(.+)$`)

func contentAndEncoding(contentType []string, url string) (string, string) {
	content := "text/html"
	encoding := "utf-8"

	if len(contentType) == 1 {
		result := regx.FindAllStringSubmatch(contentType[0], 2)

		// []
		if len(result) == 0 {
			return content, encoding
		}
		//[[text/html; charset=utf-8 text/html utf-8]]
		if len(result[0]) == 2 {
			content = result[0][1]
		} else if len(result[0]) == 3 {
			content = result[0][1]
			encoding = result[0][2]
		}
		return content, encoding
	}
	log.Printf("len(contentType) != 1, %s\n", url)
	return content, encoding
}

func request(rawurl string) *goquery.Document {
	response, err := http.Get(rawurl)
	if err != nil {
		log.Printf("http get failed!, %s\n", rawurl)
	}
	defer response.Body.Close()

	contentType := response.Header["Content-Type"]
	content, encoding := contentAndEncoding(contentType, rawurl)
	// log.Printf("%s : %s\n", contentType[0], rawurl)
	// log.Printf("%s, %s : %s\n", content, encoding, rawurl)

	if content != "text/html" {
		log.Printf("%s %s\n", content, rawurl)
		return nil
	}
	doc := newDoc(response.Body, encoding, rawurl)
	return doc
}
