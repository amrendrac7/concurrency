package deadlock

import (
	"context"
	"sync"
	"time"
)

type BankAccountWithTimeout struct {
	balance int
	mu      sync.Mutex
}

func NewBankAccountWithTimeout(balance int) *BankAccountWithTimeout {
	return &BankAccountWithTimeout{balance: balance}
}

func TransferWithTimeout(ctx context.Context, from, to *BankAccountWithTimeout, amount int) bool {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(5 * time.Second)
	}

	for {
		// Try to acquire first lock
		if from.mu.TryLock() {
			if to.mu.TryLock() {
				defer to.mu.Unlock()
				defer from.mu.Unlock()

				if from.balance >= amount {
					from.balance -= amount
					to.balance += amount
					return true
				}
				return false // Insufficient funds
			}
			// Couldn't get second lock, release first
			from.mu.Unlock()
		}
		// Check timeout
		if time.Now().After(deadline) {
			return false // Timeout
		}

		// Chcek context cancellation
		select {
		case <-ctx.Done():
			return false
		default:
		}

		// back off before retrying
		time.Sleep(time.Millisecond)
	}
}
