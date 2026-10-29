package main

import "sync"

type SemaphoreImpl struct {
	mu      sync.Mutex
	cond    *sync.Cond
	permits int
}

func NewSemaphore(initialPermits int) *SemaphoreImpl {
	s := &SemaphoreImpl{
		permits: initialPermits,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *SemaphoreImpl) Acquire() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for s.permits == 0 {
		s.cond.Wait()
	}
	s.permits--
}

func (s *SemaphoreImpl) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.permits++
	s.cond.Signal()
}

func (s *SemaphoreImpl) TryAcquire() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.permits > 0 {
		s.permits--
		return true
	}
	return false
}
