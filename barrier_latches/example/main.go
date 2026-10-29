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
	b := &Barrier{parties: parties, action: action}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	myGen := b.generation
	b.count++

	if b.count == b.parties {
		if b.action != nil {
			b.action()
		}
		b.count = 0
		b.generation++
		b.cond.Broadcast()
	} else {
		for b.generation == myGen {
			b.cond.Wait()
		}
	}
}

func main() {
	const workers = 4
	const phases = 3
	data := make([]int, 100)

	barrier := NewBarrier(workers, func() {
		fmt.Println("=== Phase complete, starting next ===")
	})

	var completionWg sync.WaitGroup
	completionWg.Add(workers)

	for w := range workers {
		start := w * 25
		end := start + 25

		go func(workerId, start, end int) {
			defer completionWg.Done()

			for phase := range phases {
				fmt.Printf("Worker %d processing phase %d [%d-%d]\n", workerId, phase, start, end)

				for i := start; i < end; i++ {
					data[i] += phase + 1
				}
				time.Sleep(50 * time.Millisecond)

				barrier.Wait()
			}
			fmt.Printf("Worker %d complete\n", workerId)
		}(w, start, end)
	}

	completionWg.Wait()
	fmt.Println("All computation complete!")
}
