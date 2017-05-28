package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	fmt.Println("Hello World")

	rawurl := "http://www.163.com"

	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "google.com"
	q := u.Query()
	q.Set("q", "golang")
	u.RawQuery = q.Encode()
	fmt.Println(u)
}
