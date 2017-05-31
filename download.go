package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func request2(rawurl string) *goquery.Document {

	response, err := http.Get(rawurl)
	if err != nil {
		log.Printf("http get failed!, %s\n", rawurl)
		return nil
	}
	defer response.Body.Close()

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

	if encoding != "" {
		return newDoc(response.Body, encoding, rawurl)
	}

	// encoding == ""
	// Content-Type and charset from html
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll failed!, %s\n", rawurl)
		return nil
	}

	doc := newDoc(bytes.NewReader(bodyBytes), "utf-8", rawurl)
	encoding = charsetInHTML(doc)
	if encoding == "" {
		return doc
	}
	doc = newDoc(bytes.NewReader(bodyBytes), encoding, rawurl)
	// fmt.Println(getTitle(doc))

	return doc
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
		doc := newDoc(bytes.NewReader(bodyBytes), "utf-8", rawurl)
		encoding = charsetInHTML(doc)
		if encoding == "" {
			return doc
		}
	}

	doc := newDoc(bytes.NewReader(bodyBytes), encoding, rawurl)
	// fmt.Println(getTitle(doc))

	return doc
}
