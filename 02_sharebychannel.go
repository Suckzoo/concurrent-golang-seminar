package main

import (
	"fmt"
	"strings"
	"time"
)

var names = [9]string{"kim", "cho", "won", "a", "b", "c", "d", "e", "ffff"}

func feedNames(source, upper chan string) {
	for _, name := range names {
		source <- name
		fmt.Println(<-upper)
	}
}

func toUpper(source, upper chan string) {
	for s := range source {
		upper <- strings.ToUpper(s)
	}
}

func main() {
	source := make(chan string)
	upper := make(chan string)
	defer close(source)
	defer close(upper)
	go feedNames(source, upper)
	go toUpper(source, upper)
	time.Sleep(10 * time.Second)
}
