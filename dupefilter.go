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

func isDuplicate(conn ResourceConn, rawURL string) bool {
	host, urlfp := hostAndFingerPrint(rawURL)
	return redisSISMember(conn, host, urlfp)
}

func isDuplicateDebug(conn ResourceConn, rawURL string) bool {
	_, urlfp := hostAndFingerPrint(rawURL)
	return redisSISMember(conn, "www.163.com", urlfp)
}

func maskDupURLDebug(conn ResourceConn, rawURL string) bool {
	_, urlfp := hostAndFingerPrint(rawURL)
	return redisSADD(conn, "www.163.com", urlfp)
}

func maskDupURL(conn ResourceConn, rawURL string) bool {
	host, urlfp := hostAndFingerPrint(rawURL)
	return redisSADD(conn, host, urlfp)
}

func unmaskDupURL(conn ResourceConn, rawURL string) bool {
	host, urlfp := hostAndFingerPrint(rawURL)
	return redisSREM(conn, host, urlfp)
}
