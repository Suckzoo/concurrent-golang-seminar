package main

import (
	"fmt"
	"math/rand"
	"time"
)

var steps = [4]string{
	"deployBroserv",
	"deployBroapp",
	"deployBromon",
	"deployPAS",
}

func runStep(stepName string) chan bool {
	ch := make(chan bool)
	go func() {
		duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(duration)
		fmt.Println("stepName: ", stepName, " elapsed: ", duration, "ms")
		ch <- true
	}()
	return ch
}
func main() {
	c := make([]chan bool, 4)
	now := time.Now()
	for i, v := range steps {
		c[i] = runStep(v)
	}
	done := 0
	ch := make(chan bool)
	for i := 0; i < 4; i++ {
		go func(c chan bool) {
			select {
			case <-c:
				done++
				ch <- true
			case <-time.After(500 * time.Millisecond):
				ch <- true
				return
			}
		}(c[i])
	}
	for i := 0; i < 4; i++ {
		<-ch
	}
	elapsed := time.Now().Sub(now)
	fmt.Println("total: ", elapsed)
	fmt.Printf("%d jobs done before timeout\n", done)
}
