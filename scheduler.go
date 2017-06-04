package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
)

type RequestWrapper struct {
	request *http.Request
	depth   int
}

func NewRequestWrapper(rawurl string, depth int) *RequestWrapper {
	request, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		log.Printf("new request failed!, %s\n", rawurl)
		return nil
	}
	// request.Header.Add("If-None-Match", `W/"wyzzy"`)
	// request.Header.Add("depth", "1")
	return &RequestWrapper{request, depth}
}

func (r *RequestWrapper) print() {
	fmt.Println(r.request.URL.String(), r.depth)
}

type RequestQueue struct {
	List *list.List
}

func (q *RequestQueue) enqueue(r *RequestWrapper) {
	q.List.PushBack(r)
}
func (q *RequestQueue) dequeue() *RequestWrapper {
	r := q.List.Front()
	q.List.Remove(r)
	return r.Value.(*RequestWrapper)
}
func (q *RequestQueue) isEmpty() bool {
	return q.List.Len() == 0
}

var requestQueue *RequestQueue

func init() {
	requestQueue = &RequestQueue{list.New()}

}
