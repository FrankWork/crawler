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

func NewDuplicateURLFilterLocal() *DuplicateURLFilterLocal {
	return &DuplicateURLFilterLocal{make(map[string]int)}
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
	rc         *RedisClient
}

func NewDuplicateURLFilterDistribute(defaultKey string, rc *RedisClient) *DuplicateURLFilterDistribute {
	return &DuplicateURLFilterDistribute{defaultKey, rc}
}

func (d *DuplicateURLFilterDistribute) isDuplicate(rawurl string) bool {
	host, urlfp := hostAndFingerPrint(rawurl)
	if d.defaultKey != "" {
		host = d.defaultKey
	}
	return d.rc.SIsMember(host, urlfp)
}

func (d *DuplicateURLFilterDistribute) addURL(rawurl string) {
	host, urlfp := hostAndFingerPrint(rawurl)
	if d.defaultKey != "" {
		host = d.defaultKey
	}
	d.rc.SAdd(host, urlfp)
}

func (d *DuplicateURLFilterDistribute) removeURL(rawurl string) {
	host, urlfp := hostAndFingerPrint(rawurl)
	if d.defaultKey != "" {
		host = d.defaultKey
	}
	d.rc.SRem(host, urlfp)
}
