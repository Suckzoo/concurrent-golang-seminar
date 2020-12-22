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

var closing chan chan error

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
func runStepBlock() error {
	fmt.Println("===============RUNNING STEPBLOCK=============")
	c := make([]chan string, 4)
	now := time.Now()
	for i, v := range steps {
		c[i] = runStepWithReplica(v)
	}
	done := 0
	ch := make(chan string)
	for i := 0; i < 4; i++ {
		go func(stepName string, c chan string) {
			select {
			case r := <-c:
				ch <- r
			case <-time.After(500 * time.Millisecond):
				ch <- "timeout"
				return
			}
		}(steps[i], c[i])
	}
	for i := 0; i < 4; i++ {
		replica := <-ch
		if replica != "timeout" {
			fmt.Println(steps[i], " is done with replica ", replica)
		}
	}
	elapsed := time.Now().Sub(now)
	fmt.Println("total: ", elapsed)
	fmt.Printf("%d jobs done before timeout\n", done)
	fmt.Println("=======================DONE==================")
	return fmt.Errorf("some error happened")
}

func Close() error {
	errc := make(chan error)
	closing <- errc
	return <-errc
}

func main() {
	req := make(chan bool)
	closing = make(chan chan error)
	go func() {
		<-time.After(5 * time.Second)
		err := Close()
		if err != nil {
			fmt.Println("Error happened!?")
			fmt.Println(err)
		}
	}()
	go func(r chan bool) {
		for {
			r <- true
			time.Sleep(1 * time.Second)
		}
	}(req)
	var err error
	err = nil
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-req:
			err = runStepBlock()
		case errc := <-closing:
			req = nil
			errc <- err
		case <-timeout:
			fmt.Println("shutting down...")
			return
		}
	}
}
