package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	safeGo(func() {
		panic("handled panic")
	})

	time.Sleep(time.Second)
	fmt.Println("Program continues")
}

func safeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Goroutine panicked: %v", r)
			}
		}()
		fn()
	}()
}
