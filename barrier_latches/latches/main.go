package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	numService := 3
	var wg sync.WaitGroup
	wg.Add(numService)

	for i := range numService {
		go func(serviceId int) {
			fmt.Printf("Service %d starting...\n", serviceId)
			time.Sleep(time.Duration(100*(serviceId+1)) * time.Millisecond)
			wg.Done()
		}(i)
	}

	fmt.Println("Main goroutine waiting for services...")
	wg.Wait()
	fmt.Println("All service is ready, stating application!")
}

// channel based latch

type Latch struct {
	c    chan struct{}
	once sync.Once
}

func NewLatch() *Latch {
	return &Latch{c: make(chan struct{})}
}

func (l *Latch) Open() {
	l.once.Do(func() {
		close(l.c)
	})
}

func (l *Latch) Wait() {
	<-l.c
}
