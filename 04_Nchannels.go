package main

import (
	"fmt"
	"time"
	"math/rand"
)

var a = [5]int{1, 2, 3, 4, 5}

func generate(marker string) chan string {
	c := make(chan string)
	go func() {
		for _, v := range a {
			c <- fmt.Sprintf("%s %d", marker, v)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func main() {
	c1 := generate("aa")
	c2 := generate("bb")
	for i := 0; i < 5; i++ {
		fmt.Println(<-c1)
		fmt.Println(<-c2)
	}
}