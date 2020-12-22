package main

import (
	"fmt"
	"time"
)

var a = [8]int{1, 2, 3, 4, 5, 6, 7, 8}

func generate() chan int {
	c := make(chan int)
	go func() {
		for _, v := range a {
			c <- v
			time.Sleep(500 * time.Millisecond)
		}
		close(c)
	}()
	return c
}
func main() {
	c := generate()
	for v := range c {
		fmt.Println(v)
	}
}
