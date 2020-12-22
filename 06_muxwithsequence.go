package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var a = [5]int{1, 2, 3, 4, 5}

func generate(marker string, wait chan bool) chan string {
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
			<-wait
		}
		close(c)
	}()
	return c
}

func fanIn(c1, c2 chan string) chan string {
	mux := make(chan string)
	go func() {
		for {
			mux <- <-c1
		}
	}()
	go func() {
		for {
			mux <- <-c2
		}
	}()
	return mux
}

func main() {
	runtime.GOMAXPROCS(10)
	wait1 := make(chan bool)
	wait2 := make(chan bool)
	c1 := generate("aa", wait1)
	c2 := generate("bb", wait2)
	mux := fanIn(c1, c2)
	defer close(mux)
	for i := 0; i < 5; i++ {
		fmt.Println(<-mux)
		fmt.Println(<-mux)
		wait1 <- true
		wait2 <- true
	}
}
