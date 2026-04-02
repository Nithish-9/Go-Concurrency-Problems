package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {

	parentCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tasks := []int{1, 2, 3, 4, 5, 6}
	numOfWorkers := 3
	taskChan := make(chan int, 3)
	resultChan := make(chan string, 3)
	errChan := make(chan string, 1)
	wg := sync.WaitGroup{}

	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go worker(parentCtx, taskChan, resultChan, errChan, &wg)
	}

	go func() {
		for _, val := range tasks {
			select {
			case <-parentCtx.Done():
				close(taskChan)
				return
			case taskChan <- val:
			}
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for {
		select {
		case <-parentCtx.Done():
			for pushedValBeforeFailure := range resultChan {
				fmt.Println(pushedValBeforeFailure)
			}
			return
		case val, ok := <-resultChan:
			if !ok {
				fmt.Println("Output Channel Closed")
				return
			}
			fmt.Println(val)
		case val, ok := <-errChan:
			if ok {
				fmt.Println(val)
				cancel()
			}
		}
	}

}

func worker(parentCtx context.Context, taskChan <-chan int, resultChan chan<- string, errChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for val := range taskChan {
		select {
		case <-parentCtx.Done():
			return
		case <-time.After(time.Duration(val) * time.Second):
			if val == 4 {
				errChan <- "Error Occurred at Task " + strconv.Itoa(val)
				return
			} else {
				resultChan <- "Task " + strconv.Itoa(val) + " completed"
			}
		}
	}

}
