package main

import (
	"fmt"
	"time"
)

func fetchData(url string, results chan<- string) {
	time.Sleep(time.Millisecond * 100)
	results <- fmt.Sprintf("Data from %s", url)
}

func main() {
	results := make(chan string, 3)

	urls := []string{"url1", "url2", "url3"}
	for _, url := range urls {
		go fetchData(url, results)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Println(<-results)
	}
}
