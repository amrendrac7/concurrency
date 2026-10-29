package main

import (
	"fmt"
	"time"
)

var result int

func main() {
	go func() {
		result = 42
	}()

	time.Sleep(time.Millisecond)
	fmt.Println("result", result)
}
