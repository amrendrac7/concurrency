package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore {
	return make(chan struct{}, n)
}

func (s Semaphore) Acquire() {
	s <- struct{}{}
}

func (s Semaphore) Release() {
	<-s
}

type Connection struct {
	ID int
}

func (c *Connection) Execute(query string) {
	fmt.Printf("Connection %d executing: %s\n", c.ID, query)
}

type ConnectionPool struct {
	semaphore Semaphore
	pool      chan *Connection
}

func NewConnectionPool(poolSize int) *ConnectionPool {
	cp := &ConnectionPool{
		semaphore: NewSemaphore(poolSize),
		pool:      make(chan *Connection, poolSize),
	}

	for i := range poolSize {
		cp.pool <- &Connection{ID: i}
	}

	return cp
}

func (cp *ConnectionPool) Acquire() *Connection {
	cp.semaphore.Acquire()
	return <-cp.pool
}

func (cp *ConnectionPool) Release(conn *Connection) {
	cp.pool <- conn
	cp.semaphore.Release()
}

func main() {
	pool := NewConnectionPool(3)

	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(taskId int) {
			defer wg.Done()

			conn := pool.Acquire()
			defer pool.Release(conn)

			conn.Execute(fmt.Sprintf("SELECT * FROM users WHERE id =%d", taskId))
			time.Sleep(100 * time.Millisecond)
		}(i)
	}

	wg.Wait()
}
