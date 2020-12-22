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

func fanIn(c1, c2 chan string, quit chan bool) chan string {
	mux := make(chan string)
	a1 := 0
	a2 := 0
	go func() {
		for {
			select {
			case v := <-c1:
				mux <- v
				a1++
				if a1 >= 5 {
					c1 = nil
				}
			case v := <-c2:
				mux <- v
				a2++
				if a2 >= 5 {
					c2 = nil
				}
			case <-quit:
				fmt.Println("cleanup...")
				quit <- true
				return
			}
		}
	}()
	return mux
}

func main() {
	runtime.GOMAXPROCS(10)
	c1 := generate("aa")
	c2 := generate("bb")
	quit := make(chan bool)
	mux := fanIn(c1, c2, quit)
	defer close(mux)
	for i := 0; i < 10; i++ {
		fmt.Println(<-mux)
	}
	quit <- true
	fmt.Println("quit: ", <-quit)
}
