package main

func run() {
	for !requestQueue.isEmpty() {
		rw := requestQueue.dequeue()
		rw.print()
	}
}
