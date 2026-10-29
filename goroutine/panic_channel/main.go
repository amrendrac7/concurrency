package main

import (
	"errors"
	"fmt"
	"log"
)

type Result struct {
	Value int
	Err   error
}

func compute(x int) Result {
	if x < 0 {
		return Result{Err: errors.New("negative input")}
	}
	return Result{Value: x * x}
}

func main() {
	results := make(chan Result)

	go func() {
		results <- compute(-5)
	}()

	r := <-results
	if r.Err != nil {
		log.Fatal(r.Err)
	}

	fmt.Println(r.Value)
}
