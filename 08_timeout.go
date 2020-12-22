package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var a = [5]int{1, 2, 3, 4, 5}

func generate(marker string) chan string {
	c := make(chan string)
	go func() {
		count := 0
		for {
			c <- fmt.Sprintf("%s %d", marker, count)
			count++
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			if marker == "aa" {
				time.Sleep(1 * time.Second)
			}
		}
		close(c)
	}()
	return c
}

func fanIn(c1, c2 chan string) chan string {
	mux := make(chan string)
	go func() {
		for {
			select {
			case v := <-c1:
				mux <- v
			case v := <-c2:
				mux <- v
			case <-time.After(500 * time.Millisecond):
				mux <- "timeout"
			}
		}
	}()
	return mux
}

func main() {
	runtime.GOMAXPROCS(10)
	c1 := generate("aa")
	c2 := generate("bb")
	mux := fanIn(c1, c2)
	defer close(mux)
	for i := 0; i < 10; i++ {
		fmt.Println(<-mux)
	}
}
