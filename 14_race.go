package main

import (
	"fmt"
	"time"
)

var sum int

func plus1() {
	for i := 0; i < 10000000; i++ {
		sum++
	}
}

// run with -race flag
func main() {
	go plus1()
	go plus1()
	time.Sleep(1 * time.Second)
	fmt.Println(sum)
}
