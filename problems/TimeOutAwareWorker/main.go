package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

/*
Problem Statement

You are given a list of tasks. Each task takes variable time to process.

Your goal is to:

Process all tasks concurrently
Each task must complete within a given timeout
If a task exceeds the timeout, it must be cancelled
Only successful task results should be printed

Input
tasks = [1,2,3,4,5]
timeout = 2 seconds
Each task simulates work like:
Task 1 → completes in 1s
Task 2 → completes in 3s (timeout)
Task 3 → completes in 500ms
Task 4 → completes in 4s (timeout)
Task 5 → completes in 1s

Output
Task 1 completed
Task 3 completed
Task 5 completed

Constraints
Each task runs in its own goroutine
You must enforce timeout per task
Timed-out tasks must NOT print results
No goroutine leaks
Do not block forever

Edge Cases
All tasks exceed timeout
All tasks complete instantly
Timeout is very small (e.g., 1ms)
Tasks finish exactly at timeout boundary
*/

type data struct {
	task          int
	timeOutInSecs int
}

func main() {
	tasks := []data{data{
		task:          1,
		timeOutInSecs: 1,
	}, data{
		task:          2,
		timeOutInSecs: 3,
	}, data{
		task:          3,
		timeOutInSecs: 5,
	}, data{
		task:          4,
		timeOutInSecs: 4,
	}, data{
		task:          5,
		timeOutInSecs: 1,
	}}

	taskTimeLimitInSecs := 2
	wg := sync.WaitGroup{}
	for _, it := range tasks {
		wg.Add(1)
		go worker(it, &wg, taskTimeLimitInSecs)
	}
	wg.Wait()
}

func worker(d data, wg *sync.WaitGroup, taskTimeLimitInSecs int) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(taskTimeLimitInSecs))
	defer cancel()

	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Duration(d.timeOutInSecs) * time.Second):
		fmt.Println("Task " + strconv.Itoa(d.task) + " Completed")
	}
}
