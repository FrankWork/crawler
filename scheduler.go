package main

import (
	"container/list"
	"fmt"
)

var urlQueue = list.New()

func enqueue(url string) {
	urlQueue.PushBack(url)
}

func dequeue(url string) *list.Element {
	e := urlQueue.Front()
	fmt.Println(e.Value)
	urlQueue.Remove(e)
	return e
}
