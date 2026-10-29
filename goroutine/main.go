package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d finished\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All workers done")
}
