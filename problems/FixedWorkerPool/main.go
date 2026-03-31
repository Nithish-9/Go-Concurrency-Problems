package main

import (
	"fmt"
	"sync"
)

/*
Problem Statement

You are given:

nums := []int{1,2,3,4,5,6,7,8,9,10}

Build a worker pool with 3 workers:

Requirements:
Create 3 worker goroutines
Workers should:
Continuously read from a jobs channel
Compute square
Send result to results channel
Main goroutine should:
Send all jobs
Close jobs channel
Collect results
Constraints
Use channels properly
Workers should exit gracefully
No goroutine leaks
Order NOT required
Follow-up

How would you:

Detect when all workers are done?
Close the results channel safely?
*/

type result struct {
	index int
	num   int
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numOfWorkers := 3
	numOfJobs := len(nums)
	jobsChan := make(chan result, numOfJobs)
	resultChan := make(chan result, numOfJobs)
	wg := sync.WaitGroup{}
	for i := 1; i <= numOfWorkers; i++ {
		wg.Add(1)
		go worker(jobsChan, resultChan, &wg)
	}

	for index, value := range nums {
		jobsChan <- result{
			index: index,
			num:   value,
		}
	}
	close(jobsChan)
	wg.Wait()
	close(resultChan)

	for it := range resultChan {
		nums[it.index] = it.num
	}

	fmt.Println(nums)
}

func worker(jobsChan <-chan result, resultChan chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range jobsChan {
		resultChan <- result{
			index: data.index,
			num:   data.num * data.num,
		}
	}
}
