package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	numOfWorkers := 3
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	taskChan := make(chan int, 3)
	outChan := make(chan string, 3)
	globalTimeout := 5
	parentCtx, cancel := context.WithTimeout(context.Background(), time.Duration(globalTimeout)*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go worker(parentCtx, taskChan, outChan, &wg)
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
		close(outChan)
	}()

	for val := range outChan {
		fmt.Println(val)
	}

}

func worker(parentCtx context.Context, taskChan <-chan int, outChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-parentCtx.Done():
		return
	default:
		for val := range taskChan {
			select {
			case <-parentCtx.Done():
				return
			case <-time.After(time.Duration(val) * time.Second):
				outChan <- "Task " + strconv.Itoa(val) + " completed"
			}
		}
	}
}
