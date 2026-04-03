package main

import (
	"context"
	"fmt"
	"time"
)

type data struct {
	index int
	value string
}

func main() {
	requests := []int{1, 1, 2, 2, 2, 3, 3, 3}
	maxRequests := 2

	ans := make([]string, len(requests))
	parentCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tokenChan := make(chan struct{}, maxRequests)
	requestChan := make(chan int)
	outChan := make(chan data)

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
		for index := range requests {
			select {
			case <-parentCtx.Done():
				return
			default:
				requestChan <- index
				time.Sleep(time.Second / time.Duration(maxRequests))
			}
		}
	}()

	// process the incoming requests
	go func() {
		defer close(outChan)
		for val := range requestChan {
			select {
			case <-parentCtx.Done():
				return
			case <-tokenChan:
				outChan <- data{
					index: val,
					value: "ALLOW",
				}
			default:
				outChan <- data{
					index: val,
					value: "DROP",
				}
			}
		}
	}()

	for it := range outChan {
		ans[it.index] = it.value
	}
	fmt.Println(ans)
}
