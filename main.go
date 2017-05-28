package main

import (
	"fmt"
	"log"
	"net/url"
)

func urlparse(rawurl string) {
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

}

func square(num int) int {
	return num * num
}

func mapper(f func(int) int, alist []int) []int {
	var a = make([]int, len(alist), len(alist))
	for index, val := range alist {

		a[index] = f(val)
	}
	return a
}

func callback() {
	alist := []int{4, 5, 6, 7}
	result := mapper(square, alist)
	fmt.Println(result)
}

func main() {
	rawurl := "http://www.163.com:80/index?page=1#active"
	// urlparse(rawurl)
	dupefilter()
	download()
	parse()
	storage()

}
