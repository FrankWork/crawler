package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

// Charset auto determine. Use golang.org/x/net/html/charset. Get page body and change it to utf-8
func changeCharsetEncodingAuto(contentTypeStr string, sor io.ReadCloser) []byte {
	var err error
	destReader, err := charset.NewReader(sor, contentTypeStr)

	if err != nil {
		log.Println(err.Error())
		destReader = sor
	}

	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
		log.Println(err.Error())
		// For gb2312, an error will be returned.
		// Error like: simplifiedchinese: invalid GBK encoding
		return nil
	}

	return sorbody
}

func request(rawurl string) *goquery.Document {
	response, err := http.Get(rawurl)
	if err != nil {
		log.Printf("http get failed!, %s\n", rawurl)
		return nil
	}
	defer response.Body.Close()

	// FIXME : index out of range

	content := response.Header["Content-Type"]
	if len(content) == 0 {
		log.Printf("No Content-Type!, %s\n", rawurl)
		return nil
	}
	bodyBytes := changeCharsetEncodingAuto(content[0], response.Body)
	if bodyBytes == nil {
		log.Printf("Convert encoding failed!, %s\n", rawurl)
		return nil
	}
	doc := newDoc(bodyBytes, rawurl)

	return doc
}
