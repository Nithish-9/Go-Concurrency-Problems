package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	numOfChannels := 3
	chans := make([]chan int, numOfChannels)
	for i := 0; i < numOfChannels; i++ {
		chans[i] = make(chan int)
	}

	go func() {
		for i := 0; i < 9; i++ {
			chans[0] <- i
		}
		close(chans[0])
	}()

	go func() {
		for i := 10; i < 19; i++ {
			chans[1] <- i
		}
		close(chans[1])
	}()

	go func() {
		for i := 20; i < 29; i++ {
			chans[2] <- i
		}
		close(chans[2])
	}()

	outChan := make(chan int)

	go func() {
		for val := range outChan {
			fmt.Print(strconv.Itoa(val) + " ")
		}
	}()

	wg := sync.WaitGroup{}
	for _, ch := range chans {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				outChan <- val
			}
		}(ch)
	}

	wg.Wait()
	close(outChan)

}
