package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *Counter) IncrementManual() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *Counter) TryIncrement() bool {
	if c.mu.TryLock() {
		defer c.mu.Unlock()
		c.count++
		return true
	}
	return false
}

func (c *Counter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	counter := &Counter{}
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	fmt.Printf("Final count: %d\n", counter.GetCount())
}
