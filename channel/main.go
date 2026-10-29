package main

import (
	"fmt"
)

func producer(ch chan<- int) {
	for i := range 100 {
		ch <- i
	}
	close(ch)
}

func consumer(ch <-chan int) {
	for val := range ch {
		fmt.Println("Done processing:", val)
	}
}

func main() {
	ch := make(chan int, 10)
	go producer(ch)
	consumer(ch)
}
