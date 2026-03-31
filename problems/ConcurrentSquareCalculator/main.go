package main

import (
	"fmt"
	"sync"
)

/*
Problem Statement

You are given:

nums := []int{1, 2, 3, 4, 5}

Do the following using goroutines + channels:

Launch a separate goroutine for each number
Each goroutine computes the square of the number
Send the result to a channel
Collect all results in the main goroutine and print them
Constraints
No global variables
Use channels correctly (no race conditions)
No goroutine leaks
Output order does NOT need to be sorted
Follow-up

How would you modify your solution to preserve the order of results?
*/
type order struct {
	index int
	num   int
}

func main() {
	nums := []int{1, 2, 3, 4, 5}

	wg := sync.WaitGroup{}
	orderChan := make(chan order, len(nums))
	for i, num := range nums {
		wg.Add(1)
		go func(i int, num int) {
			defer wg.Done()
			orderChan <- order{
				index: i,
				num:   num * num,
			}
		}(i, num)
	}
	wg.Wait()
	close(orderChan)

	for val := range orderChan {
		nums[val.index] = val.num
	}

	fmt.Print(nums)
}
