package deadlock

import "sync"

// AcquireAll tries to acquire all locks atomically.
// Returns true if all locks acquired, false otherwise.
// On failure, no locks are held.
func AcquireAll(mutexes ...*sync.Mutex) bool {
	acquired := make([]*sync.Mutex, 0, len(mutexes))

	for _, mu := range mutexes {
		if mu.TryLock() {
			acquired = append(acquired, mu)
		} else {
			for _, held := range acquired {
				held.Unlock()
			}
			return false
		}
	}

	return true
}

// ReleaseAll unlocks all provided mutexes
func ReleaseAll(mutextes ...*sync.Mutex) {
	for _, mu := range mutextes {
		mu.Unlock()
	}
}

// Usage example:
func operationNeedingAllLocks(mu1, mu2, mu3 *sync.Mutex) bool {
	if AcquireAll(mu1, mu2, mu3) {
		defer ReleaseAll(mu1, mu2, mu3)

		// Critical section all three locks held
		doWork()
		return true
	}
	return false
}

func doWork() {}
