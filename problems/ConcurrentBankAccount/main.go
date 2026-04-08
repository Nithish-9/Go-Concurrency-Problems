package main

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	mu     sync.RWMutex
	amount int
}

func NewBankAccount(initialBalance int) *BankAccount {
	return &BankAccount{amount: initialBalance}
}

func (b *BankAccount) Deposit(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.amount += amount
}

func (b *BankAccount) Withdraw(amount int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.amount-amount < 0 {
		return false
	}
	b.amount -= amount
	return true
}

func (b *BankAccount) Balance() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.amount
}

func main() {
	account := NewBankAccount(100)
	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Deposit(1)
		}()
	}

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Withdraw(1)
		}()
	}

	wg.Wait()
	fmt.Println("Final Balance:", account.Balance())
}
