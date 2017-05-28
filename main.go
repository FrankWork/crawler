package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	fmt.Println("Hello World")

	rawurl := "http://www.163.com:80/index?page=1#active"
	// scheme://[userinfo@]host[:port]/path[?query][#fragment]
	// scheme:opaque[?query][#fragment]

	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.Scheme)
	fmt.Println(u.Host)
	fmt.Println(u.Path)
	fmt.Println(u.Fragment)
	fmt.Println(u.RawQuery)

	// u.Scheme = "https"
	// u.Host = "google.com"
	// q := u.Query()
	// q.Set("q", "golang")
	// u.RawQuery = q.Encode()
	// fmt.Println(u)
}
