package main 

import (
	"fmt"
	"sync"
	"time"
)

func merge(cs ... <-chan int) <-chan int {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c<-chan int) {
			for v := range c {
				ch <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

func toChan(ar ...int) <-chan int {
	ch := make(chan int)
	go func() {
		for _, i := range ar {
			ch<-i
			time.Sleep(time.Second / 5)
		}
		close(ch)
	}()
	return ch
}

func main() {
	a := toChan(1, 3, 5, 7)
	b := toChan(2, 4, 6, 8)
	c := merge(a, b)
	for i := range c {
		fmt.Println(i)
	}
}