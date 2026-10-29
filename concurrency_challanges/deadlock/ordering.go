package main

import "sync"

type BankAccount struct {
	id      int64
	balance int
	mu      sync.Mutex
}

func NewBankAccount(id int64, initialBalance int) *BankAccount {
	return &BankAccount{id: id, balance: initialBalance}
}

func Transfer(from, to *BankAccount, amount int) {
	// Always lock in order of account ID to prevent deadlock
	first, second := from, to
	if from.id > to.id {
		first, second = to, from
	}

	first.mu.Lock()
	defer first.mu.Unlock()

	second.mu.Lock()
	defer second.mu.Unlock()

	if from.balance > amount {
		from.balance -= amount
		to.balance += amount
	}
}
