package main

import (
	"fmt"
	"time"
)

func toss(how string, ch chan string) {
	for {
		fmt.Println(<-ch)
		time.Sleep(100 * time.Millisecond)
		ch <- how
	}
}

// run with GOTRACEBACK=all
func main() {
	table := make(chan string)
	go toss("ping", table)
	go toss("pong", table)
	// table <- "serve"
	time.Sleep(1 * time.Second)
	<-table
	panic("show me the stack")
}
