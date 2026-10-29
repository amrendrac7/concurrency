package main

import (
	"fmt"
	"sync"
	"time"
)

type Barrier struct {
	parties    int
	count      int
	generation int
	mu         sync.Mutex
	cond       *sync.Cond
	action     func()
}

func NewBarrier(parties int, action func()) *Barrier {
	b := &Barrier{
		parties: parties,
		count:   0,
		action:  action,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()

	myGeneration := b.generation
	b.count++

	if b.count == b.parties {
		if b.action != nil {
			b.action()
		}
		b.count = 0
		b.generation++
		b.cond.Broadcast()
	} else {
		for b.generation == myGeneration {
			b.cond.Wait()
		}
	}
}

func main() {
	numGoroutines := 3

	barrier := NewBarrier(numGoroutines, func() {
		fmt.Println("All goroutines arrived, startin next phase!")
	})

	var wg sync.WaitGroup
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			fmt.Printf("Goroutine %d doing phase 1\n", id)
			time.Sleep(time.Duration(100*id) * time.Millisecond)

			barrier.Wait()

			fmt.Printf("Goroutine %d doint phase 2\n", id)
		}(i)
	}

	wg.Wait()
}
