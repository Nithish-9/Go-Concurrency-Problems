package main

import (
	"fmt"
	"sync"
)

/*
Problem Statement

You are given:

nums := []int{1,2,3,4,5}

Build a 2-stage pipeline:

Stage 1:
Take input numbers
Compute square
Send to next stage
Stage 2:
Take squared values
Compute double
Send to final output
Example:
Input: 1 → square: 1 → double: 2
Input: 2 → square: 4 → double: 8
Requirements
Each stage runs in separate goroutines
Use channels between stages
Ensure proper closing of channels
No goroutine leaks
Follow-up (VERY IMPORTANT)
What happens if Stage 2 is slow?
How do you add backpressure handling?
*/

type data struct {
	index int
	value int
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	numOfWorkers := 3

	squareChan := make(chan data)
	doubleChan := make(chan data)
	outChan := make(chan data)

	wg1 := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	for i := 1; i <= numOfWorkers; i++ {
		wg1.Add(1)
		wg2.Add(1)
		go squareWorker(squareChan, doubleChan, &wg1)
		go doubleWorker(doubleChan, outChan, &wg2)
	}

	for index, value := range nums {
		squareChan <- data{
			index: index,
			value: value,
		}
	}
	close(squareChan)

	go func() {
		wg1.Wait()
		close(doubleChan)
	}()

	go func() {
		wg2.Wait()
		close(outChan)
	}()

	for it := range outChan {
		nums[it.index] = it.value
	}

	fmt.Println(nums)

}

func squareWorker(squareChan <-chan data, doubleChan chan<- data, wg1 *sync.WaitGroup) {
	defer wg1.Done()
	for it := range squareChan {
		doubleChan <- data{
			index: it.index,
			value: it.value * it.value,
		}
	}
}

func doubleWorker(doubleChan <-chan data, outChan chan<- data, wg2 *sync.WaitGroup) {
	defer wg2.Done()
	for it := range doubleChan {
		outChan <- data{
			index: it.index,
			value: it.value * 2,
		}
	}
}
