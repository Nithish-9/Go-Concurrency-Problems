package main

import (
	"fmt"
	"sync"
)

/*
Problem Name: Parallel Sum Aggregator
Problem Statement

You are given:

nums := []int{1,2,3,4,5,6,7,8,9,10}
Task:
Split the array into N chunks (say 3)
Each chunk is processed by a goroutine
Each goroutine computes sum of its chunk
Send partial sums to a channel
Main goroutine aggregates final sum
Example:
Chunk1 → sum = 6
Chunk2 → sum = 15
Chunk3 → sum = 34

Final → 55
Constraints
Use goroutines + channels
No race conditions
No global variables
Handle uneven splits properly
Follow-up (important)
What happens if one goroutine is very slow?
How would you add timeout or cancellation?
*/
func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	numOfWorkers := 3
	chunkSize := (len(nums) + numOfWorkers - 1) / numOfWorkers
	i := 0
	j := chunkSize
	resultChan := make(chan int, chunkSize)
	finalSum := 0
	wg := sync.WaitGroup{}
	for {
		wg.Add(1)
		if j >= len(nums) {
			go worker(nums[i:], resultChan, &wg)
			break
		} else {
			go worker(nums[i:j], resultChan, &wg)
		}
		i = j
		j += chunkSize
	}
	wg.Wait()
	close(resultChan)

	for value := range resultChan {
		finalSum += value
	}

	fmt.Println(finalSum)
}

func worker(chunk []int, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for _, value := range chunk {
		sum += value
	}
	resultChan <- sum
}
