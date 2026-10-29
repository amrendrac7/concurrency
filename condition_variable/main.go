package main

import (
	"fmt"
	"sync"
)

type DataProcessor struct {
	mu        sync.Mutex
	cond      *sync.Cond
	dataReady bool
	data      string
}

func NewDataProcessor() *DataProcessor {
	dp := &DataProcessor{}
	dp.cond = sync.NewCond(&dp.mu)
	return dp
}
func (dp *DataProcessor) WaitForData() string {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	for !dp.dataReady {
		dp.cond.Wait()
	}
	dp.dataReady = false
	return dp.data
}

func (dp *DataProcessor) ProvideData(newData string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	dp.data = newData
	dp.dataReady = true
	dp.cond.Signal()
}

func main() {
	processor := NewDataProcessor()
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		data := processor.WaitForData()
		fmt.Printf("Received: %s\n", data)
	}()

	go func() {
		defer wg.Done()
		processor.ProvideData("Hello, World!")
	}()

	wg.Wait()
}
