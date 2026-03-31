package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Problem 5: Rate Limited Worker
Problem Statement

You are given:

tasks := []int{1,2,3,4,5,6,7,8,9,10}
🎯 Task

Process these tasks using goroutines, but enforce:

❗ Maximum 2 tasks per second

✅ Requirements
Use goroutines + channels
Enforce rate limiting (time-based, not worker count)
Maintain concurrency (don’t make it purely sequential)
No goroutine leaks
Process all tasks
⏱ Expected Behavior (Example)
Time 0s → process task 1,2
Time 1s → process task 3,4
Time 2s → process task 5,6
...
⚠️ Constraints / Hints
You’ll need something like:
time.Ticker OR
time.After OR
a custom rate limiter

Think carefully:

Rate limit ≠ number of workers

🔥 Follow-up (VERY IMPORTANT)
1. How is this different from limiting number of workers?
2. Where is this used in real systems?

(Examples: APIs, external services, etc.)
*/

func main() {
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	wg := sync.WaitGroup{}

	j := 0

	for j < len(tasks) {
		<-ticker.C
		for i := 1; i <= 3; i++ {
			if j < len(tasks) {
				wg.Add(1)
				go func(value int) {
					defer wg.Done()
					fmt.Println("Value is ", value)
				}(tasks[j])
			} else {
				break
			}
			j++
		}
	}

	wg.Wait()
}
