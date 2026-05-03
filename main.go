package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch1)
		defer close(ch2)
		for i := range 10 {
			if i%2 == 0 {
				ch1 <- i
			} else {
				ch2 <- i
			}
		}
	}()

	data := channelmerge(ch1, ch2)

	for res := range data {
		fmt.Println(res, "была получена из канала")
	}

}
func channelmerge(channels ...chan int) chan int {

	Outchan := make(chan int)

	for _, ch := range channels {
		wg.Add(1)
		go func(ch chan int) {
			defer wg.Done()
			for val := range ch {
				fmt.Println(val, " была отправлена в канал")
				Outchan <- val

			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(Outchan)
	}()

	return Outchan
}
