package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type data struct {
	index int
	value string
}

func main() {
	requests := [][]int{{1, 2}, {2, 1e9}} // first value -> task, second value -> num of times it fails
	maxRequests := 2

	ans := make([]string, len(requests))
	parentCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tokenChan := make(chan struct{}, maxRequests)
	requestChan := make(chan []int)
	allowDropChan := make(chan data)

	numOfWorkers := 3
	jobChan := make(chan []int, 3)
	processChan := make(chan string, numOfWorkers)
	wg := sync.WaitGroup{}
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go processAllowedRequests(jobChan, processChan, &wg)
	}

	// Token Refiller Channel
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(maxRequests))
		defer ticker.Stop()
		for {
			select {
			case <-parentCtx.Done():
				return
			case <-ticker.C:
				select {
				case tokenChan <- struct{}{}:
				default:
				}
			}
		}
	}()

	// Incoming Requests Channel
	go func() {
		defer close(requestChan)
		for index, inner := range requests {
			select {
			case <-parentCtx.Done():
				return
			default:
				requestChan <- []int{index, inner[0], inner[1]}
				time.Sleep(time.Second / time.Duration(maxRequests))
			}
		}
	}()

	// allow/drop the incoming requests
	go func() {
		defer close(allowDropChan)
		defer close(jobChan)
		for it := range requestChan {
			select {
			case <-parentCtx.Done():
				return
			case <-tokenChan:
				allowDropChan <- data{
					index: it[0],
					value: "ALLOW",
				}
				jobChan <- it
			default:
				allowDropChan <- data{
					index: it[0],
					value: "DROP",
				}
			}
		}
	}()

	go func() {
		for it := range allowDropChan {
			ans[it.index] = it.value
		}
		fmt.Println(ans)
	}()

	go func() {
		for it := range processChan {
			fmt.Println(it)
		}
	}()

	wg.Wait()
	close(processChan)
}

func processAllowedRequests(jobChan <-chan []int, processChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	maxRetries := 2
	for task := range jobChan {
		fmt.Println("Request Processing for ", task[1])
		attempt := 0
		var status string
		for {
			if attempt >= task[2] {
				status = "SUCCESS"
				break
			}

			if attempt == maxRetries {
				status = "FAILED"
				break
			}

			time.Sleep(time.Duration(1<<attempt) * time.Second)

			attempt++
		}
		processChan <- fmt.Sprintf("Task %d %s", task[1], status)
	}
}
