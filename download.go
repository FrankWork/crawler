package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	qconv "gopkg.in/iconv.v1"
)

func convertStr(str string, fromEncode string, toEncode string) string {
	cd, err := qconv.Open(toEncode, fromEncode)
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return ""
	}
	defer cd.Close()

	return cd.ConvString(str)
}

func request(rawurl string) *goquery.Document {
	response, err := http.Get(rawurl)
	if err != nil {
		log.Printf("http get failed!, %s\n", rawurl)
		return nil
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll failed!, %s\n", rawurl)
		return nil
	}

	// Content-Type and charset from Header
	contentType := response.Header["Content-Type"]
	content, encoding := contentAndEncoding(contentType[0])

	if content == "" {
		log.Printf("content == `` %s\n", rawurl)
		return nil
	} else if content != "text/html" {
		log.Printf("%s %s\n", content, rawurl)
		return nil
	}

	if encoding == "" {
		doc := newDoc(bodyBytes, rawurl)
		encoding = charsetInHTML(doc)
		if encoding == "" {
			return doc
		}
	}
	bodyStr := convertStr(string(bodyBytes), encoding, "utf-8")
	log.Println(bodyStr)
	doc := newDoc([]byte(bodyStr), rawurl)

	return doc
}
