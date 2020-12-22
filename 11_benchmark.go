package main

import (
	"fmt"
)

func f(l, r chan int) {
	l <- (1 + <-r)
}

func main() {
	leftmost := make(chan int)
	left := leftmost
	right := leftmost
	for i := 0; i < 100000; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	go func(c chan int) { c <- 1 }(right)
	fmt.Println(<-leftmost)
}
