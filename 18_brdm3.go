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

func First(channels ...chan string) chan string {
	ch := make(chan string)
	for _, c := range channels {
		go func(c chan string) {
			select {
			case v := <-c:
				ch <- v
			case <-time.After(500 * time.Millisecond):
				return
			}
		}(c)
	}
	return ch
}

func runStep(stepName string, replica string) chan string {
	ch := make(chan string)
	go func() {
		duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(duration)
		fmt.Println("stepName: ", stepName, " elapsed: ", duration, "ms")
		ch <- replica
	}()
	return ch
}

func runStepWithReplica(stepName string) chan string {
	return First(
		runStep(stepName, "r1"),
		runStep(stepName, "r2"),
		runStep(stepName, "r3"),
	)
}
func main() {
	c := make([]chan string, 4)
	now := time.Now()
	for i, v := range steps {
		c[i] = runStepWithReplica(v)
	}
	done := 0
	ch := make(chan bool)
	for i := 0; i < 4; i++ {
		go func(stepName string, c chan string) {
			select {
			case r := <-c:
				fmt.Println(stepName, " is done with Replica ", r)
				done++
				ch <- true
			case <-time.After(500 * time.Millisecond):
				ch <- true
				return
			}
		}(steps[i], c[i])
	}
	for i := 0; i < 4; i++ {
		<-ch
	}
	elapsed := time.Now().Sub(now)
	fmt.Println("total: ", elapsed)
	fmt.Printf("%d jobs done before timeout\n", done)
}
