package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var urlQueue = list.New()

type RequestWrapper struct {
	request *http.Request
	depth   int
}

func (r *RequestWrapper) print() {
	fmt.Println(r.request.URL.String(), r.depth)
}

type Queue struct {
	List *list.List
}

func (q *Queue) enqueue(r *RequestWrapper) {
	q.List.PushBack(r)
}
func (q *Queue) dequeue() *RequestWrapper {
	r := q.List.Front()
	q.List.Remove(r)
	return r.Value.(*RequestWrapper)
}
func (q *Queue) isEmpty() bool {
	return q.List.Len() == 0
}
func foo() {
	queue := Queue{list.New()}
	r1, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	queue.enqueue(&RequestWrapper{r1, 0})
	r2, _ := http.NewRequest("GET", "http://www.163.com", nil)
	queue.enqueue(&RequestWrapper{r2, 1})
	for !queue.isEmpty() {
		r := queue.dequeue()
		r.print()
	}
}

func client() {
	u, _ := url.Parse("http://www.163.com")

	client := &http.Client{
	// CheckRedirect: redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	req.Header.Add("depth", "1")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", resp.Header.Get("Content-Type"))
	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", result[:10])
}
