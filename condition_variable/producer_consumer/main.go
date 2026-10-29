package main

import (
	"fmt"
	"sync"
	"time"
)

type BoundedBuffer struct {
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	buffer   []interface{}
	caapcity int
}

func NewBoundedBuffer(capacity int) *BoundedBuffer {
	bb := &BoundedBuffer{
		buffer:   make([]interface{}, 0, capacity),
		caapcity: capacity,
	}
	bb.notEmpty = sync.NewCond(&bb.mu)
	bb.notFull = sync.NewCond(&bb.mu)

	return bb
}

func (bb *BoundedBuffer) Put(item interface{}) {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	for len(bb.buffer) == bb.caapcity {
		bb.notFull.Wait()
	}

	bb.buffer = append(bb.buffer, item)
	bb.notEmpty.Signal()
}

func (bb *BoundedBuffer) Take() interface{} {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	for len(bb.buffer) == 0 {
		bb.notEmpty.Wait()
	}

	item := bb.buffer[0]
	bb.buffer = bb.buffer[1:]
	bb.notFull.Signal()
	return item
}

func main() {
	buffer := NewBoundedBuffer(5)
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range 20 {
			buffer.Put(i)
			fmt.Printf("\033[31m Produced: %d \033[0m \n", i)
		}
	}()

	go func() {
		defer wg.Done()
		for range 20 {
			item := buffer.Take()
			fmt.Printf("\033[32m Consumed: %v \033[0m \n", item)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	wg.Wait()
}
