package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/url"
)

func fingerPrint(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)

	return fmt.Sprintf("%x", bytes)
}

func urlparse(rawurl string) {
	// scheme://[userinfo@]host[:port]/path[?query][#fragment]
	// scheme:opaque[?query][#fragment]

	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(u.Scheme)
	log.Println(u.Host)
	log.Println(u.Path)
	log.Println(u.Fragment)
	log.Println(u.RawQuery)
}

func isDuplicate(rawURL string) bool {
	urlStrct, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	urlfb := fingerPrint(rawURL)

	return redisSISMember(urlStrct.Host, urlfb)
}
