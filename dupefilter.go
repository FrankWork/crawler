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

func hostAndFingerPrint(rawURL string) (string, string) {
	urlStrct, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	urlfb := fingerPrint(rawURL)

	return urlStrct.Host, urlfb
}

func isDuplicate(rawURL string) bool {
	host, urlfp := hostAndFingerPrint(rawURL)
	return redisSISMember(host, urlfp)
}

func maskDupURL(rawURL string) bool {
	host, urlfp := hostAndFingerPrint(rawURL)
	return redisSADD(host, urlfp)
}

func unmaskDupURL(rawURL string) bool {
	host, urlfp := hostAndFingerPrint(rawURL)
	return redisSREM(host, urlfp)
}
