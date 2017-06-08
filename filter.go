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

type DuplicateURLFilter interface {
	isDuplicate(rawurl string) bool
	addURL(rawurl string)
	removeURL(rawurl string)
}

type DuplicateURLFilterLocal struct {
	urlSet map[string]int
}

func (d *DuplicateURLFilterLocal) isDuplicate(rawurl string) bool {
	_, urlfp := hostAndFingerPrint(rawurl)

	if d.urlSet[urlfp] == 0 {
		return false
	}
	return true
}

func (d *DuplicateURLFilterLocal) addURL(rawurl string) {
	_, urlfp := hostAndFingerPrint(rawurl)
	d.urlSet[urlfp]++
}

func (d *DuplicateURLFilterLocal) removeURL(rawurl string) {
	_, urlfp := hostAndFingerPrint(rawurl)
	d.urlSet[urlfp] = 0
}

type DuplicateURLFilterDistribute struct {
	defaultKey string
}

func (d *DuplicateURLFilterDistribute) isDuplicate(rawurl string) bool {
	conn, resource := redisPoolConnect()
	defer RedisResourcePool.Put(resource)

	host, urlfp := hostAndFingerPrint(rawurl)
	if d.defaultKey != "" {
		host = d.defaultKey
	}
	return redisSISMember(conn, host, urlfp)
}

func (d *DuplicateURLFilterDistribute) addURL(rawurl string) {
	conn, resource := redisPoolConnect()
	defer RedisResourcePool.Put(resource)

	host, urlfp := hostAndFingerPrint(rawurl)
	if d.defaultKey != "" {
		host = d.defaultKey
	}
	redisSADD(conn, host, urlfp)
}

func (d *DuplicateURLFilterDistribute) removeURL(rawurl string) {
	conn, resource := redisPoolConnect()
	defer RedisResourcePool.Put(resource)

	host, urlfp := hostAndFingerPrint(rawurl)
	if d.defaultKey != "" {
		host = d.defaultKey
	}
	redisSREM(conn, host, urlfp)
}
