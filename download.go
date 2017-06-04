package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"bytes"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

var client *http.Client

func init() {
	client = &http.Client{
	// CheckRedirect: redirectPolicyFunc,
	}
}

func convertEncoding(encoding string, bodyReader io.Reader) []byte {
	reader, err := charset.NewReaderLabel(encoding, bodyReader)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return bodyBytes
}

func request(rw *RequestWrapper) *goquery.Document {
	rawurl := rw.request.URL.String()
	request := rw.request

	response, err := client.Do(request)
	if err != nil {
		log.Printf("http get failed!, %s\n", rawurl)
		return nil
	}

	// response, err := http.Get(rawurl)
	// if err != nil {
	// 	log.Printf("http get failed!, %s\n", rawurl)
	// 	return nil
	// }
	defer response.Body.Close()

	contentType := response.Header["Content-Type"]
	if len(contentType) == 0 {
		log.Printf("No Content-Type!, %s\n", rawurl)
		return nil
	}

	content, encoding := contentAndEncoding(contentType[0])
	if content == "" {
		log.Printf("Empty Content-Type!, %s\n", rawurl)
		return nil
	} else if content != "text/html" {
		log.Printf("Not html: %s, %s\n", content, rawurl)
		return nil
	}

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll(response.Body) failed!, %s\n", rawurl)
		return nil
	}

	if encoding == "utf-8" {
		doc := newDocFromByte(bodyBytes, rawurl)
		return doc
	} else if encoding == "" {
		doc := newDocFromByte(bodyBytes, rawurl)
		encoding = charsetInHTML(doc)
		if encoding == "" {
			log.Printf("No encoding!, %s\n", rawurl)
			return nil
		}
	}

	if encoding == "utf-8" {
		doc := newDocFromByte(bodyBytes, rawurl)
		return doc
	}
	bodyBytes = convertEncoding(encoding, bytes.NewReader(bodyBytes))
	if bodyBytes == nil {
		log.Printf("Convert encoding failed!, %s\n", rawurl)
		return nil
	}
	doc := newDocFromByte(bodyBytes, rawurl)

	return doc
}
